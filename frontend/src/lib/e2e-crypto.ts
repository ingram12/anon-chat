export class E2ECryptoHandler {
    privateKey!: CryptoKey;
    publicKey!: CryptoKey;
    publicKeyBase64!: string;
    peerPublicKey!: CryptoKey;
  
    // Generate RSA key pair and export publicKey as base64
    async generateKeyPair(): Promise<void> {
      const keyPair = await crypto.subtle.generateKey(
        {
          name: 'RSA-OAEP',
          modulusLength: 2048,
          publicExponent: new Uint8Array([1, 0, 1]),
          hash: 'SHA-256',
        },
        true,
        ['encrypt', 'decrypt']
      );
  
      this.privateKey = keyPair.privateKey;
      this.publicKey = keyPair.publicKey;
  
      const spki = await crypto.subtle.exportKey('spki', this.publicKey);
      this.publicKeyBase64 = btoa(String.fromCharCode(...new Uint8Array(spki)));
    }
  
    // Set peer's public key from base64
    async setPeerPublicKey(base64: string): Promise<void> {
      const binary = Uint8Array.from(atob(base64), c => c.charCodeAt(0));
      this.peerPublicKey = await crypto.subtle.importKey(
        'spki',
        binary,
        { name: 'RSA-OAEP', hash: 'SHA-256' },
        true,
        ['encrypt']
      );
    }
  
    // Encrypt message: AES → encrypt AES key with RSA
    async encrypt(message: string): Promise<string> {
      if (!this.peerPublicKey) throw new Error('Peer public key is not set.');
  
      const aesKey = await crypto.subtle.generateKey(
        { name: 'AES-GCM', length: 256 },
        true,
        ['encrypt', 'decrypt']
      );
  
      const iv = crypto.getRandomValues(new Uint8Array(12));
      const encoded = new TextEncoder().encode(message);
  
      const ciphertext = await crypto.subtle.encrypt(
        { name: 'AES-GCM', iv },
        aesKey,
        encoded
      );
  
      const rawAesKey = await crypto.subtle.exportKey('raw', aesKey);
      const encryptedAesKey = await crypto.subtle.encrypt(
        { name: 'RSA-OAEP' },
        this.peerPublicKey,
        rawAesKey
      );
  
      return JSON.stringify({
        aesKey: btoa(String.fromCharCode(...new Uint8Array(encryptedAesKey))),
        ciphertext: btoa(String.fromCharCode(...new Uint8Array(ciphertext))),
        iv: btoa(String.fromCharCode(...iv)),
      });
    }
  
    // Decrypt: RSA → AES → message
    async decrypt(payload: string): Promise<string> {
      if (!this.privateKey) throw new Error('Private key is not set.');
  
      const { aesKey, ciphertext, iv } = JSON.parse(payload);
  
      const encKey = Uint8Array.from(atob(aesKey), c => c.charCodeAt(0));
      const rawKey = await crypto.subtle.decrypt(
        { name: 'RSA-OAEP' },
        this.privateKey,
        encKey
      );
  
      const aesKeyObj = await crypto.subtle.importKey(
        'raw',
        rawKey,
        { name: 'AES-GCM' },
        true,
        ['decrypt']
      );
  
      const ct = Uint8Array.from(atob(ciphertext), c => c.charCodeAt(0));
      const ivBytes = Uint8Array.from(atob(iv), c => c.charCodeAt(0));
  
      const decrypted = await crypto.subtle.decrypt(
        { name: 'AES-GCM', iv: ivBytes },
        aesKeyObj,
        ct
      );
  
      return new TextDecoder().decode(decrypted);
    }
  }

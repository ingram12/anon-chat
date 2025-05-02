import { E2ECryptoHandler } from '../lib/e2e-crypto';

describe('E2ECryptoHandler', () => {
  let alice: E2ECryptoHandler;
  let bob: E2ECryptoHandler;
  
  beforeEach(async () => {
    alice = new E2ECryptoHandler();
    bob = new E2ECryptoHandler();
    
    // Generate keypairs for both parties
    await alice.generateKeyPair();
    await bob.generateKeyPair();
  });

  test('should generate different key pairs', () => {
    expect(alice.publicKeyBase64).not.toBe(bob.publicKeyBase64);
  });

  test('should successfully encrypt and decrypt message', async () => {
    // Setup peers
    await alice.setPeerPublicKey(bob.publicKeyBase64);
    await bob.setPeerPublicKey(alice.publicKeyBase64);

    const originalMessage = 'Hello, World!';
    
    // Alice encrypts message for Bob
    const encrypted = await alice.encrypt(originalMessage);
    
    // Bob decrypts message from Alice
    const decrypted = await bob.decrypt(encrypted);
    
    expect(decrypted).toBe(originalMessage);
  });

  test('should handle empty messages', async () => {
    await alice.setPeerPublicKey(bob.publicKeyBase64);
    
    const encrypted = await alice.encrypt('');
    const decrypted = await bob.decrypt(encrypted);
    
    expect(decrypted).toBe('');
  });

  test('should handle Unicode characters', async () => {
    await alice.setPeerPublicKey(bob.publicKeyBase64);
    
    const message = '👋 Hello, 世界!';
    const encrypted = await alice.encrypt(message);
    const decrypted = await bob.decrypt(encrypted);
    
    expect(decrypted).toBe(message);
  });

  test('should throw error when peer public key is not set', async () => {
    await expect(alice.encrypt('test'))
      .rejects
      .toThrow('Peer public key is not set');
  });

  test('should throw error when private key is not set', async () => {
    const handler = new E2ECryptoHandler();
    await expect(handler.decrypt('{}'))
      .rejects
      .toThrow('Private key is not set');
  });

  test('should throw error on invalid encrypted payload', async () => {
    await bob.generateKeyPair();
    
    await expect(bob.decrypt('{"invalid":"json"}'))
      .rejects
      .toThrow();
  });

  test('should not decrypt message encrypted for another recipient', async () => {
    const eve = new E2ECryptoHandler();
    await eve.generateKeyPair();
    
    await alice.setPeerPublicKey(bob.publicKeyBase64);
    const encrypted = await alice.encrypt('secret message');

    await expect(eve.decrypt(encrypted))
      .rejects
      .toThrow();
  });
});
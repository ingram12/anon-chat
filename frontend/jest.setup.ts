import { webcrypto } from 'node:crypto';

if (!globalThis.crypto?.subtle) {
    // Set up the crypto environment
    globalThis.crypto = webcrypto as unknown as Crypto;
}
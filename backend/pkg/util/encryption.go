package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// ---------------------------------------------------------------------
// Utilities
// ---------------------------------------------------------------------

// generateRandomBytes returns securely generated random bytes.
func generateRandomBytes(n int) []byte {
	b := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}
	return b
}

// encryptAESGCM encrypts plaintext using key with a random nonce.
func encryptAESGCM(key, plaintext []byte) (ciphertext, nonce []byte, tag []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	nonce = generateRandomBytes(aead.NonceSize())
	encrypted := aead.Seal(nil, nonce, plaintext, nil)

	// Split into ciphertext + tag
	ciphertext = encrypted[:len(encrypted)-aead.Overhead()]
	tag = encrypted[len(encrypted)-aead.Overhead():]
	return
}

// decryptAESGCM decrypts ciphertext using key, nonce, and tag.
func decryptAESGCM(key, ciphertext, nonce, tag []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	combined := append(ciphertext, tag...)
	plaintext, err := aead.Open(nil, nonce, combined, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// ---------------------------------------------------------------------
// Envelope Encryption (Service Secret Demo)
// ---------------------------------------------------------------------

type EncryptedSecret struct {
	Ciphertext []byte
	Nonce      []byte
	Tag        []byte

	WrappedDEK []byte
	WrapNonce  []byte
	WrapTag    []byte
	KekVersion int
}

// encryptSecretWithKEK performs envelope encryption:
//  1. Generate DEK
//  2. Encrypt secret with DEK (AES-GCM)
//  3. Wrap DEK with KEK[v] (AES-GCM)
func encryptSecretWithKEK(secret, kek []byte, kekVersion int) EncryptedSecret {
	// Step 1: generate DEK
	dek := generateRandomBytes(32)

	// Step 2: encrypt secret with DEK
	ct, nonce, tag := encryptAESGCM(dek, secret)

	// Step 3: wrap DEK with KEK
	wrapped, wrapNonce, wrapTag := encryptAESGCM(kek, dek)

	return EncryptedSecret{
		Ciphertext: ct,
		Nonce:      nonce,
		Tag:        tag,

		WrappedDEK: wrapped,
		WrapNonce:  wrapNonce,
		WrapTag:    wrapTag,
		KekVersion: kekVersion,
	}
}

// decryptSecretWithKEK reverses the envelope encryption.
func decryptSecretWithKEK(record EncryptedSecret, kek []byte) ([]byte, error) {
	// Step 1: unwrap DEK
	dek, err := decryptAESGCM(kek, record.WrappedDEK, record.WrapNonce, record.WrapTag)
	if err != nil {
		return nil, fmt.Errorf("failed to unwrap DEK: %w", err)
	}

	// Step 2: decrypt secret with DEK
	plaintext, err := decryptAESGCM(dek, record.Ciphertext, record.Nonce, record.Tag)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt secret: %w", err)
	}

	// Wipe DEK (best-effort)
	for i := range dek {
		dek[i] = 0
	}

	return plaintext, nil
}

// ---------------------------------------------------------------------
// Demo
// ---------------------------------------------------------------------

func Demo() {
	// Simulate KEK[v1] stored securely (here just in memory)
	kekV1 := generateRandomBytes(32)

	secret := []byte("my-github-api-token-123456")

	// ---- ENCRYPT & STORE ----
	record := encryptSecretWithKEK(secret, kekV1, 1)

	fmt.Println("=== Stored in DB ===")
	fmt.Println("Ciphertext:", base64.StdEncoding.EncodeToString(record.Ciphertext))
	fmt.Println("Nonce:", base64.StdEncoding.EncodeToString(record.Nonce))
	fmt.Println("Tag:", base64.StdEncoding.EncodeToString(record.Tag))
	fmt.Println("WrappedDEK:", base64.StdEncoding.EncodeToString(record.WrappedDEK))
	fmt.Println("WrapNonce:", base64.StdEncoding.EncodeToString(record.WrapNonce))
	fmt.Println("WrapTag:", base64.StdEncoding.EncodeToString(record.WrapTag))
	fmt.Println("KEK Version:", record.KekVersion)

	// ---- READ & DECRYPT ----
	plaintext, err := decryptSecretWithKEK(record, kekV1)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n=== Decrypted Secret ===")
	fmt.Printf("Plaintext: %s\n", plaintext)
}

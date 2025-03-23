package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
)

// Encrypt Returns a Base64 encoded ciphertext encrypted with AES-GCM 256 using key
func Encrypt(plaintextString string, keyB64 string) (string, error) {
	// Format plaintext
	plaintext := []byte(plaintextString)

	// Decode key for encryption
	key, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil {
		return "", err
	}

	// Generate a nonce
	nonce := make([]byte, 12) // AES-GCM recommends a 12-byte nonce
	_, err = rand.Read(nonce)
	if err != nil {
		return "", err
	}

	// Create the AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create GCM cipher mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Encrypt the data
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	// Encode the ciphertext into Base64
	encodedCiphertext := base64.StdEncoding.EncodeToString(ciphertext)

	return encodedCiphertext, nil
}

// Decrypt decrypts the plaintext in AES-GCM 256 using key
func Decrypt(ciphertextB64 string, keyB64 string) (string, error) {
	// Decode ciphertext
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return "", err
	}

	// Decode key
	key, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil {
		return "", err
	}

	// Separate nonce and ciphertext
	nonce, ciphertext := ciphertext[:12], ciphertext[12:]

	// Create the AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create the GCM cipher mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Decrypt the data
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	plaintextString := string(plaintext)

	return plaintextString, nil
}

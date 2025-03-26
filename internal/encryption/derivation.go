package encryption

import (
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/argon2"
)

// =-- Argon2 Parameters --= //

// Argon2Params holds the standardized Argon2 parameters for key derivation/hash
type Argon2Params struct {
	Time    uint32 // Number of iterations
	Memory  uint32 // Memory cost (KB)
	Threads uint8  // Number of threads
	KeyLen  uint32 // Length of derived key (bytes)
}

var DefaultArgon2Params = Argon2Params{
	Time:    6,         // 6 iterations
	Memory:  64 * 1024, // 64 MB
	Threads: 4,         // 4 threads
	KeyLen:  64,        // For two 32-byte key for AES-256 (K_enc, K_auth)
}

// =-- Primary Functions --= //

// GenerateSalt Generates a randomized salt given size in bytes
func GenerateSalt(size int) (string, error) {
	salt := make([]byte, size) // Create a byte slice of size
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	// Encode salt into Base64
	encodedSalt := base64.StdEncoding.EncodeToString(salt)

	return encodedSalt, nil
}

// DeriveMasterKeys Derives two keys (encryption, auth) from password using argon2 given salt
func DeriveMasterKeys(password string, saltB64 string) (string, string, error) {
	// Decode and format salt
	salt, err := base64.StdEncoding.DecodeString(saltB64)
	if err != nil {
		return "", "", err
	}

	// Derive key
	params := DefaultArgon2Params
	masterKey := argon2.IDKey([]byte(password), salt, params.Time, params.Memory, params.Threads, params.KeyLen)

	// Split the master 64-byte key into K_auth and K_enc
	encryptionKey := masterKey[:32]
	authenticationKey := masterKey[32:]

	return base64.StdEncoding.EncodeToString(encryptionKey), base64.StdEncoding.EncodeToString(authenticationKey), nil
}

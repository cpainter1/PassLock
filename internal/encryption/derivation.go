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
	KeyLen:  32,        // 32-byte key for AES-256
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

// DeriveKey Derives a key from password using argon2 given salt
func DeriveKey(password string, saltB64 string) string {
	// Decode and format salt
	salt, err := base64.StdEncoding.DecodeString(saltB64)

	if err != nil {
		panic(err)
	}

	// Derive key
	params := DefaultArgon2Params
	key := argon2.IDKey([]byte(password), salt, params.Time, params.Memory, params.Threads, params.KeyLen)

	// Encode key into Base64
	encodedKey := base64.StdEncoding.EncodeToString(key)

	return encodedKey
}

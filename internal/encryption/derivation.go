package encryption

import (
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/argon2"
)

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

	time := uint32(2)        // 2 iterations
	mem := uint32(64 * 1024) // 64 MB
	threads := uint8(4)      // Use 4 threads for computation
	keyLen := uint32(32)     // 32 bytes for 256-bit key

	// Derive key
	key := argon2.Key([]byte(password), salt, time, mem, threads, keyLen)

	// Encode key into Base64
	encodedKey := base64.StdEncoding.EncodeToString(key)

	return encodedKey
}

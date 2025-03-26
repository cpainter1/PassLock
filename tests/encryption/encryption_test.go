package tests

import (
	"github.com/cpainter1/PassLock/internal/encryption"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	// Generate a test key
	password := "verysimplepassword"
	salt, err := encryption.GenerateSalt(16)
	if err != nil {
		t.Fatalf("GenerateSalt failed in TestEncryptDecrypt: %v", err)
	}

	key, _, _ := encryption.DeriveMasterKeys(password, salt)
	t.Logf("Derived key: %s", key)

	// Encrypt a sample plaintext
	plaintext := "Here's a very secret message..."

	ciphertext, err := encryption.Encrypt(plaintext, key)
	t.Logf("Encrypted ciphertext: %s", ciphertext)

	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Decrypt the ciphertext
	decryptedPlaintext, err := encryption.Decrypt(ciphertext, key)
	t.Logf("Decrypted plaintext: %s", decryptedPlaintext)

	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	// Check if the decrypted plaintext and the original unencrypted plaintext are equal
	if decryptedPlaintext != plaintext {
		t.Errorf("Decrypt failed: decrypted plaintext does not match")
	}
}

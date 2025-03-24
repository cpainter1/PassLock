package tests

import (
	"github.com/cpainter1/PassLock/internal/encryption"
	"testing"
)

func TestGenerateSalt(t *testing.T) {
	// Constant size
	saltSizeBytes := 16
	saltSizeB64 := 24 // Since 16-bytes encoded in Base64 is ~24 characters in length

	// Generate salt1
	salt1, err := encryption.GenerateSalt(saltSizeBytes)
	t.Logf("Generate salt1: %s", salt1)
	if err != nil {
		t.Fatalf("GenerateSalt salt1 failed: %v", err)
	}

	// Compare salt1 and size
	if len(salt1) != saltSizeB64 {
		t.Errorf("Expected salt1 length of %d, but got %d", saltSizeB64, len(salt1))
	}

	// Generate salt2
	salt2, err := encryption.GenerateSalt(saltSizeBytes)
	t.Logf("Generate salt2 %s", salt2)
	if err != nil {
		t.Fatalf("GenerateSalt salt2 failed: %v", err)
	}

	// Compare salt2 and size
	if len(salt2) != saltSizeB64 {
		t.Errorf("Expected salt2 length of %d, but got %d", saltSizeB64, len(salt2))
	}
}

func TestDeriveKey(t *testing.T) {
	// Constant passwords
	password1 := "securepassword"
	password2 := "superspecialpassword123"

	// Derive key1 and key2 on salt1
	salt1, err := encryption.GenerateSalt(16)
	if err != nil {
		t.Fatalf("GenerateSalt failed in TestDeriveKey: %v", err)
	}
	key1 := encryption.DeriveKey(password1, salt1)
	t.Logf("DeriveKey 1 key: %x", key1)

	key2 := encryption.DeriveKey(password1, salt1)
	t.Logf("DeriveKey 2 key: %x", key2)

	// Ensure key1 and key2 are of equal length
	if key1 != key2 {
		t.Errorf("DeriveKey failed: expected %x and %x to be equal!", key1, key2)
	}

	// Derive key3 on salt2
	salt2, err := encryption.GenerateSalt(16)
	if err != nil {
		t.Fatalf("GenerateSalt (diff) failed in TestDeriveKey: %v", err)
	}
	key3 := encryption.DeriveKey(password1, salt2)
	t.Logf("DeriveKey 3 key: %x", key3)

	if key1 == key3 {
		t.Errorf("DeriveKey failed: expected %x and %x to differ since they have differet salt", key1, key3)
	}

	// Derive key 4
	key4 := encryption.DeriveKey(password2, salt2)
	t.Logf("DeriveKey 4 key: %x", key4)
}

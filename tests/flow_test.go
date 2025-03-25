package tests

import (
	"github.com/cpainter1/PassLock/internal/database"
	"github.com/cpainter1/PassLock/internal/encryption"
	"testing"
)

// TestFullEncryptionFlow tests the full application flow without UI (encryption -> db storage -> decryption)
func TestFullEncryptionFlow(t *testing.T) {
	// =-- Input parameters (simulated user-provided) --= //
	inputService := "https://mail.google.com"
	inputUsername := "example@gmail.com"
	inputPassword := "password123"
	inputNotes := "my gmail account password"

	inputMasterPassword := "supersecretpassword321"

	// =-- Encryption --= //

	// Derivation & salt
	encryptionSalt, err := encryption.GenerateSalt(16)
	if err != nil {
		t.Errorf("Error generating salt: %v", err)
	}

	derivedMasterPassword := encryption.DeriveKey(inputMasterPassword, encryptionSalt)
	t.Logf("DerivedMasterPassword: %v", derivedMasterPassword)

	// Password and Notes Encryption
	encryptedPassword, err := encryption.Encrypt(inputPassword, derivedMasterPassword)
	if err != nil {
		t.Errorf("Error encrypting password: %v", err)
	}
	t.Logf("Encrypted password: %s", encryptedPassword)

	encryptedNotes, err := encryption.Encrypt(inputNotes, derivedMasterPassword)
	if err != nil {
		t.Errorf("Error encrypting notes: %v", err)
	}
	t.Logf("Encrypted notes: %s", encryptedNotes)

	// =-- Storing in Database --= //
	entry := database.PasswordEntry{
		Service:           inputService,
		Username:          inputUsername,
		EncryptedPassword: encryptedPassword,
		EncryptedNotes:    encryptedNotes,
	}

	db, err := database.InitDB()
	if err != nil {
		t.Errorf("Error connecting to database: %v", err)
	}

	returnedInformation, err := database.StorePassword(db, entry)
	if err != nil {
		t.Errorf("Error storing password: %v", err)
	}
	t.Logf("Returned information: %v", returnedInformation)

	// =-- Retrieving and Decrypting from Database --= //
	retrievedEntry, err := database.GetEntryFromID(db, returnedInformation.ID)
	if err != nil {
		t.Errorf("Error retrieving entry: %v", err)
	}
	t.Logf("Retrieved entry: %v", retrievedEntry)

	// Decrypting
	decryptedPassword, err := encryption.Decrypt(retrievedEntry.EncryptedPassword, derivedMasterPassword)
	if err != nil {
		t.Errorf("Error decrypting password: %v", err)
	}
	t.Logf("Decrypted password: %s", decryptedPassword)

	decryptedNotes, err := encryption.Decrypt(retrievedEntry.EncryptedNotes, derivedMasterPassword)
	if err != nil {
		t.Errorf("Error decrypting notes: %v", err)
	}
	t.Logf("Decrypted notes: %s", decryptedNotes)

	// =-- Comparisons --= //
	if decryptedNotes != inputNotes || decryptedPassword != inputPassword {
		t.Fatalf("Decrypted password %v does not match the input password %v", decryptedPassword, inputPassword)
	}

	// =-- Database Clearing --= //
	err = database.ClearDatabase(db)
	if err != nil {
		t.Errorf("Error clearing database: %v", err)
	}

}

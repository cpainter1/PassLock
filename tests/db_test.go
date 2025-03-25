package tests

import (
	"github.com/cpainter1/PassLock/internal/database"
	"testing"
)

// TestStorePassword tests the password storing database function, simultaneously tests the GetEntriesFromService func
func TestStorePassword(t *testing.T) {
	// Initialize db instance
	db, err := database.InitDB()
	if err != nil {
		t.Errorf("Error connecting to database: %v", err)
	}

	// Sample password information
	var entry = database.PasswordEntry{
		Service:           "https://mail.google.com",
		Username:          "example@gmail.com",
		EncryptedPassword: "password123",
		EncryptedNotes:    "this is a test",
	}

	// Store the password entry
	returnInformation, err := database.StorePassword(db, entry)
	if err != nil {
		t.Errorf("Error storing password: %v", err)
	}
	t.Logf("Returned inserted entry information: %v", returnInformation)

	// Confirm results
	confirmedInformation, err := database.GetEntriesFromService(db, "https://mail.google.com")
	if err != nil {
		t.Errorf("Error retrieving entries: %v", err)
	}
	t.Logf("First queried inserted entry: %v", confirmedInformation[0])

	if confirmedInformation[0].EncryptedPassword != returnInformation.EncryptedPassword {
		t.Fatalf("Returned password does not match the queried password")
	}

	// Clear the database for the next test
	err = database.ClearDatabase(db)
	if err != nil {
		t.Errorf("Error clearing database: %v", err)
	}
}

// TestGetEntryFromID creates a sample password entry, obtains its ID, and then queries it using GetEntryFromID
func TestGetEntryFromID(t *testing.T) {
	// Initialize db instance
	db, err := database.InitDB()
	if err != nil {
		t.Errorf("Error connecting to database: %v", err)
	}

	// Create and store password entry
	var entry = database.PasswordEntry{
		Service:           "https://mail.google.com",
		Username:          "example@gmail.com",
		EncryptedPassword: "password123",
		EncryptedNotes:    "this is a test",
	}

	returnInformation, err := database.StorePassword(db, entry)
	if err != nil {
		t.Errorf("Error storing password: %v", err)
	}
	t.Logf("Returned inserted entry information: %v", returnInformation)

	returnID := returnInformation.ID

	// Query the entry using GetEntryFromID
	queriedEntry, err := database.GetEntryFromID(db, returnID)
	if err != nil {
		t.Errorf("Error retrieving entry information: %v", err)
	}
	t.Logf("Returned entry information: %v", queriedEntry)

	// Compare results
	if returnInformation.EncryptedPassword != queriedEntry.EncryptedPassword {
		t.Fatalf("Returned password %v does not match the queried password %v",
			returnInformation.EncryptedPassword, queriedEntry.EncryptedPassword)
	}

	// Clear database for the next test
	err = database.ClearDatabase(db)
	if err != nil {
		t.Errorf("Error clearing database: %v", err)
	}
}

// TestDeleteEntryFromID creates a password entry and deletes it, querying the whole database and testing GetAllEntries
func TestDeleteEntryFromID(t *testing.T) {
	// Initialize db instance
	db, err := database.InitDB()
	if err != nil {
		t.Errorf("Error connecting to database: %v", err)
	}

	// Create and store password entry
	var entry = database.PasswordEntry{
		Service:           "https://mail.google.com",
		Username:          "example@gmail.com",
		EncryptedPassword: "password123",
		EncryptedNotes:    "this is a test",
	}

	returnInformation, err := database.StorePassword(db, entry)
	if err != nil {
		t.Errorf("Error storing password: %v", err)
	}
	t.Logf("Returned inserted entry information: %v", returnInformation)

	// Delete the password entry
	err = database.DeleteEntryFromID(db, returnInformation.ID)
	if err != nil {
		t.Errorf("Error deleting entry information: %v", err)
	}

	// Query all database entries to confirm deletion
	databaseEntries, err := database.GetAllEntries(db)
	if err != nil {
		t.Errorf("Error retrieving all database entries: %v", err)
	}
	t.Logf("All database entries: %v", databaseEntries)

	if len(databaseEntries) != 0 {
		t.Fatalf("Database entries should be empty: %v", databaseEntries)
	}

	// Clear database for next test
	err = database.ClearDatabase(db)
	if err != nil {
		t.Errorf("Error clearing database: %v", err)
	}
}

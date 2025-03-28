package database

import (
	"database/sql"
	"errors"
	"log"
)

// =-- Standardized EncryptedPassword Entry Data Structures --= //

// PasswordInformation stores information for **output** password entry row dumps
type PasswordInformation struct {
	ID                int    // Unique ID
	Service           string // Service (e.g., "github.com")
	Username          string // Username for the account
	EncryptedPassword string // Encrypted password
	EncryptedNotes    string // Encrypted notes
	CreatedAt         string // Timestamp for entry creation
}

// PasswordEntry stores information for **input** password entries
type PasswordEntry struct {
	Service           string // Service (e.g., "github.com")
	Username          string // Username for the account
	EncryptedPassword string // Encrypted password
	EncryptedNotes    string // Encrypted notes
}

// =-- Database Management Functions --= //

// StorePassword stores a new password entry in the database and returns the row information as PasswordInformation struct
func StorePassword(db *sql.DB, entry PasswordEntry) (*PasswordInformation, error) {
	// Insert SQL query to add a new entry to the passwords table
	insertSQL := `
    INSERT INTO passwords (service, username, password, notes) 
    VALUES (?, ?, ?, ?);`

	// Execute the query with the parameters (service, username, encrypted password, and encrypted notes)
	result, err := db.Exec(
		insertSQL,
		entry.Service,
		entry.Username,
		entry.EncryptedPassword,
		entry.EncryptedNotes)
	if err != nil {
		log.Printf("Error inserting password: %v", err)
		return nil, err
	}

	// Get last inserted row
	lastID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error retrieving last insert ID: %v", err)
		return nil, err
	}

	// Retrieve inserted row
	query := `
    SELECT id, service, username, password, notes, created_at 
    FROM passwords 
    WHERE id = ? LIMIT 1;`

	row := db.QueryRow(query, lastID)

	var inserted PasswordInformation
	err = row.Scan(&inserted.ID, &inserted.Service, &inserted.Username, &inserted.EncryptedPassword, &inserted.EncryptedNotes, &inserted.CreatedAt)
	if err != nil {
		log.Printf("Error fetching inserted password entry with ID %d: %v", lastID, err)
		return nil, err
	}

	return &inserted, nil
}

// GetEntryFromID retrieves all entry information for a unique entry ID
func GetEntryFromID(db *sql.DB, id int) (*PasswordInformation, error) {
	// Query to retrieve the entire row information for a specific ID
	query := `
    SELECT id, service, username, password, notes, created_at 
    FROM passwords 
    WHERE id = ? LIMIT 1;`

	// Query the database and fetch the result
	row := db.QueryRow(query, id)

	// Initialize a PasswordInformation struct to hold the data
	var entry PasswordInformation

	// Scan the row into the PasswordInformation struct
	err := row.Scan(&entry.ID, &entry.Service, &entry.Username, &entry.EncryptedPassword, &entry.EncryptedNotes, &entry.CreatedAt)
	if err != nil {
		log.Printf("Error fetching password entry with ID %d: %v", id, err)
		return nil, err
	}

	// Return the populated PasswordInformation struct
	return &entry, nil
}

// GetEntriesFromService GetEntryFromService retrieves all password entries from a given service
func GetEntriesFromService(db *sql.DB, service string) ([]*PasswordInformation, error) {
	// Query to retrieve all entries for the given service
	query := `
    SELECT id, service, username, password, notes, created_at 
    FROM passwords 
    WHERE service = ?;`

	// Execute the query and get the rows
	rows, err := db.Query(query, service)
	if err != nil {
		log.Printf("Error fetching entries for service '%s': %v", service, err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}(rows)

	// Slice to hold all the password entries
	var entries []*PasswordInformation

	// Loop through the rows and scan each one into a PasswordInformation
	for rows.Next() {
		var entry PasswordInformation
		err := rows.Scan(&entry.ID, &entry.Service, &entry.Username, &entry.EncryptedPassword, &entry.EncryptedNotes, &entry.CreatedAt)
		if err != nil {
			log.Printf("Error reading row for service '%s': %v", service, err)
			return nil, err
		}

		// Append the entry to the slice
		entries = append(entries, &entry)
	}

	// Check if there was an error during iteration
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return nil, err
	}

	// Return all the entries for the given service
	return entries, nil
}

// GetAllEntries returns all entries in the database as a list of PasswordInformation structs
func GetAllEntries(db *sql.DB) ([]*PasswordInformation, error) {
	// Set up SQL query
	query := `
	SELECT id, service, username, password, notes, created_at
	FROM passwords;`

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error fetching entries for all entries: %v", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}(rows)

	// Iterate through rows
	var entries []*PasswordInformation
	for rows.Next() {
		var entry PasswordInformation
		// Input row information into PasswordInformation struct
		err := rows.Scan(
			&entry.ID,
			&entry.Service,
			&entry.Username,
			&entry.EncryptedPassword,
			&entry.EncryptedNotes,
			&entry.CreatedAt)
		if err != nil {
			log.Printf("Error reading row for entry with ID %d: %v", entry.ID, err)
			return nil, err
		}
		entries = append(entries, &entry)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return nil, err
	}

	return entries, nil
}

// DeleteEntryFromID deletes a specific password entry based on unique ID
func DeleteEntryFromID(db *sql.DB, id int) error {
	// Query to delete the entry with the given ID
	deleteSQL := `DELETE FROM passwords WHERE id = ?;`

	// Execute the DELETE query
	_, err := db.Exec(deleteSQL, id)
	if err != nil {
		log.Printf("Error deleting password entry with ID %d: %v", id, err)
		return err
	}

	log.Printf("EncryptedPassword entry with ID %d has been deleted.", id)
	return nil
}

// ClearDatabase clears all entries in the passwords table
func ClearDatabase(db *sql.DB) error {
	// SQL query to delete all rows from the passwords table
	deleteSQL := `DELETE FROM passwords;`

	// Execute the delete query
	_, err := db.Exec(deleteSQL)
	if err != nil {
		log.Printf("Error clearing database: %v", err)
		return err
	}

	return nil
}

// GetSaltFromVault returns the AuthKey salt from a given vault
func GetSaltFromVault(vaultName string) (string, error) {
	db, err := InitDB(vaultName)
	if err != nil {
		return "", err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}(db)

	var salt string
	err = db.QueryRow("SELECT salt FROM vault_metadata WHERE vault_name = ?;", vaultName).Scan(&salt)
	if err != nil {
		log.Printf("Error reading vault metadata: %v", err)
		return "", err
	}

	return salt, nil
}

// AuthenticateVault returns whether the user is authenticated for a specific vault given an authKey
func AuthenticateVault(vaultName string, authKey string) (bool, error) {
	db, err := InitDB(vaultName)
	if err != nil {
		return false, err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println(err)
		}
	}(db)

	// Retrieve the stored authentication key for the given key
	var storedAuthKey string
	err = db.QueryRow("SELECT auth_key FROM vault_metadata WHERE vault_name = ?", vaultName).Scan(&storedAuthKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Vault does not exist
			log.Printf("No vault found with name %s", vaultName)
			return false, nil
		}
		log.Printf("Error retrieving metadata: %s", err)
		return false, err
	}

	// Compare provided authKey with vault metadata authKey
	verificationResult := authKey == storedAuthKey

	if verificationResult {
		return true, nil // Authenticated
	} else {
		log.Printf("Vault %s not authenticated", vaultName)
		return false, nil
	}
}

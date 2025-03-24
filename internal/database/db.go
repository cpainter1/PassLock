package database

import (
	"database/sql"
	"log"
)

// =-- Standardized Password Entry Data Structure --= //

type PasswordEntry struct {
	ID        int    // Unique ID
	Service   string // Service (e.g., "github.com")
	Username  string // Username for the account
	Password  []byte // Encrypted password
	Notes     []byte // Encrypted notes
	CreatedAt string // Timestamp for entry creation
}

// =-- Database Management Functions --= //

// StorePassword stores a new password entry in the database
func StorePassword(db *sql.DB, service, username string, encryptedPassword, encryptedNotes []byte) error {
	// Insert SQL query to add a new entry to the passwords table
	insertSQL := `
    INSERT INTO passwords (service, username, password, notes) 
    VALUES (?, ?, ?, ?);`

	// Execute the query with the parameters (service, username, encrypted password, and encrypted notes)
	_, err := db.Exec(insertSQL, service, username, encryptedPassword, encryptedNotes)
	if err != nil {
		log.Printf("Error inserting password: %v", err)
		return err
	}
	return nil
}

// GetEntryFromID retrieves all entry information for a unique entry ID
func GetEntryFromID(db *sql.DB, id int) (*PasswordEntry, error) {
	// Query to retrieve the entire row information for a specific ID
	query := `
    SELECT id, service, username, password, notes, created_at 
    FROM passwords 
    WHERE id = ? LIMIT 1;`

	// Query the database and fetch the result
	row := db.QueryRow(query, id)

	// Initialize a PasswordEntry struct to hold the data
	var entry PasswordEntry

	// Scan the row into the PasswordEntry struct
	err := row.Scan(&entry.ID, &entry.Service, &entry.Username, &entry.Password, &entry.Notes, &entry.CreatedAt)
	if err != nil {
		log.Printf("Error fetching password entry with ID %d: %v", id, err)
		return nil, err
	}

	// Return the populated PasswordEntry struct
	return &entry, nil
}

// GetEntryFromService retrieves all password entries from a given service
func GetEntriesFromService(db *sql.DB, service string) ([]PasswordEntry, error) {
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
	defer rows.Close()

	// Slice to hold all the password entries
	var entries []PasswordEntry

	// Loop through the rows and scan each one into a PasswordEntry
	for rows.Next() {
		var entry PasswordEntry
		err := rows.Scan(&entry.ID, &entry.Service, &entry.Username, &entry.Password, &entry.Notes, &entry.CreatedAt)
		if err != nil {
			log.Printf("Error reading row for service '%s': %v", service, err)
			return nil, err
		}

		// Append the entry to the slice
		entries = append(entries, entry)
	}

	// Check if there was an error during iteration
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return nil, err
	}

	// Return all the entries for the given service
	return entries, nil
}

// DeletePassword deletes a specific password entry based on unique ID
func DeletePassword(db *sql.DB, id int) error {
	// Query to delete the entry with the given ID
	deleteSQL := `DELETE FROM passwords WHERE id = ?;`

	// Execute the DELETE query
	_, err := db.Exec(deleteSQL, id)
	if err != nil {
		log.Printf("Error deleting password entry with ID %d: %v", id, err)
		return err
	}

	log.Printf("Password entry with ID %d has been deleted.", id)
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

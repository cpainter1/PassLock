package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"runtime"

	_ "github.com/mattn/go-sqlite3" // Used to init SQLite driver
)

// GetDatabasePath returns the OS-specific path for plockdb.sqlite
func GetDatabasePath() string {
	var basePath string

	switch runtime.GOOS {
	case "darwin": // MacOS
		basePath = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "PassLock")
	case "linux":
		basePath = filepath.Join(os.Getenv("HOME"), ".config", "passlock")
	case "windows":
		basePath = filepath.Join(os.Getenv("APPDATA"), "PassLock")
	default:
		basePath = "." // Fallback to current directory
	}

	// Ensure the directory exists
	err := os.MkdirAll(basePath, 0700)
	if err != nil {
		log.Printf("Could not create database directory: %s", err)
		return ""
	}

	return filepath.Join(basePath, "plockdb.sqlite")
}

func InitDB() (*sql.DB, error) {
	dbPath := GetDatabasePath()

	// Check if the database exists; if not, create it
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		file, err := os.OpenFile(dbPath, os.O_CREATE, 0600)
		if err != nil {
			log.Printf("Error creating database: %s", err)
			return nil, err
		}

		// Close the file
		err = file.Close()
		if err != nil {
			log.Printf("Error closing file: %v", err)
			return nil, err
		}
	}

	// Open SQLite database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return nil, err
	}

	// Create passwords table if it does not already exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS passwords (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		service TEXT NOT NULL,
		username TEXT NOT NULL,
		password BLOB NOT NULL, -- AES-256 encrypted
		notes BLOB,             -- Optional encrypted field
		created_at TIMESTAMP DEFAULT CURrENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return nil, err
	}

	return db, nil
}

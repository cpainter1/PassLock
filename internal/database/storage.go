package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"runtime"

	_ "github.com/mattn/go-sqlite3" // REQUIRED - Used to init SQLite driver
)

// GetVaultDirectoryPath returns the vault directory path depending on OS
func GetVaultDirectoryPath() string {
	var vaultPath string

	switch runtime.GOOS {
	case "darwin": // MacOS
		vaultPath = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "PassLock", "vaults")
	case "linux":
		vaultPath = filepath.Join(os.Getenv("HOME"), ".config", "passlock", "vaults")
	case "windows":
		vaultPath = filepath.Join(os.Getenv("APPDATA"), "PassLock", "vaults")
	default:
		vaultPath = "./vaults" // Fallback to current directory
	}

	return vaultPath
}

// GetDatabasePath returns the OS-specific path for a specific vault
func GetDatabasePath(vaultName string) string {
	var basePath string

	switch runtime.GOOS {
	case "darwin": // MacOS
		basePath = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "PassLock", "vaults")
	case "linux":
		basePath = filepath.Join(os.Getenv("HOME"), ".config", "passlock", "vaults")
	case "windows":
		basePath = filepath.Join(os.Getenv("APPDATA"), "PassLock", "vaults")
	default:
		basePath = "./vaults" // Fallback to current directory
	}

	// Ensure the directory exists
	err := os.MkdirAll(basePath, 0700)
	if err != nil {
		log.Printf("Could not create database directory: %s", err)
		return ""
	}

	return filepath.Join(basePath, vaultName+".sqlite")
}

// CreateVault creates an SQLite vault given vaultName and authentication key authKey
func CreateVault(vaultName string, hashedAuthKey string, authKeySalt string) error {
	dbPath := GetDatabasePath(vaultName)

	// Ensure the vault does not already exist
	if _, err := os.Stat(dbPath); err == nil {
		log.Printf("Vault %s already exists", vaultName)
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	// Create the SQLite database file
	file, err := os.OpenFile(dbPath, os.O_CREATE, 0600)
	if err != nil {
		log.Printf("Error creating database: %s", err)
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Error closing database: %s", err)
		}
	}(file)

	// Open the database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("Error opening database: %s", err)
		return err
	}

	// Create necessary tables
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS passwords (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    service TEXT NOT NULL,
	    username TEXT NOT NULL,
	    password TEXT NOT NULL, -- AES-256 encrypted
	    notes TEXT,             -- Optional encrypted field
	    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS vault_metadata (
	    vault_name TEXT PRIMARY KEY,
	    auth_key TEXT NOT NULL,
	    salt TEXT NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Printf("Error creating table: %s", err)
		return err
	}

	// Store the authentication key in vault_metadata
	_, err = db.Exec(
		"INSERT INTO vault_metadata (vault_name, auth_key, salt) VALUES (?, ?, ?)",
		vaultName,
		hashedAuthKey,
		authKeySalt,
	)
	if err != nil {
		log.Printf("Error inserting metadata: %s", err)
		return err
	}
	return nil
}

// InitDB returns a database instance for an existing database
func InitDB(vaultName string) (*sql.DB, error) {
	dbPath := GetDatabasePath(vaultName)

	// Check if the database exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Printf("Vault %s does not exist", vaultName)
		return nil, err
	}

	// Open SQLite database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("Error opening database: %s", err)
		return nil, err
	}

	return db, nil
}

// ListVaults lists all vaults in the vault directory
func ListVaults() ([]string, error) {
	vaultDir := GetVaultDirectoryPath()

	// Read files in directory
	files, err := os.ReadDir(vaultDir)
	if err != nil {
		return nil, err
	}

	// Iterate through files
	var vaults []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sqlite" {
			vaultName := file.Name()[:len(file.Name())-len(".sqlite")] // Remove .sqlite extension
			vaults = append(vaults, vaultName)
		}
	}

	return vaults, nil
}

// DeleteVault deletes a specific vault given its file name
func DeleteVault(vaultName string) error {
	vaultPath := GetDatabasePath(vaultName)

	// Check if vault exists
	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		log.Printf("Error: vault %s does not exist", vaultName)
		return err
	}

	// Attempt to remove vault
	err := os.Remove(vaultPath)
	if err != nil {
		log.Printf("Error deleting vault: %v", err)
		return err
	}

	return nil
}

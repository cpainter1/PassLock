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

func InitDB(vaultName string) (*sql.DB, error) {
	dbPath := GetDatabasePath(vaultName)

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
		password TEXT NOT NULL, -- AES-256 encrypted
		notes TEXT,             -- Optional encrypted field
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Printf("Error creating table: %v", err)
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

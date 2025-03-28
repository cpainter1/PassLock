package tests

import (
	"github.com/cpainter1/PassLock/internal/database"
	"testing"
)

// TestGetDatabasePath returns the database path of a sample vault
func TestGetDatabasePath(t *testing.T) {
	databasePath := database.GetDatabasePath("TestingVault")
	t.Logf("Database path: %s", databasePath)
}

// TestListVaults
func TestListVaults(t *testing.T) {
	// Create sample vaults
	err := database.CreateVault("TestingVault1", "hashedKey1", "salt1")
	if err != nil {
		t.Errorf("Error creating vault: %v", err)
	}
	_, err = database.InitDB("TestingVault1")
	if err != nil {
		t.Errorf("Error connecting to database: %v", err)
	}

	err = database.CreateVault("TestingVault2", "hashedKey2", "salt2")
	_, err = database.InitDB("TestingVault2")
	if err != nil {
		t.Errorf("Error connecting to database: %v", err)
	}

	err = database.CreateVault("TestingVault3", "hashedKey3", "salt3")
	_, err = database.InitDB("TestingVault3")
	if err != nil {
		t.Errorf("Error connecting to database: %v", err)
	}

	// List vaults
	vaultList, err := database.ListVaults()
	if err != nil {
		t.Errorf("Error in ListVaults: %v", err)
	}
	t.Logf("Vault list: %v", vaultList)

	// Delete created vaults
	err = database.DeleteVault("TestingVault1")
	if err != nil {
		t.Errorf("Error deleting vault: %v", err)
	}

	err = database.DeleteVault("TestingVault2")
	if err != nil {
		t.Errorf("Error deleting vault: %v", err)
	}

	err = database.DeleteVault("TestingVault3")
	if err != nil {
		t.Errorf("Error deleting vault: %v", err)
	}
}

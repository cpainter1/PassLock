package tests

import (
	"github.com/cpainter1/PassLock/internal/database"
	"testing"
)

func TestGetDatabasePath(t *testing.T) {
	databasePath := database.GetDatabasePath()
	t.Logf("Database path: %s", databasePath)
}

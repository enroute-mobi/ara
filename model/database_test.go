package model

import (
	"testing"

	"github.com/af83/edwig/config"
)

func initTestDb(t *testing.T) {
	// Load configuration
	err := config.LoadConfig("")
	if err != nil {
		t.Fatal(err)
	}

	// Initialize Database
	Database = InitDB(config.Config.TestDB)
}

func cleanTestDb(t *testing.T) {
	err := Database.TruncateTables()
	if err != nil {
		t.Fatal(err)
	}
	CloseDB(Database)
}

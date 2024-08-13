package database_test

import (
	"testing"

	"bitbucket.org/junglee_games/getsetgo/pandora/database"
)

func TestCreateNewTestDatabase(t *testing.T) {
	// Create the initial database
	database.CreateNewTestDatabase(t, "pandora_db_testing")

	// Attempting to create it again should crash with db already exists
	database.CreateNewTestDatabase(t, "pandora_db_testing")
}

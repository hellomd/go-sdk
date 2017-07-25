package mongo

import (
	"testing"
)

func TestTestDB(t *testing.T) {
	testDB := NewTestDB()

	if testDB.DB == nil {
		t.Error("Expected DB not to be nil")
	}

	if testDB.DBName == "" {
		t.Error("Expected DBName not to be empty")
	}

	// Make sure session is still open
	testDB.DB.Session.Copy()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Session was not closed")
		}
	}()

	testDB.Close()
	testDB.DB.Session.Copy()
}

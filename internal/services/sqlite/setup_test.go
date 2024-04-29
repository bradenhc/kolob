// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package sqlite_test

import (
	"database/sql"
	"path"
	"testing"

	"github.com/bradenhc/kolob/internal/services/sqlite"
)

func TestCreateTables(t *testing.T) {
	tempdir := t.TempDir()
	dbpath := path.Join(tempdir, "create-table-test.db")

	db, err := sqlite.Open(dbpath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	err = sqlite.CreateTables(db)
	if err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}

	checkTable(t, db, "info")
	checkTable(t, db, "member")
	checkTable(t, db, "conversation")
	checkTable(t, db, "message")
	checkTable(t, db, "mediates")
}

func checkTable(t *testing.T, db *sql.DB, name string) {
	var ret string
	err := db.QueryRow("SELECT name FROM sqlite_master WHERE name = ?", name).Scan(&ret)
	if err != nil {
		t.Errorf("failed to check for %s table: %v", name, err)
	}
	if ret != name {
		t.Errorf("missing table: %s", name)
	}
}

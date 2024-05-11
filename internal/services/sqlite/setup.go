// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package sqlite

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/bradenhc/kolob/internal/fail"
	_ "modernc.org/sqlite"
)

// Open creates a new connection to an SQLite database on the filesystem at the provided path. It
// also initialies the database and sets up the tables needed by Kolob.
//
// The returned DB handle is safe to use throughout the lifetime of the program and by multiple
// goroutines; therefore, Open should only be called once when the program starts.
func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	err = CreateTables(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	return db, nil
}

// CreateTables uses the provided database connection to create the database tables used by Kolob.
// The tables are created inside a transaction: either the all get created or none of them get
// created. It is safe to call this function multiple times even if the tables already exist, or if
// new tables have been added since the last call.
//
// CreateTables is called explicitly in the call to Open, so you rarely need to call this method
// directly; however, it is exported in case you have your own database connection you want to
// setup tables for.
func CreateTables(db *sql.DB) error {
	_, err := db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return fail.Format("failed to enable foreign keys", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to open transaction to create tabels: %v", err)
	}
	defer tx.Rollback()

	slog.Info("Setting up table: group")
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS [group] (
			id 		TEXT,
			idhash	BLOB,
			psalt	BLOB,
			phash	BLOB,
			ekey 	BLOB,
			created TEXT,
			updated TEXT,
			data	BLOB,
	
			PRIMARY KEY (id),
			UNIQUE (idhash)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create group table: %v", err)
	}

	slog.Info("Setting up table: member")
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS member (
			id 			TEXT,
			created		TEXT,
			updated		TEXT,
			idhash		BLOB,
			phash		BLOB,
			data 		BLOB,
			
			PRIMARY KEY (id),
			UNIQUE (idhash)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create member table: %v", err)
	}

	slog.Info("Setting up table: conversation")
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS conversation (
			id			TEXT,
			created		TEXT,
			updated		TEXT,
			data		BLOB,

			PRIMARY KEY (id)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create conversation table: %v", err)
	}

	slog.Info("Setting up table: message")
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS message (
			id				TEXT,
			conversation	TEXT,
			author			TEXT,
			created			TEXT,
			updated			TEXT,
			data			BLOB,

			PRIMARY KEY (id),
			FOREIGN KEY (conversation) 	REFERENCES conversation(id) ON DELETE CASCADE,
			FOREIGN KEY (author) 		REFERENCES member(id) 		ON DELETE SET NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create message table: %v", err)
	}

	slog.Info("Setting up table: moderates")
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS moderates (
			mid TEXT,
			cid TEXT,

			PRIMARY KEY(mid, cid),
			FOREIGN KEY (mid) REFERENCES member(id) 	  ON DELETE CASCADE,
			FOREIGN KEY (cid) REFERENCES conversation(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create moderates table: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit create tables transaction: %v", err)
	}

	return nil
}

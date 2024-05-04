// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package sqlite

import (
	"database/sql"
	"fmt"
	"log/slog"
)

// CreateTables uses the provided database connection to create the database tables used by Kolob.
// The tables are created inside a transaction: either the all get created or none of them get
// created. It is safe to call this function multiple times even if the tables already exist, or if
// new tables have been added since the last call.
func CreateTables(db *sql.DB) error {
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
			data	BLOB,
	
			PRIMARY KEY (id),
			UNIQUE (ghash)
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
			UNIQUE (uhash)
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
			FOREIGN KEY (conversation) 	REFERENCES conversation(id),
			FOREIGN KEY (author) 		REFERENCES member(id)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create message table: %v", err)
	}

	slog.Info("Setting up table: mediates")
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS mediates (
			mid TEXT,
			cid TEXT,

			PRIMARY KEY(mid, cid),
			FOREIGN KEY (mid) REFERENCES member(id),
			FOREIGN KEY (cid) REFERENCES conversation(id)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create mediates table: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit create tables transaction: %v", err)
	}

	return nil
}

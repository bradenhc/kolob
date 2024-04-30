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

	slog.Info("Setting up table: info")
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS info (
			id 		TEXT,
			ghash	BLOB,
			psalt	BLOB,
			phash	BLOB,
			ekey 	BLOB,
			data	BLOB,
	
			PRIMARY KEY (id),
			UNIQUE (gid)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create info table: %v", err)
	}

	slog.Info("Setting up table: member")
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS member (
			id 			TEXT,
			created		TEXT,
			updated		TEXT,
			uhash		BLOB,
			phash		BLOB,
			data 		BLOB,
			
			PRIMARY KEY (id)
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
			id			TEXT,
			author		TEXT,
			created		TEXT,
			updated		TEXT,
			data		BLOB,

			PRIMARY KEY (id),
			FOREIGN KEY (author) REFERENCES member(id)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create message table: %v", err)
	}

	slog.Info("Setting up table: mediates")
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS mediates (
			m_id TEXT,
			c_id TEXT,

			PRIMARY KEY(m_id, c_id),
			FOREIGN KEY (m_id) REFERENCES member(id),
			FOREIGN KEY (c_id) REFERENCES conversation(id)
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

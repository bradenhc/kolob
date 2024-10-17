// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package sqlite

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// Open creates a new connection to an SQLite database on the filesystem at the provided path and
// initializes it with settings needed to support the various stores.
//
// The returned DB handle is safe to use throughout the lifetime of the program and by multiple
// goroutines; therefore, Open should only be called once when the program starts.
func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %v", err)
	}

	return db, nil
}

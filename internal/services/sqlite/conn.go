// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package sqlite

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

// Open creates a new connection to an SQLite database on the filesystem at the provided path.
//
// The returned DB handle is safe to use throughout the lifetime of the program and by multiple
// goroutines; therefore, Open should only be called once when the program starts.
func Open(path string) (*sql.DB, error) {
	return sql.Open("sqlite", path)
}

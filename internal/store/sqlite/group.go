// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/bradenhc/kolob/internal/store"
)

type GroupStore struct {
	db *sql.DB
}

func NewGroupStore(db *sql.DB) (GroupStore, error) {
	slog.Info("Setting up table: group")
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS [group] (
			id 		TEXT,
			ghash	BLOB,
			psalt	BLOB,
			phash	BLOB,
			ekey 	BLOB,
			created INTEGER,
			updated INTEGER,
			data	BLOB,
	
			PRIMARY KEY (id),
			UNIQUE (idhash)
		)
	`)
	if err != nil {
		var s GroupStore
		return s, fmt.Errorf("failed to create group table: %v", err)
	}

	return GroupStore{db}, nil
}

func (s GroupStore) SetGroupData(ctx context.Context, m store.GroupMetadata, d []byte) error {
	slog.Info("Adding group information to SQLite database")
	_, err := s.db.ExecContext(
		ctx,
		"INSERT INTO [group] VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		m.Id[:], m.GroupHash, m.PassSalt, m.PassHash, m.EncryptedKey, m.CreatedAt,
		m.UpdatedAt, d,
	)
	if err != nil {
		return fmt.Errorf("failed to store group information in database: %v", err)
	}

	return nil
}

func (s GroupStore) IsGroupDataSet(ctx context.Context) (bool, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM [group]").Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check for existing group information", err)
	}

	return count != 0, nil
}

func (s GroupStore) GetGroupMetadata(ctx context.Context) (store.GroupMetadata, error) {
	var m store.GroupMetadata
	err := s.db.QueryRowContext(
		ctx, "SELECT id, ghash, psalt, phash, ekey, created, updated FROM [group]",
	).Scan(
		&m.Id, &m.GroupHash, &m.PassSalt, &m.PassHash, &m.EncryptedKey, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return m, fmt.Errorf("failed to get group data from database: %v", err)
	}

	return m, nil
}

func (s GroupStore) GetGroupData(ctx context.Context) ([]byte, error) {
	var d []byte
	err := s.db.QueryRowContext(ctx, "SELECT data FROM [group]").Scan(&d)
	if err != nil {
		return nil, fmt.Errorf("failed to get group data from database: %v", err)
	}

	return d, nil
}

func (s GroupStore) UpdateGroupData(ctx context.Context, m store.GroupMetadata, d []byte) error {
	query := "UPDATE [group] SET ghash = ?, psalt = ?, phash = ?, ekey = ?, updated = ?, data = ?"
	_, err := s.db.ExecContext(
		ctx, query, m.GroupHash, m.PassSalt, m.PassHash, m.EncryptedKey, m.UpdatedAt, d,
	)
	if err != nil {
		return fmt.Errorf("failed to store updated group data in database: %v", err)
	}

	return nil
}

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
			UNIQUE (ghash)
		)
	`)
	if err != nil {
		var s GroupStore
		return s, fmt.Errorf("failed to create group table: %v", err)
	}

	return GroupStore{db}, nil
}

func (s GroupStore) AddGroupEntity(ctx context.Context, e store.GroupEntity) error {
	slog.Info("Adding group information to sqlite database")
	_, err := s.db.ExecContext(
		ctx,
		"INSERT INTO [group] VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		e.Id, e.GroupHash[:], e.PassSalt, e.PassHash, e.EncryptedKey, e.CreatedAt,
		e.UpdatedAt, e.EncryptedData,
	)
	if err != nil {
		return fmt.Errorf("failed to store group entity in sqlite database: %v", err)
	}

	return nil
}

func (s GroupStore) IsGroupDataSet(ctx context.Context) (bool, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM [group]").Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check for existing group information: %v", err)
	}

	return count != 0, nil
}

func (s GroupStore) GetGroupEntity(ctx context.Context) (store.GroupEntity, error) {
	var e store.GroupEntity
	var ghash []byte
	err := s.db.QueryRowContext(ctx, "SELECT * FROM [group]").Scan(
		&e.Id, &ghash, &e.PassSalt, &e.PassHash, &e.EncryptedKey, &e.CreatedAt, &e.UpdatedAt,
		&e.EncryptedData,
	)
	if err != nil {
		return e, fmt.Errorf("failed to get group entity from sqlite database: %v", err)
	}

	copy(e.GroupHash[:], ghash)

	return e, nil
}

func (s GroupStore) UpdateGroupEntity(ctx context.Context, e store.GroupEntity) error {
	query := "UPDATE [group] SET ghash = ?, psalt = ?, phash = ?, ekey = ?, updated = ?, data = ?"
	_, err := s.db.ExecContext(
		ctx, query, e.GroupHash[:], e.PassSalt, e.PassHash, e.EncryptedKey, e.UpdatedAt,
		e.EncryptedData,
	)
	if err != nil {
		return fmt.Errorf("failed to update group entity in sqlite database: %v", err)
	}

	return nil
}

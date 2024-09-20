// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/bradenhc/kolob/internal/crypto"
	"github.com/bradenhc/kolob/internal/model"
	"github.com/bradenhc/kolob/internal/store"
)

type MemberStore struct {
	db *sql.DB
}

func NewMemberStore(db *sql.DB) (MemberStore, error) {

	slog.Info("Setting up table: member")
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS member (
			id 			TEXT,
			uhash		BLOB,
			phash		BLOB,
			created		INTEGER,
			updated		INTEGER,
			data 		BLOB,
			
			PRIMARY KEY (id),
			UNIQUE (uhash)
		)
	`)
	if err != nil {
		var ms MemberStore
		return ms, fmt.Errorf("failed to create member table: %v", err)
	}

	return MemberStore{db}, nil
}

func (s MemberStore) AddMemberData(ctx context.Context, m store.MemberMetadata, d []byte) error {
	slog.Info("Adding member information to SQLite database")
	_, err := s.db.ExecContext(
		ctx,
		"INSERT INTO member VALUES (?, ?, ?, ?, ?, ?)",
		m.Id, m.CreatedAt, m.UpdatedAt, m.UsernameHash[:], m.PassHash, d,
	)
	if err != nil {
		return fmt.Errorf("failed to store new member in database: %v", err)
	}

	return nil
}

func (s MemberStore) GetMemberData(
	ctx context.Context, id model.Uuid,
) (store.MemberMetadata, []byte, error) {
	var m store.MemberMetadata
	var d []byte
	err := s.db.QueryRowContext(
		ctx, "SELECT id, uhash, phash, created, updated, data FROM [member] WHERE id = ?", id,
	).Scan(
		&m.Id, &m.UsernameHash, &m.PassHash, &m.CreatedAt, &m.UpdatedAt, &d,
	)
	if err != nil {
		return m, d, fmt.Errorf("failed to get member data from SQLite database: %v", err)
	}

	return m, d, nil
}

func (s MemberStore) GetMemberDataByUname(
	ctx context.Context, uhash crypto.DataHash,
) (store.MemberMetadata, []byte, error) {
	var m store.MemberMetadata
	var d []byte
	err := s.db.QueryRowContext(
		ctx, "SELECT id, uhash, phash, created, updated, data FROM [member] WHERE uhash = ?",
		uhash[:],
	).Scan(
		&m.Id, &m.UsernameHash, &m.PassHash, &m.CreatedAt, &m.UpdatedAt, &d,
	)
	if err != nil {
		return m, d, fmt.Errorf("failed to get member data by uname from SQLite database: %v", err)
	}

	return m, d, nil
}

func (s MemberStore) UpdateMemberData(ctx context.Context, m store.MemberMetadata, d []byte) error {
	query := "UPDATE [member] SET uhash = ?, phash = ?, updated = ?, data = ? WHERE id = ?"
	_, err := s.db.ExecContext(
		ctx, query, m.UsernameHash, m.PassHash, m.UpdatedAt, d, m.Id[:],
	)
	if err != nil {
		return fmt.Errorf("failed to store updated member data in database: %v", err)
	}

	return nil
}

func (s MemberStore) RemoveMemberData(ctx context.Context, id model.Uuid) error {
	query := "DELETE FROM [member] WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to remove member with id %s: %v", id, err)
	}
	return nil
}

func (s MemberStore) ListMemberData(ctx context.Context) ([][]byte, error) {
	query := "SELECT data FROM [member]"
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get member list from database: %v", err)
	}
	defer rows.Close()

	ms := make([][]byte, 0)
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, fmt.Errorf("failed to scan member row: %v", err)
		}

		ms = append(ms, data)
	}

	return ms, nil
}

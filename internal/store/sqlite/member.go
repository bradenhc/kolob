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

func (s MemberStore) AddMemberEntity(ctx context.Context, e store.MemberEntity) error {
	slog.Info("Adding member information to sqlite database")
	_, err := s.db.ExecContext(
		ctx,
		"INSERT INTO member VALUES (?, ?, ?, ?, ?, ?)",
		e.Id, e.UsernameHash[:], e.PassHash, e.CreatedAt, e.UpdatedAt, e.EncryptedData,
	)
	if err != nil {
		return fmt.Errorf("failed to store new member in database: %v", err)
	}

	return nil
}

func (s MemberStore) GetMemberEntity(
	ctx context.Context, id model.Uuid,
) (store.MemberEntity, error) {
	var e store.MemberEntity
	var uhash []byte
	err := s.db.QueryRowContext(
		ctx, "SELECT * FROM [member] WHERE id = ?", id,
	).Scan(
		&e.Id, &uhash, &e.PassHash, &e.CreatedAt, &e.UpdatedAt, &e.EncryptedData,
	)
	if err != nil {
		var e store.MemberEntity
		return e, fmt.Errorf("failed to get member data from sqlite database: %v", err)
	}

	copy(e.UsernameHash[:], uhash)

	return e, nil
}

func (s MemberStore) GetMemberEntityByUname(
	ctx context.Context, uhash crypto.DataHash,
) (store.MemberEntity, error) {
	var e store.MemberEntity
	var uh []byte
	err := s.db.QueryRowContext(
		ctx, "SELECT * FROM [member] WHERE uhash = ?",
		uhash[:],
	).Scan(
		&e.Id, &uh, &e.PassHash, &e.CreatedAt, &e.UpdatedAt, &e.EncryptedData,
	)
	if err != nil {
		var e store.MemberEntity
		return e, fmt.Errorf("failed to get member data by uname from sqlite database: %v", err)
	}

	copy(e.UsernameHash[:], uh)

	return e, nil
}

func (s MemberStore) UpdateMemberEntity(ctx context.Context, e store.MemberEntity) error {
	query := "UPDATE [member] SET uhash = ?, phash = ?, updated = ?, data = ? WHERE id = ?"
	_, err := s.db.ExecContext(
		ctx, query, e.UsernameHash[:], e.PassHash, e.UpdatedAt, e.EncryptedData, e.Id[:],
	)
	if err != nil {
		return fmt.Errorf("failed to store updated member data in database: %v", err)
	}

	return nil
}

func (s MemberStore) RemoveMemberEntity(ctx context.Context, id model.Uuid) error {
	query := "DELETE FROM [member] WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to remove member with id %s: %v", id, err)
	}
	return nil
}

func (s MemberStore) ListMemberEntities(ctx context.Context) ([]store.MemberEntity, error) {
	query := "SELECT * FROM [member]"
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get member list from database: %v", err)
	}
	defer rows.Close()

	ms := make([]store.MemberEntity, 0)
	for rows.Next() {
		var e store.MemberEntity
		var uhash []byte
		err := rows.Scan(
			&e.Id, &uhash, &e.PassHash, &e.CreatedAt, &e.UpdatedAt, &e.EncryptedData,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan member row: %v", err)
		}

		copy(e.UsernameHash[:], uhash)

		ms = append(ms, e)
	}

	return ms, nil
}

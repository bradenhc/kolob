// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/bradenhc/kolob/internal/model"
	"github.com/bradenhc/kolob/internal/store"
)

type ConversationStore struct {
	db *sql.DB
}

func NewConversationStore(db *sql.DB) (ConversationStore, error) {
	slog.Info("Setting up table: conversation")
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS conversation (
			id			TEXT,
			created		INTEGER,
			updated		INTEGER,
			data		BLOB,

			PRIMARY KEY (id)
		)
	`)
	if err != nil {
		var store ConversationStore
		return store, fmt.Errorf("failed to create conversation table: %v", err)
	}

	return ConversationStore{db}, nil
}

func (s ConversationStore) AddConversationEntity(
	ctx context.Context, e store.ConversationEntity,
) error {
	slog.Info("Adding conversation information to sqlite database")
	_, err := s.db.ExecContext(
		ctx,
		"INSERT INTO [conversation] VALUES (?, ?, ?, ?)",
		e.Id, e.CreatedAt, e.UpdatedAt, e.EncryptedData,
	)
	if err != nil {
		return fmt.Errorf("failed to store new conversation in database: %v", err)
	}

	return nil
}

func (s ConversationStore) GetConversationEntity(
	ctx context.Context, id model.Uuid,
) (store.ConversationEntity, error) {
	var e store.ConversationEntity
	err := s.db.QueryRowContext(
		ctx, "SELECT id, created, updated, data FROM [conversation] WHERE id = ?", id,
	).Scan(
		&e.Id, &e.CreatedAt, &e.UpdatedAt, &e.EncryptedData,
	)
	if err != nil {
		var e store.ConversationEntity
		return e, fmt.Errorf("failed to get conversation entity from sqlite db: %v", err)
	}

	return e, nil
}

func (s ConversationStore) UpdateConversationEntity(
	ctx context.Context, e store.ConversationEntity,
) error {
	query := "UPDATE [conversation] SET updated = ?, data = ? WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, e.UpdatedAt, e.EncryptedData, e.Id)
	if err != nil {
		return fmt.Errorf("failed to store updated conversation entity in sqlite db: %v", err)
	}

	return nil
}

func (s ConversationStore) RemoveConversationEntity(ctx context.Context, id model.Uuid) error {
	query := "DELETE FROM [conversation] WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to remove conversation with id %s from sqlite db: %v", id, err)
	}

	return nil
}

func (s ConversationStore) ListConversationEntities(
	ctx context.Context,
) ([]store.ConversationEntity, error) {
	query := "SELECT * FROM [conversation]"
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation list from sqlite db: %v", err)
	}
	defer rows.Close()

	cs := make([]store.ConversationEntity, 0)
	for rows.Next() {
		var e store.ConversationEntity
		if err := rows.Scan(&e.Id, &e.CreatedAt, &e.UpdatedAt, &e.EncryptedData); err != nil {
			return nil, fmt.Errorf("failed to scan conversation row: %v", err)
		}

		cs = append(cs, e)
	}

	return cs, nil
}

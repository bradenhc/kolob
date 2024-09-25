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

func (s ConversationStore) AddConversationData(ctx context.Context, m store.ConversationMetadata, d []byte) error {
	slog.Info("Adding conversation information to SQLite database")
	_, err := s.db.ExecContext(
		ctx,
		"INSERT INTO [conversation] VALUES (?, ?, ?, ?)",
		m.Id, m.CreatedAt, m.UpdatedAt, d,
	)
	if err != nil {
		return fmt.Errorf("failed to store new conversation in database: %v", err)
	}

	return nil

}

func (s ConversationStore) GetConversationData(ctx context.Context, id model.Uuid) (store.ConversationMetadata, []byte, error) {
	var m store.ConversationMetadata
	var d []byte
	err := s.db.QueryRowContext(
		ctx, "SELECT id, created, updated, data FROM [conversation] WHERE id = ?", id,
	).Scan(
		&m.Id, &m.CreatedAt, &m.UpdatedAt, &d,
	)
	if err != nil {
		return m, d, fmt.Errorf("failed to get conversation data from SQLite database: %v", err)
	}

	return m, d, nil
}

func (s ConversationStore) UpdateConversationData(ctx context.Context, m store.ConversationMetadata, d []byte) error {
	query := "UPDATE [conversation] SET updated = ?, data = ? WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, m.UpdatedAt, d, m.Id)
	if err != nil {
		return fmt.Errorf("failed to store updated conversation data in database: %v", err)
	}

	return nil
}

func (s ConversationStore) RemoveConversationData(ctx context.Context, id model.Uuid) error {
	query := "DELETE FROM [conversation] WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to remove conversation with id %s: %v", id, err)
	}

	return nil
}

func (s ConversationStore) ListConversationData(ctx context.Context) ([][]byte, error) {
	query := "SELECT data FROM [conversation]"
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation list from database: %v", err)
	}
	defer rows.Close()

	cs := make([][]byte, 0)
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, fmt.Errorf("failed to scan conversation row: %v", err)
		}

		cs = append(cs, data)
	}

	return cs, nil
}

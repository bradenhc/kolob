// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/bradenhc/kolob/internal/model"
	"github.com/bradenhc/kolob/internal/store"
)

type MessageStore struct {
	db *sql.DB
}

func NewMessageStore(db *sql.DB) (MessageStore, error) {
	slog.Info("Setting up table: message")
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS message (
			id				TEXT,
			conversation	TEXT,
			author			TEXT,
			created			INTEGER,
			updated			INTEGER,
			data			BLOB,

			PRIMARY KEY (id),
			FOREIGN KEY (conversation) 	REFERENCES conversation(id) ON DELETE CASCADE,
			FOREIGN KEY (author) 		REFERENCES member(id) 		ON DELETE SET NULL
		)
	`)
	if err != nil {
		var s MessageStore
		return s, fmt.Errorf("failed to create message table: %v", err)
	}

	return MessageStore{db}, nil
}

func (s MessageStore) AddMessageData(
	ctx context.Context, m store.MessageMetadata, d []byte,
) error {
	_, err := s.db.ExecContext(
		ctx,
		"INSERT INTO message VALUES (?, ?, ?, ?, ?, ?)",
		m.Id, m.Conversation, m.Author, m.CreatedAt, m.UpdatedAt, d,
	)
	if err != nil {
		return fmt.Errorf("failed to store new message in database: %v", err)
	}

	return nil
}

func (s MessageStore) GetMessageData(
	ctx context.Context, id model.Uuid,
) (store.MessageMetadata, []byte, error) {
	var m store.MessageMetadata
	var data []byte
	query := `
		SELECT id, author, conversation, content, created, updated, data
		FROM message
		WHERE id = ?
	`
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&m.Id, &m.Author, &m.Conversation, &m.CreatedAt, &m.UpdatedAt, &data,
	)
	if err != nil {
		return m, data, fmt.Errorf("failed to get message from database: %v", err)
	}

	return m, data, nil
}

func (s MessageStore) UpdateMessageData(
	ctx context.Context, m store.MessageMetadata, d []byte,
) error {
	query := "UPDATE message SET updated = ?, data = ? WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, m.UpdatedAt, d, m.Id)
	if err != nil {
		return fmt.Errorf("failed to update message in database: %v", err)
	}
	return nil
}

func (s MessageStore) RemoveMessageData(ctx context.Context, id model.Uuid) error {
	query := "DELETE FROM message WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to remove message from database: %v", err)
	}
	return nil
}

func (s MessageStore) ListMessageData(
	ctx context.Context, cid model.Uuid, q store.ListMessageDataQuery,
) ([][]byte, error) {
	where := make([]string, 1)
	params := make([]any, 1)
	where = append(where, "converation = ?")
	params = append(params, cid)
	if q.Author != nil {
		where = append(where, "author = ?")
		params = append(params, *q.Author)
	}
	if q.CreatedAfter != nil {
		where = append(where, "created >= ?")
		params = append(params, *q.CreatedAfter)
	}
	if q.CreatedBefore != nil {
		where = append(where, "created <= ?")
		params = append(params, *q.CreatedBefore)
	}
	query := fmt.Sprintf("SELECT data from message WHERE %s", strings.Join(where, " AND "))
	rows, err := s.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch messages from database: %v", err)
	}

	datalist := make([][]byte, 0)
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, fmt.Errorf("failed to scan message row: %v", err)
		}

		datalist = append(datalist, data)
	}

	return datalist, nil
}

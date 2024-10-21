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

func (s MessageStore) AddMessageEntity(ctx context.Context, e store.MessageEntity) error {
	_, err := s.db.ExecContext(
		ctx,
		"INSERT INTO message VALUES (?, ?, ?, ?, ?, ?)",
		e.Id, e.Conversation, e.Author, e.CreatedAt, e.UpdatedAt, e.EncryptedData,
	)
	if err != nil {
		return fmt.Errorf("failed to store new message in database: %v", err)
	}

	return nil
}

func (s MessageStore) GetMessageEntity(
	ctx context.Context, id model.Uuid,
) (store.MessageEntity, error) {
	var e store.MessageEntity
	query := "SELECT * FROM message WHERE id = ?"
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&e.Id, &e.Author, &e.Conversation, &e.CreatedAt, &e.UpdatedAt, &e.EncryptedData,
	)
	if err != nil {
		var e store.MessageEntity
		return e, fmt.Errorf("failed to get message from database: %v", err)
	}

	return e, nil
}

func (s MessageStore) UpdateMessageEntity(
	ctx context.Context, e store.MessageEntity,
) error {
	query := "UPDATE message SET updated = ?, data = ? WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, e.UpdatedAt, e.EncryptedData, e.Id)
	if err != nil {
		return fmt.Errorf("failed to update message in database: %v", err)
	}
	return nil
}

func (s MessageStore) RemoveMessageEntity(ctx context.Context, id model.Uuid) error {
	query := "DELETE FROM message WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to remove message from database: %v", err)
	}
	return nil
}

func (s MessageStore) ListMessageEntities(
	ctx context.Context, cid model.Uuid, q store.ListMessageDataQuery,
) ([]store.MessageEntity, error) {
	where := make([]string, 0)
	params := make([]any, 0)
	where = append(where, "conversation = ?")
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
	query := fmt.Sprintf("SELECT * from [message] WHERE %s", strings.Join(where, " AND "))
	rows, err := s.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch messages from database: %v", err)
	}

	entities := make([]store.MessageEntity, 0)
	for rows.Next() {
		var e store.MessageEntity
		err := rows.Scan(
			&e.Id, &e.Author, &e.Conversation, &e.CreatedAt, &e.UpdatedAt, &e.EncryptedData,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message row: %v", err)
		}

		entities = append(entities, e)
	}

	return entities, nil
}

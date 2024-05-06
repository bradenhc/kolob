// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/bradenhc/kolob/internal/fail"
	"github.com/bradenhc/kolob/internal/model"
)

type ConversationService struct {
	db QueryExecutor
}

func NewConversationService(db *sql.DB) ConversationService {
	return ConversationService{db}
}

func (s *ConversationService) CreateConversation(
	ctx context.Context, p model.CreateConversationParams,
) (model.Conversation, error) {
	// Create the conversation object
	c, err := model.NewConversation(p.Name, p.Moderator)
	if err != nil {
		return fail.Zero[model.Conversation]("failed to create conversation object", err)
	}

	// Format the datetimes as RFC3339-compliant strings for storage in the DB
	created := c.CreatedAt.Format(time.RFC3339)
	updated := c.UpdatedAt.Format(time.RFC3339)

	// Encrypt the conversation info before storing it in the database
	eda, err := NewEncryptedDataAccessor[model.Conversation](ctx, s.db, "conversation", p.PassKey)
	if err != nil {
		return fail.Zero[model.Conversation]("failed to create encrypted data accessor", err)
	}
	data, err := eda.Encrypt(c)
	if err != nil {
		return fail.Zero[model.Conversation]("failed to encrypt conversation info", err)
	}

	// Store the conversation in the database
	_, err = s.db.ExecContext(
		ctx,
		"INSERT INTO conversation VALUES (?, ?, ?, ?)",
		c.Id, created, updated, data,
	)
	if err != nil {
		return fail.Zero[model.Conversation]("failed to store conversation in database", err)
	}

	// Add the moderator
	_, err = s.db.ExecContext(ctx, "INSERT INTO moderates VALUES (?, ?)", c.Moderators[0], c.Id)
	if err != nil {
		return fail.Zero[model.Conversation]("failed to store moderates relation in database", err)
	}

	return c, nil
}

func (s *ConversationService) UpdateConversation(
	ctx context.Context, p model.UpdateConversationParams,
) error {
	return nil
}

func (s *ConversationService) RemoveConversation(
	ctx context.Context,
	p model.RemoveConversationParams,
) error {
	return nil
}

func (s *ConversationService) ListConversations(
	ctx context.Context, p model.ListConversationsParams,
) ([]model.Conversation, error) {
	return nil, nil
}

func (s *ConversationService) FindConversationById(
	ctx context.Context, p model.FindConversationByIdParams,
) (model.Conversation, error) {
	var m model.Conversation
	return m, nil
}

// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package sqlite

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/bradenhc/kolob/internal/fail"
	"github.com/bradenhc/kolob/internal/model"
)

type ConversationService struct {
	db *sql.DB
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
	eda, err := NewEncryptedDataAccessor[model.Conversation](ctx, s.db, "conversation", p.PassKey)
	if err != nil {
		return fail.Format("failed to create encrypted data accessor", err)
	}

	c, err := eda.Get(ctx, s.db, p.Id)
	if err != nil {
		return fail.Format("failed to get encrypted conversation", err)
	}

	if p.Name != nil {
		c.Name = *p.Name
	}

	c.UpdatedAt = time.Now()

	err = eda.Set(ctx, s.db, p.Id, c)
	if err != nil {
		return fail.Format("failed to set updated conversation information", err)
	}

	return nil
}

func (s *ConversationService) RemoveConversation(
	ctx context.Context,
	p model.RemoveConversationParams,
) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM conversation WHERE id = ?", p.Id)
	if err != nil {
		return fail.Format("failed to remove member from database", err)
	}

	return nil
}

func (s *ConversationService) ListConversations(
	ctx context.Context, p model.ListConversationsParams,
) ([]model.Conversation, error) {
	eda, err := NewEncryptedDataAccessor[model.Conversation](ctx, s.db, "conversation", p.PassKey)
	if err != nil {
		return nil, fail.Format("failed to create encrypted data accessor", err)
	}

	cs, err := eda.GetList(ctx, s.db)
	if err != nil {
		return nil, fail.Format("failed to get conversation list", err)
	}

	return cs, nil
}

func (s *ConversationService) FindConversationById(
	ctx context.Context, p model.FindConversationByIdParams,
) (model.Conversation, error) {
	eda, err := NewEncryptedDataAccessor[model.Conversation](ctx, s.db, "conversation", p.PassKey)
	if err != nil {
		return fail.Zero[model.Conversation]("failed to create encrypted data accessor", err)
	}

	c, err := eda.Get(ctx, s.db, p.Id)
	if err != nil {
		return fail.Zero[model.Conversation]("failed to get encrypted conversation", err)
	}

	return c, nil
}

func (s *ConversationService) AddConversationMods(
	ctx context.Context, p model.AddConversationModsParams,
) error {
	vals := make([]string, 0)
	for range len(p.Moderators) {
		vals = append(vals, "(?, ?)")
	}
	query := "INSERT OR IGNORE INTO moderates VALUES " + strings.Join(vals, ", ")
	ids := make([]any, 0)
	for _, id := range p.Moderators {
		ids = append(ids, id, p.Id)
	}
	_, err := s.db.ExecContext(ctx, query, ids...)
	if err != nil {
		return fail.Format("failed to add moderator entries in database", err)
	}

	return nil
}

func (s *ConversationService) ListConversationMods(
	ctx context.Context, p model.ListConversationModsParams,
) ([]model.Uuid, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT mid FROM moderates WHERE cid = ?", p.Id)
	if err != nil {
		return nil, fail.Format("failed to get conversation moderators from database", err)
	}

	ids := make([]model.Uuid, 0)
	for rows.Next() {
		var id string
		rows.Scan(&id)
		ids = append(ids, model.Uuid(id))
	}

	return ids, nil
}

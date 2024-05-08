// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package sqlite

import (
	"context"
	"regexp"
	"time"

	"github.com/bradenhc/kolob/internal/fail"
	"github.com/bradenhc/kolob/internal/model"
)

type MessageService struct {
	db QueryExecutor
}

func NewMessageService(db QueryExecutor) MessageService {
	return MessageService{db}
}

func (s *MessageService) CreateMessage(
	ctx context.Context, p model.CreateMessageParams,
) (model.Message, error) {
	// Create the new message object
	m, err := model.NewMessage(p.Author, p.Content)
	if err != nil {
		return fail.Zero[model.Message]("failed to create new message", err)
	}

	// Format the datetimes as RFC3339-compliant strings for storage in the DB
	created := m.CreatedAt.Format(time.RFC3339)
	updated := m.UpdatedAt.Format(time.RFC3339)

	// Encrypt the member info before storing it in the database
	eda, err := NewEncryptedDataAccessor[model.Message](ctx, s.db, "member", p.PassKey)
	if err != nil {
		return fail.Zero[model.Message]("failed to create encrypted data accessor", err)
	}
	data, err := eda.Encrypt(m)
	if err != nil {
		return fail.Zero[model.Message]("failed to encrypt message info", err)
	}

	// Store the message in the database
	_, err = s.db.ExecContext(
		ctx,
		"INSERT INTO message VALUES (?, ?, ?, ?, ?, ?)",
		m.Id, p.ConversationId, m.Author, created, updated, data,
	)
	if err != nil {
		return fail.Zero[model.Message]("failed to store new member in database", err)
	}

	return m, nil
}

func (s *MessageService) GetMessage(
	ctx context.Context, p model.GetMessageParams,
) (model.Message, error) {
	eda, err := NewEncryptedDataAccessor[model.Message](ctx, s.db, "message", p.PassKey)
	if err != nil {
		return fail.Zero[model.Message]("failed to create encrypted data accessor", err)
	}

	return eda.Get(ctx, p.Id)
}

func (s *MessageService) UpdateMessage(ctx context.Context, p model.UpdateMessageParams) error {
	// Get the original message object
	eda, err := NewEncryptedDataAccessor[model.Message](ctx, s.db, "member", p.PassKey)
	if err != nil {
		return fail.Format("failed to create encrypted data accessor", err)
	}

	m, err := eda.Get(ctx, p.Id)
	if err != nil {
		return fail.Format("failed to get original message prior to update", err)
	}

	// Apply updates
	if p.Content != nil {
		m.Content = *p.Content
	}

	m.UpdatedAt = time.Now()

	// Update the entry in the database
	err = eda.Set(ctx, m.Id, m)
	if err != nil {
		return fail.Format("failed to store updated message in database", err)
	}

	return nil
}

func (s *MessageService) RemoveMessage(ctx context.Context, p model.RemoveMessageParams) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM message WHERE id = ?", p.Id)
	if err != nil {
		return fail.Format("failed to remove message from database", err)
	}
	return nil
}

func (s *MessageService) ListMessages(
	ctx context.Context, p model.ListMessagesParams,
) ([]model.Message, error) {
	eda, err := NewEncryptedDataAccessor[model.Message](ctx, s.db, "message", p.PassKey)
	if err != nil {
		return nil, fail.Format("failed to get encrypted data accessor", err)
	}

	pred := make(map[string]any)
	pred["conversation = ?"] = p.ConversationId
	if p.StartDate != nil {
		pred["created >= ?"] = p.StartDate.Format(time.RFC3339)
	}
	if p.EndDate != nil {
		pred["created <= ?"] = p.EndDate.Format(time.RFC3339)
	}
	mlist, err := eda.GetListFilt(ctx, pred)
	if err != nil {
		return nil, fail.Format("failed to get filtered list", err)
	}

	if p.ContentPattern == nil {
		return mlist, err
	}

	// We have a content pattern, so we need to filter the list further
	r, err := regexp.Compile(*p.ContentPattern)
	if err != nil {
		return nil, fail.Format("failed to match content", err)
	}

	flist := make([]model.Message, 0)
	for _, v := range mlist {
		if r.MatchString(v.Content) {
			flist = append(flist, v)
		}
	}

	return flist, nil
}

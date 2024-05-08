// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package model

import (
	"context"
	"fmt"
	"time"

	"github.com/bradenhc/kolob/internal/crypto"
)

type Message struct {
	Id        Uuid      `json:"id"`
	Author    Uuid      `json:"author"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewMessage(author Uuid, content string) (Message, error) {
	var m Message

	uuid, err := NewUuid()
	if err != nil {
		return m, fmt.Errorf("failed to create new message: %v", err)
	}

	now := time.Now()

	m.Id = uuid
	m.Author = author
	m.Content = content
	m.CreatedAt = now
	m.UpdatedAt = now

	return m, nil
}

func (a *Message) Equal(b *Message) bool {
	if a != b {
		if a == nil || b == nil ||
			a.Id != b.Id ||
			a.Author != b.Author ||
			a.Content != b.Content ||
			!a.CreatedAt.Equal(b.CreatedAt) ||
			!a.UpdatedAt.Equal(b.UpdatedAt) {
			return false
		}
	}
	return true
}

type MessageService interface {
	CreateMessage(ctx context.Context, p CreateMessageParams) (Message, error)
	GetMessage(ctx context.Context, p GetMessageParams) (Message, error)
	UpdateMessage(ctx context.Context, p UpdateMessageParams) error
	RemoveMessage(ctx context.Context, p RemoveMessageParams) error
	ListMessages(ctx context.Context, p ListMessagesParams) ([]Message, error)
}

type CreateMessageParams struct {
	ConversationId Uuid       `json:"conversation"`
	Author         Uuid       `json:"author"`
	Content        string     `json:"content"`
	PassKey        crypto.Key `json:"-"`
}

type GetMessageParams struct {
	Id      Uuid       `json:"id"`
	PassKey crypto.Key `json:"-"`
}

type UpdateMessageParams struct {
	Id      Uuid       `json:"id"`
	Content *string    `json:"content"`
	PassKey crypto.Key `json:"-"`
}

type RemoveMessageParams struct {
	Id Uuid `json:"id"`
}

type ListMessagesParams struct {
	ConversationId Uuid       `json:"conversation"`
	ContentPattern *string    `json:"pattern"`
	StartDate      *time.Time `json:"start"`
	EndDate        *time.Time `json:"end"`
	PassKey        crypto.Key `json:"-"`
}

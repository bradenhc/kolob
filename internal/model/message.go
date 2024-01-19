package model

import (
	"context"
	"fmt"
	"time"
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

	m.Id = uuid
	m.Author = author
	m.Content = content
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()

	return m, nil
}

type MessageService interface {
	CreateMessage(ctx context.Context, p CreateMessageParams) (Message, error)
	UpdateMessage(ctx context.Context, p UpdateMessageParams) error
	RemoveMessage(ctx context.Context, p RemoveMessageParams) error
	ListMessages(ctx context.Context, p ListMessagesParams)
}

type CreateMessageParams struct {
	GroupId        Uuid   `json:"group"`
	ConversationId Uuid   `json:"conversation"`
	Author         Uuid   `json:"author"`
	Content        string `json:"content"`
}

type UpdateMessageParams struct {
	GroupId        Uuid    `json:"group"`
	ConversationId Uuid    `json:"conversation"`
	Id             Uuid    `json:"id"`
	Content        *string `json:"content"`
}

type RemoveMessageParams struct {
	GroupId        Uuid `json:"group"`
	ConversationId Uuid `json:"conversation"`
	Id             Uuid `json:"id"`
}

type ListMessagesParams struct {
	GroupId        Uuid    `json:"group"`
	ConversationId Uuid    `json:"conversation"`
	Pattern        *string `json:"pattern"`
}

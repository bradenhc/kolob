package model

import (
	"context"
	"fmt"
	"time"
)

type Conversation struct {
	Id        Uuid      `json:"id"`
	Name      string    `json:"name"`
	Moderator Uuid      `json:"moderator"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewConversation(name string, moderator Uuid) (Conversation, error) {
	var c Conversation

	uuid, err := NewUuid()
	if err != nil {
		return c, fmt.Errorf("failed to create conversation: %v", err)
	}

	c.Id = uuid
	c.Name = name
	c.Moderator = moderator
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	return c, nil
}

type ConversationService interface {
	CreateConversation(ctx context.Context, p CreateConversationParams) (Conversation, error)
	UpdateConversation(ctx context.Context, p UpdateConversationParams) error
	RemoveConversation(ctx context.Context, p RemoveConversationParams) error
	ListConversations(ctx context.Context, p ListConversationsParams) ([]Conversation, error)
	FindConversationById(ctx context.Context, p FindConversationByIdParams) (Conversation, error)
}

type CreateConversationParams struct {
	GroupId   Uuid   `json:"group"`
	Name      string `json:"name"`
	Moderator Uuid   `json:"moderator"`
}

type UpdateConversationParams struct {
	GroupId   Uuid    `json:"group"`
	Id        Uuid    `json:"id"`
	Name      *string `json:"name"`
	Moderator *Uuid   `json:"moderator"`
}

type RemoveConversationParams struct {
	GroupId Uuid `json:"group"`
	Id      Uuid `json:"id"`
}

type ListConversationsParams struct {
	GroupId Uuid    `json:"group"`
	Pattern *string `json:"pattern"`
}

type FindConversationByIdParams struct {
	GroupId Uuid `json:"group"`
	Id      Uuid `json:"id"`
}

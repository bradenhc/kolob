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

type Conversation struct {
	Id         Uuid      `json:"id"`
	Name       string    `json:"name"`
	Moderators []Uuid    `json:"moderators"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func NewConversation(name string, moderator Uuid) (Conversation, error) {
	var c Conversation

	uuid, err := NewUuid()
	if err != nil {
		return c, fmt.Errorf("failed to create conversation: %v", err)
	}

	now := time.Now()

	c.Id = uuid
	c.Name = name
	c.Moderators = []Uuid{moderator}
	c.CreatedAt = now
	c.UpdatedAt = now

	return c, nil
}

type ConversationService interface {
	CreateConversation(ctx context.Context, p CreateConversationParams) (Conversation, error)
	UpdateConversation(ctx context.Context, p UpdateConversationParams) error
	RemoveConversation(ctx context.Context, p RemoveConversationParams) error
	ListConversations(ctx context.Context, p ListConversationsParams) ([]Conversation, error)
	FindConversationById(ctx context.Context, p FindConversationByIdParams) (Conversation, error)
	AddConversationMods(ctx context.Context, p AddConversationModsParams) error
	ListConversationMods(ctx context.Context, p ListConversationModsParams) ([]Uuid, error)
}

type CreateConversationParams struct {
	Name      string     `json:"name"`
	Moderator Uuid       `json:"moderator"`
	PassKey   crypto.Key `json:"-"`
}

type UpdateConversationParams struct {
	Id      Uuid       `json:"id"`
	Name    *string    `json:"name"`
	PassKey crypto.Key `json:"-"`
}

type RemoveConversationParams struct {
	Id Uuid `json:"id"`
}

type ListConversationsParams struct {
	Pattern *string    `json:"pattern"`
	PassKey crypto.Key `json:"-"`
}

type FindConversationByIdParams struct {
	Id      Uuid       `json:"id"`
	PassKey crypto.Key `json:"-"`
}

type AddConversationModsParams struct {
	Id         Uuid       `json:"id"`
	Moderators []Uuid     `json:"moderators"`
	PassKey    crypto.Key `json:"-"`
}

type ListConversationModsParams struct {
	Id      Uuid       `json:"id"`
	PassKey crypto.Key `json:"-"`
}

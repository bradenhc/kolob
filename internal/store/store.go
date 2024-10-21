// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package store

import (
	"context"

	"github.com/bradenhc/kolob/internal/crypto"
	"github.com/bradenhc/kolob/internal/model"
)

type GroupStore interface {
	AddGroupEntity(ctx context.Context, e GroupEntity) error
	IsGroupDataSet(ctx context.Context) (bool, error)
	GetGroupEntity(ctx context.Context) (GroupEntity, error)
	UpdateGroupEntity(ctx context.Context, e GroupEntity) error
}

type MemberStore interface {
	AddMemberEntity(ctx context.Context, e MemberEntity) error
	GetMemberEntity(ctx context.Context, id model.Uuid) (MemberEntity, error)
	GetMemberEntityByUname(ctx context.Context, uhash crypto.DataHash) (MemberEntity, error)
	UpdateMemberEntity(ctx context.Context, e MemberEntity) error
	RemoveMemberEntity(ctx context.Context, id model.Uuid) error
	ListMemberEntities(ctx context.Context) ([]MemberEntity, error)
}

type ConversationStore interface {
	AddConversationEntity(ctx context.Context, e ConversationEntity) error
	GetConversationEntity(ctx context.Context, id model.Uuid) (ConversationEntity, error)
	UpdateConversationEntity(ctx context.Context, e ConversationEntity) error
	RemoveConversationEntity(ctx context.Context, id model.Uuid) error
	ListConversationEntities(ctx context.Context) ([]ConversationEntity, error)
}

type MessageStore interface {
	AddMessageEntity(ctx context.Context, m MessageEntity) error
	GetMessageEntity(ctx context.Context, id model.Uuid) (MessageEntity, error)
	UpdateMessageEntity(ctx context.Context, m MessageEntity) error
	RemoveMessageEntity(ctx context.Context, id model.Uuid) error
	ListMessageEntities(ctx context.Context, cid model.Uuid, q ListMessageDataQuery) ([]MessageEntity, error)
}

type ListMessageDataQuery struct {
	Author        *model.Uuid
	CreatedAfter  *int64
	CreatedBefore *int64
}

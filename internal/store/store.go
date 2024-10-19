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
	AddMessageData(ctx context.Context, m MessageMetadata, d []byte) error
	GetMessageData(ctx context.Context, id model.Uuid) (MessageMetadata, []byte, error)
	UpdateMessageData(ctx context.Context, m MessageMetadata, d []byte) error
	RemoveMessageData(ctx context.Context, id model.Uuid) error
	ListMessageData(ctx context.Context, cid model.Uuid, q ListMessageDataQuery) ([][]byte, error)
}

type MessageMetadata struct {
	Id           model.Uuid
	Author       model.Uuid
	Conversation model.Uuid
	CreatedAt    int64
	UpdatedAt    int64
}

type ListMessageDataQuery struct {
	Author        *model.Uuid
	CreatedAfter  *int64
	CreatedBefore *int64
}

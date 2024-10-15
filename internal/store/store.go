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
	SetGroupData(ctx context.Context, m GroupMetadata, d []byte) error
	IsGroupDataSet(ctx context.Context) (bool, error)
	GetGroupMetadata(ctx context.Context) (GroupMetadata, error)
	GetGroupData(ctx context.Context) ([]byte, error)
	UpdateGroupData(ctx context.Context, m GroupMetadata, d []byte) error
}

type GroupMetadata struct {
	Id           model.Uuid
	GroupHash    crypto.DataHash
	PassSalt     crypto.Salt
	PassHash     crypto.PassHash
	EncryptedKey []byte
	CreatedAt    int64
	UpdatedAt    int64
}

type MemberStore interface {
	AddMemberData(ctx context.Context, m MemberMetadata, d []byte) error
	GetMemberData(ctx context.Context, id model.Uuid) (MemberMetadata, []byte, error)
	GetMemberDataByUname(ctx context.Context, uhash crypto.DataHash) (MemberMetadata, []byte, error)
	UpdateMemberData(ctx context.Context, m MemberMetadata, d []byte) error
	RemoveMemberData(ctx context.Context, id model.Uuid) error
	ListMemberData(ctx context.Context) ([][]byte, error)
}

type MemberMetadata struct {
	Id           model.Uuid
	UsernameHash crypto.DataHash
	PassHash     crypto.PassHash
	CreatedAt    int64
	UpdatedAt    int64
}

type ConversationStore interface {
	AddConversationData(ctx context.Context, m ConversationMetadata, d []byte) error
	GetConversationData(ctx context.Context, id model.Uuid) (ConversationMetadata, []byte, error)
	UpdateConversationData(ctx context.Context, m ConversationMetadata, d []byte) error
	RemoveConversationData(ctx context.Context, id model.Uuid) error
	ListConversationData(ctx context.Context) ([][]byte, error)
}

type ConversationMetadata struct {
	Id        model.Uuid
	CreatedAt int64
	UpdatedAt int64
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

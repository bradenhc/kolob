// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package model

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/bradenhc/kolob/internal/crypto"
	"github.com/bradenhc/kolob/internal/types"
	flatbuffers "github.com/google/flatbuffers/go"
)

func NewMember(username, name string) (*types.Member, error) {
	uuid, err := NewUuid()
	if err != nil {
		return nil, fmt.Errorf("failed to create new member: %v", err)
	}

	builder := flatbuffers.NewBuilder(256)

	mi := builder.CreateString(string(uuid))
	mu := builder.CreateString(username)
	mn := builder.CreateString(name)

	now := time.Now()

	types.MemberStart(builder)
	types.MemberAddId(builder, mi)
	types.MemberAddUname(builder, mu)
	types.MemberAddName(builder, mn)
	types.MemberAddCreated(builder, now.UnixMilli())
	types.MemberAddUpdated(builder, now.UnixMilli())

	m := types.MemberEnd(builder)

	builder.Finish(m)

	return types.GetRootAsMember(builder.FinishedBytes(), 0), nil
}

func MemberEqual(a, b *types.Member) bool {
	if a != b {
		if a == nil || b == nil ||
			!slices.Equal(a.Id(), b.Id()) ||
			!slices.Equal(a.Uname(), b.Uname()) ||
			!slices.Equal(a.Name(), b.Name()) ||
			a.Created() != b.Created() ||
			a.Updated() != b.Updated() {
			return false
		}
	}
	return true
}

type MemberService interface {
	CreateMember(ctx context.Context, p CreateMemberParams) (*types.Member, error)
	AuthenticateMember(ctx context.Context, p AuthenticateMemberParams) error
	UpdateMember(ctx context.Context, p UpdateMemberParams) error
	RemoveMember(ctx context.Context, p RemoveMemberParams) error
	ListMembers(ctx context.Context, p ListMembersParams) ([]*types.Member, error)
	FindMemberByUsername(ctx context.Context, p FindMemberByUsernameParams) (*types.Member, error)
}

type CreateMemberParams struct {
	Username string          `json:"username"`
	Name     string          `json:"name"`
	Password crypto.Password `json:"pass"`
	PassKey  crypto.Key      `json:"-"`
}

type AuthenticateMemberParams struct {
	Username string          `json:"username"`
	Password crypto.Password `json:"password"`
}

type UpdateMemberParams struct {
	Id       Uuid       `json:"id"`
	Username *string    `json:"username"`
	Name     *string    `json:"name"`
	PassKey  crypto.Key `json:"-"`
}

type RemoveMemberParams struct {
	Id Uuid `json:"id"`
}

type ListMembersParams struct {
	PassKey crypto.Key `json:"-"`
}

type FindMemberByUsernameParams struct {
	Username string     `json:"username"`
	PassKey  crypto.Key `json:"-"`
}

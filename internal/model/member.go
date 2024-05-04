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

type Member struct {
	Id        Uuid      `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewMember(username, name string) (Member, error) {
	var m Member
	uuid, err := NewUuid()
	if err != nil {
		return m, fmt.Errorf("failed to create new member: %v", err)
	}

	m.Id = uuid
	m.Username = username
	m.Name = name
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()

	return m, nil
}

func (a *Member) Equal(b *Member) bool {
	if a != b {
		if a == nil || b == nil ||
			a.Id != b.Id ||
			a.Username != b.Username ||
			a.Name != b.Name ||
			!a.CreatedAt.Equal(b.CreatedAt) ||
			!a.UpdatedAt.Equal(b.UpdatedAt) {
			return false
		}
	}
	return true
}

type MemberService interface {
	CreateMember(ctx context.Context, p CreateMemberParams) (Member, error)
	AuthenticateMember(ctx context.Context, p AuthenticateMemberParams) error
	UpdateMember(ctx context.Context, p UpdateMemberParams) error
	RemoveMember(ctx context.Context, p RemoveMemberParams) error
	ListMembers(ctx context.Context, p ListMembersParams) ([]Member, error)
	FindMemberByUsername(ctx context.Context, p FindMemberByUsernameParams) (Member, error)
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

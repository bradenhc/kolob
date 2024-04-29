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
	Key      crypto.Key      `json:"-"`
}

type AuthenticateMemberParams struct {
	Username string          `json:"username"`
	Password crypto.Password `json:"password"`
	Key      crypto.Key      `json:"-"`
}

type UpdateMemberParams struct {
	Id       Uuid       `json:"id"`
	Username *string    `json:"username"`
	Name     *string    `json:"name"`
	Key      crypto.Key `json:"-"`
}

type RemoveMemberParams struct {
	Id Uuid `json:"id"`
}

type ListMembersParams struct {
	NamePattern *string    `json:"pattern"`
	Key         crypto.Key `json:"-"`
}

type FindMemberByUsernameParams struct {
	Username string     `json:"username"`
	Key      crypto.Key `json:"-"`
}

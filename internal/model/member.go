// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package model

import (
	"context"
	"fmt"
	"time"
)

type Member struct {
	Id        Uuid      `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	IsOwner   bool      `json:"isOwner"`
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
	UpdateMember(ctx context.Context, p UpdateMemberParams) error
	RemoveMember(ctx context.Context, p RemoveMemberParams) error
	ListMembers(ctx context.Context, p ListMembersParams) ([]Member, error)
	FindMemberByUsername(ctx context.Context, p FindMemberByUsernameParams) (Member, error)
}

type CreateMemberParams struct {
	Username string `json:"username"`
	Name     string `json:"name"`
}

type UpdateMemberParams struct {
	Id       Uuid    `json:"id"`
	Username *string `json:"username"`
	Name     *string `json:"name"`
}

type RemoveMemberParams struct {
	Id Uuid `json:"id"`
}

type ListMembersParams struct {
	NamePattern *string `json:"pattern"`
}

type FindMemberByUsernameParams struct {
	Username string `json:"username"`
}

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

type Group struct {
	Id          Uuid      `json:"id"`
	GroupId     string    `json:"gid"`
	Name        string    `json:"name"`
	Description string    `json:"desc"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func NewGroup(gid, name, desc string) (Group, error) {
	id, err := NewUuid()
	if err != nil {
		var g Group
		return g, fmt.Errorf("failed to create UUID for new group: %v", err)
	}

	now := time.Now()

	g := Group{
		Id:          id,
		GroupId:     gid,
		Name:        name,
		Description: desc,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return g, nil
}

func (a *Group) Equal(b *Group) bool {
	if a != b {
		if a == nil || b == nil ||
			a.Id != b.Id ||
			a.GroupId != b.GroupId ||
			a.Name != b.Name ||
			a.Description != b.Description ||
			!a.CreatedAt.Equal(b.CreatedAt) ||
			!a.UpdatedAt.Equal(b.UpdatedAt) {
			return false
		}
	}

	return true
}

type GroupService interface {
	InitGroup(ctx context.Context, p InitGroupParams) (Group, error)
	GetGroupInfo(ctx context.Context, p GetGroupInfoParams) (Group, error)
	GetGroupPassKey(ctx context.Context, p GetGroupPassKeyParams) (crypto.Key, error)
	GetGroupDataKey(ctx context.Context, p GetGroupDataKeyParams) (crypto.Key, error)
	AuthenticateGroup(ctx context.Context, p AuthenticateGroupParams) error
	UpdateGroup(ctx context.Context, p UpdateGroupParams) error
	ChangeGroupPassword(ctx context.Context, p ChangeGroupPasswordParams) error
}

type InitGroupParams struct {
	GroupId     string          `json:"gid"`
	Name        string          `json:"name"`
	Description string          `json:"desc"`
	Password    crypto.Password `json:"pass"`
}

type GetGroupInfoParams struct {
	Id      Uuid
	PassKey crypto.Key
}

type AuthenticateGroupParams struct {
	GroupId  string          `json:"gid"`
	Password crypto.Password `json:"pass"`
}

type UpdateGroupParams struct {
	Id          Uuid       `json:"id"`
	GroupId     *string    `json:"gid"`
	Name        *string    `json:"name"`
	Description *string    `json:"desc"`
	PassKey     crypto.Key `json:"-"`
}

type ChangeGroupPasswordParams struct {
	Id      Uuid            `json:"id"`
	OldPass crypto.Password `json:"old"`
	NewPass crypto.Password `json:"new"`
}

type GetGroupPassKeyParams struct {
	Password crypto.Password
}

type GetGroupDataKeyParams struct {
	PassKey crypto.Key
}

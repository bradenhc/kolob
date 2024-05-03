// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package model

import (
	"context"
	"fmt"

	"github.com/bradenhc/kolob/internal/crypto"
)

type Group struct {
	Id          Uuid   `json:"id"`
	GroupId     string `json:"gid"`
	Name        string `json:"name"`
	Description string `json:"desc"`
}

func NewGroup(gid, name, desc string) (Group, error) {
	id, err := NewUuid()
	if err != nil {
		var g Group
		return g, fmt.Errorf("failed to create UUID for new group: %v", err)
	}

	g := Group{
		Id:          id,
		GroupId:     gid,
		Name:        name,
		Description: desc,
	}

	return g, nil
}

type GroupService interface {
	CreateGroup(ctx context.Context, p CreateGroupParams) (Group, error)
	GetGroupInfo(ctx context.Context, p GetGroupInfoParams) (Group, error)
	GetGroupPassKey(ctx context.Context, p GetGroupPassKeyParams) (crypto.Key, error)
	GetGroupDataKey(ctx context.Context, p GetGroupDataKeyParams) (crypto.Key, error)
	AuthenticateGroup(ctx context.Context, p AuthenticateGroupParams) error
	UpdateGroup(ctx context.Context, p UpdateGroupParams) error
	ChangeGroupPassword(ctx context.Context, p ChangeGroupPasswordParams) error
}

type CreateGroupParams struct {
	GroupId     string          `json:"gid"`
	Name        string          `json:"name"`
	Description string          `json:"desc"`
	Password    crypto.Password `json:"pass"`
}

type GetGroupInfoParams struct {
	PassKey crypto.Key
}

type AuthenticateGroupParams struct {
	GroupId  string          `json:"gid"`
	Password crypto.Password `json:"pass"`
}

type UpdateGroupParams struct {
	GroupId     *string    `json:"gid"`
	Name        *string    `json:"name"`
	Description *string    `json:"desc"`
	PassKey     crypto.Key `json:"-"`
}

type ChangeGroupPasswordParams struct {
	OldPass crypto.Password `json:"old"`
	NewPass crypto.Password `json:"new"`
}

type GetGroupPassKeyParams struct {
	Password crypto.Password
}

type GetGroupDataKeyParams struct {
	PassKey crypto.Key
}

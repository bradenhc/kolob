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
	GetGroupName(ctx context.Context, p GetGroupNameParams) (string, error)
	GetGroupDescription(ctx context.Context, p GetGroupDescriptionParams) (string, error)
	GetGroupPassKey(ctx context.Context, p GetGroupPassKeyParams) (crypto.Key, error)
	GetGroupDataKey(ctx context.Context, p GetGroupDataKeyParams) (crypto.Key, error)
	AuthenticateGroup(ctx context.Context, p AuthenticateGroupParams) error
	UpdateGroup(ctx context.Context, p UpdateGroupParams) error
	ChangeGroupPassword(ctx context.Context, p ChangeGroupPasswordParams) error
}

type CreateGroupParams struct {
	GroupId     string
	Name        string
	Description string
	Password    crypto.Password
}

type GetGroupNameParams struct {
	PassKey crypto.Key
}

type GetGroupDescriptionParams struct {
	PassKey crypto.Key
}

type AuthenticateGroupParams struct {
	GroupId  string
	Password crypto.Password
}

type UpdateGroupParams struct {
	GroupId     *string
	Name        *string
	Description *string
	PassKey     crypto.Key
}

type ChangeGroupPasswordParams struct {
	OldPass crypto.Password
	NewPass crypto.Password
}

type GetGroupPassKeyParams struct {
	Password crypto.Password
}

type GetGroupDataKeyParams struct {
	PassKey crypto.Key
}

// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package model

import (
	"context"

	"github.com/bradenhc/kolob/internal/crypto"
)

type GroupService interface {
	CreateGroup(ctx context.Context, params CreateGroupParams) (Uuid, error)
	GetGroupName(ctx context.Context) (string, error)
	GetGroupDescription(ctx context.Context) (string, error)
	GetGroupPassKey(ctx context.Context, params GetGroupPassKeyParams) (crypto.Key, error)
	GetGroupDataKey(ctx context.Context, params GetGroupDataKeyParams) (crypto.Key, error)
	AuthenticateGroup(ctx context.Context, params AuthenticateGroupParams) error
	UpdateGroup(ctx context.Context, params UpdateGroupParams) error
	ChangeGroupPassword(ctx context.Context, params ChangeGroupPasswordParams) error
}

type CreateGroupParams struct {
	GroupId     string
	Name        string
	Description string
	Password    crypto.Password
}

type AuthenticateGroupParams struct {
	GroupId  string
	Password crypto.Password
}

type UpdateGroupParams struct {
	GroupId     string
	Name        string
	Description string
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

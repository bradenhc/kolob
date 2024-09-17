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
	SetGroupData(ctx context.Context, meta GroupMetadata, data []byte) error
	IsGroupDataSet(ctx context.Context) (bool, error)
	GetGroupMetadata(ctx context.Context) (GroupMetadata, error)
	GetGroupData(ctx context.Context) ([]byte, error)
	UpdateGroupData(ctx context.Context, meta GroupMetadata, data []byte) error
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

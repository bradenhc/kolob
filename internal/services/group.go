// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/bradenhc/kolob/internal/crypto"
	"github.com/bradenhc/kolob/internal/model"
	"github.com/bradenhc/kolob/internal/store"
)

type GroupService struct {
	store store.GroupStore
}

func NewGroupService(store store.GroupStore) GroupService {
	return GroupService{store}
}

func (svc GroupService) Create(ctx context.Context, req *GroupInitRequest) (*model.Group, error) {
	// Make sure there isn't already a set of group information in the database
	exists, err := svc.store.IsGroupDataSet(ctx)
	if exists {
		return nil, fmt.Errorf("group already exists")
	}

	// Generate a new random key that will encrypt all data for the group (data key)
	slog.Info("Generating data encryption key for group")
	dkey, err := crypto.NewRandomKey()
	if err != nil {
		return nil, fmt.Errorf("failed to create new data encryption key for group: %v", err)
	}

	// Generate a second key using the user password that will encrypt the data key (pass key)
	psalt, err := crypto.NewSalt()
	if err != nil {
		return nil, fmt.Errorf("failed to create salt for group: %v", err)
	}
	slog.Info("Deriving password key")
	pkey := crypto.NewDerivedKey(crypto.Password(req.Password()), psalt)

	// Encyrypt the data key using the pass key before we store it in the database
	slog.Info("Encrypting data key")
	ekey, err := crypto.Encrypt(pkey, dkey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data key: %v", err)
	}

	// Hash the group id so we can use it for authentication without leaking information
	ghash := crypto.HashData(req.GroupId())

	// Hash the password so we can use it for authentication
	phash, err := crypto.HashPassword(crypto.Password(req.Password()))
	if err != nil {
		return nil, fmt.Errorf("failed to generated group password hash: %v", err)
	}

	// Create the group data to encrypt and store
	group, err := model.NewGroup(string(req.GroupId()), string(req.Name()), string(req.Description()))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize group: %v", err)
	}

	// Encrypt the group information to protect privacy
	slog.Info("Encrypting group information")
	edata, err := crypto.Encrypt(dkey, group.Table().Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt group data before storing in database: %v", err)
	}

	// Store the group information
	meta := store.GroupMetadata{
		Id:           model.Uuid(group.Id()),
		GroupHash:    ghash,
		PassSalt:     psalt,
		PassHash:     phash,
		EncryptedKey: ekey,
		CreatedAt:    group.Created(),
		UpdatedAt:    group.Updated(),
	}
	svc.store.SetGroupData(ctx, meta, edata)

	return group, nil
}

func (g GroupService) GetInfo(ctx context.Context, dkey crypto.Key) (*model.Group, error) {
	d, err := g.store.GetGroupData(ctx)
	if err != nil {
		return nil, fmt.Errorf("store error: %v", err)
	}

	data, err := crypto.Decrypt(dkey, d)
	if err != nil {
		return nil, fmt.Errorf("decrypt error: %v", err)
	}

	return model.GetRootAsGroup(data, 0), nil
}

func (g GroupService) Authenticate(
	ctx context.Context, req *GroupAuthenticateRequest,
) (crypto.Key, error) {
	m, err := g.store.GetGroupMetadata(ctx)
	if err != nil {
		return nil, fmt.Errorf("store error: %v", err)
	}

	pass := crypto.Password(req.Password())
	if !crypto.CheckDataHash(req.GroupId(), m.GroupHash) {
		return nil, fmt.Errorf("incorrect credentials")
	}
	if !crypto.CheckPasswordHash(pass, m.PassHash) {
		return nil, fmt.Errorf("incorrect credentials")

	}

	pkey := crypto.NewDerivedKey(pass, m.PassSalt)
	dkey, err := crypto.Decrypt(pkey, m.EncryptedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt group data key: %v", err)
	}

	return dkey, nil
}

func (g GroupService) Update(
	ctx context.Context, req *GroupUpdateRequest, dkey crypto.Key,
) (*model.Group, error) {
	m, err := g.store.GetGroupMetadata(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get metadata from store: %v", err)
	}

	edata, err := g.store.GetGroupData(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get data from store: %v", err)
	}

	data, err := crypto.Decrypt(dkey, edata)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt group data: %v", err)
	}
	prev := model.GetRootAsGroup(data, 0)

	group := model.GroupCloneWithUpdates(prev, req.GroupId(), req.Name(), req.Description())

	m.GroupHash = crypto.HashData(group.Gid())
	m.UpdatedAt = group.Updated()

	edata, err = crypto.Encrypt(dkey, group.Table().Bytes)
	g.store.UpdateGroupData(ctx, m, edata)

	return group, nil
}

func (g GroupService) ChangePassword(
	ctx context.Context, req *GroupChangePasswordRequest, dkey crypto.Key,
) error {
	m, err := g.store.GetGroupMetadata(ctx)
	if err != nil {
		return fmt.Errorf("metadata fetch failed: %v", err)
	}

	if !crypto.CheckPasswordHash(crypto.Password(req.OldPassword()), m.PassHash) {
		return fmt.Errorf("old password validation failed")
	}

	edata, err := g.store.GetGroupData(ctx)
	if err != nil {
		return fmt.Errorf("data fetch failed: %v", err)
	}

	data, err := crypto.Decrypt(dkey, edata)
	if err != nil {
		return fmt.Errorf("failed to decrypt group data: %v", err)
	}

	psalt, err := crypto.NewSalt()
	if err != nil {
		return fmt.Errorf("failed to generate salt for new password: %v", err)
	}

	npass, err := crypto.NewPassword(string(req.NewPassword()))
	if err != nil {
		return fmt.Errorf("invalid password: %v", err)
	}

	phash, err := crypto.HashPassword(npass)
	if err != nil {
		return fmt.Errorf("failed to generate hash for new password: %v", err)
	}
	m.PassSalt = psalt
	m.PassHash = phash

	pkey := crypto.NewDerivedKey(npass, psalt)
	ekey, err := crypto.Encrypt(pkey, dkey)
	if err != nil {
		return fmt.Errorf("failed to encrypt data key with new password key: %v", err)
	}
	m.EncryptedKey = ekey

	updated := time.Now().UnixMilli()
	m.UpdatedAt = updated
	group := model.GetRootAsGroup(data, 0)
	group.MutateUpdated(updated)

	edata, err = crypto.Encrypt(dkey, data)
	err = g.store.UpdateGroupData(ctx, m, edata)
	if err != nil {
		return fmt.Errorf("failed to update group: %v", err)
	}

	return nil
}

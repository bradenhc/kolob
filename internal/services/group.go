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
	if err != nil {
		return nil, fmt.Errorf("failed to check group initialization: %v", err)
	}
	if exists {
		return nil, fmt.Errorf("group already exists")
	}

	// Generate a new random key that will encrypt all data for the group (data key)
	slog.Info("Generating data encryption key for group")
	dkey, err := crypto.NewRandomKey()
	if err != nil {
		return nil, fmt.Errorf("failed to create new data encryption key for group: %v", err)
	}

	// Create the group data to encrypt and store
	group, err := model.NewGroup(string(req.GroupId()), string(req.Name()), string(req.Description()))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize group: %v", err)
	}

	entity, err := store.NewGroupEntity(group, crypto.Password(req.Password()), dkey)
	if err != nil {
		return nil, fmt.Errorf("failed to create group store entity: %v", err)
	}

	svc.store.AddGroupEntity(ctx, entity)

	return group, nil
}

func (g GroupService) Get(ctx context.Context, dkey crypto.Key) (*model.Group, error) {
	e, err := g.store.GetGroupEntity(ctx)
	if err != nil {
		return nil, fmt.Errorf("store error: %v", err)
	}

	data, err := crypto.Decrypt(dkey, e.EncryptedData)
	if err != nil {
		return nil, fmt.Errorf("decrypt error: %v", err)
	}

	return model.GetRootAsGroup(data, 0), nil
}

func (g GroupService) Authenticate(
	ctx context.Context, req *GroupAuthenticateRequest,
) (crypto.Key, error) {
	m, err := g.store.GetGroupEntity(ctx)
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
	e, err := g.store.GetGroupEntity(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get group entity from store: %v", err)
	}

	group, err := e.Update(dkey, req.GroupId(), req.Name(), req.Description())
	if err != nil {
		return nil, fmt.Errorf("could not update group: %v", err)
	}

	err = g.store.UpdateGroupEntity(ctx, e)
	if err != nil {
		return nil, fmt.Errorf("failed to store updated group entity: %v", err)
	}

	return group, nil
}

func (g GroupService) ChangePassword(
	ctx context.Context, req *GroupChangePasswordRequest, dkey crypto.Key,
) error {
	e, err := g.store.GetGroupEntity(ctx)
	if err != nil {
		return fmt.Errorf("metadata fetch failed: %v", err)
	}

	if !crypto.CheckPasswordHash(crypto.Password(req.OldPassword()), e.PassHash) {
		return fmt.Errorf("old password validation failed")
	}

	data, err := crypto.Decrypt(dkey, e.EncryptedData)
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
	e.PassSalt = psalt
	e.PassHash = phash

	pkey := crypto.NewDerivedKey(npass, psalt)
	ekey, err := crypto.Encrypt(pkey, dkey)
	if err != nil {
		return fmt.Errorf("failed to encrypt data key with new password key: %v", err)
	}
	e.EncryptedKey = ekey

	updated := time.Now().UnixMilli()
	e.UpdatedAt = updated
	group := model.GetRootAsGroup(data, 0)
	group.MutateUpdated(updated)

	e.EncryptedData, err = crypto.Encrypt(dkey, data)
	if err != nil {
		return fmt.Errorf("failed to encrypt update group data: %v", err)
	}

	err = g.store.UpdateGroupEntity(ctx, e)
	if err != nil {
		return fmt.Errorf("failed to update group: %v", err)
	}

	return nil
}

// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/bradenhc/kolob/internal/crypto"
	"github.com/bradenhc/kolob/internal/model"
)

type GroupService struct {
	db *sql.DB
}

func NewGroupService(db *sql.DB) GroupService {
	return GroupService{db}
}

func (s GroupService) CreateGroup(
	ctx context.Context, params model.CreateGroupParams,
) (model.Uuid, error) {
	var id model.Uuid

	// Generate a new random key that will encrypt all data for the group
	randomKey, err := crypto.NewRandomKey()
	if err != nil {
		return id, fmt.Errorf("failed to create new data encryption key for group: %v", err)
	}

	// Generate a second key using the user password that will encrypt the data key
	salt, err := crypto.NewSalt()
	if err != nil {
		return id, fmt.Errorf("failed to create salt for new group: %v", err)
	}

	derivedKey := crypto.NewDerivedKey(params.Password, salt)

	// Encyrypt the data encryption key before we store it in the database
	ekey, err := crypto.Encrypt(derivedKey, randomKey)
	if err != nil {
		return id, fmt.Errorf("failed to encrypt data key: %v", err)
	}

	// Hash the password so we can use it for authentication
	hash, err := crypto.HashPassword(params.Password)
	if err != nil {
		return id, fmt.Errorf("failed to generated group password hash: %v", err)
	}

	// Store everything inside the database
	id, err = model.NewUuid()
	if err != nil {
		return id, fmt.Errorf("failed to create new group UUID: %v", err)
	}

	_, err = s.db.ExecContext(
		ctx,
		"INSERT INTO info VALUES (?, ?, ?, ?, ?, ?, ?)",
		id, params.GroupId, params.Name, params.Description, salt, hash, ekey,
	)
	if err != nil {
		return "", fmt.Errorf("failed to store group information in database: %v", err)
	}

	return id, nil
}

func (s GroupService) GetGroupName(ctx context.Context) (string, error) {
	var name string
	err := s.db.QueryRowContext(ctx, "SELECT name FROM info").Scan(&name)
	if err != nil {
		return "", fmt.Errorf("failed to get group name from database: %v", err)
	}
	return name, nil
}

func (s GroupService) GetGroupDescription(ctx context.Context) (string, error) {
	var desc string
	err := s.db.QueryRowContext(ctx, "SELECT desc FROM info").Scan(&desc)
	if err != nil {
		return "", fmt.Errorf("failed to get group description from database: %v", err)
	}
	return desc, nil
}

func (s GroupService) GetGroupPassKey(
	ctx context.Context, params model.GetGroupPassKeyParams,
) (crypto.Key, error) {
	var salt crypto.Salt
	var hash string
	err := s.db.QueryRowContext(ctx, "SELECT psalt, phash FROM info").Scan(&salt, &hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get group password salt and hash: %v", err)
	}

	if !crypto.CheckPasswordHash(params.Password, hash) {
		return nil, fmt.Errorf("failed to check group password against hash: %v", err)
	}

	key := crypto.NewDerivedKey(params.Password, salt)
	return key, nil
}

func (s GroupService) GetGroupDataKey(
	ctx context.Context, params model.GetGroupDataKeyParams,
) (crypto.Key, error) {
	var ekey []byte
	err := s.db.QueryRowContext(ctx, "SELECT ekey FROM info").Scan(&ekey)
	if err != nil {
		return nil, fmt.Errorf("failed to get encrypted group data key: %v", err)
	}

	dkey, err := crypto.Decrypt(params.PassKey, ekey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt group data key: %v", err)
	}

	return dkey, nil
}

func (s GroupService) AuthenticateGroup(
	ctx context.Context, params model.AuthenticateGroupParams,
) error {
	var hash string
	err := s.db.QueryRowContext(
		ctx,
		"SELECT phash FROM info WHERE gid = ?",
		params.GroupId,
	).Scan(&hash)

	if err != nil {
		return fmt.Errorf("failed to get password hash from database: %v", err)
	}

	if !crypto.CheckPasswordHash(params.Password, hash) {
		return fmt.Errorf("password authentication failed for group: %v", err)
	}

	return nil
}

func (s GroupService) UpdateGroup(ctx context.Context, params model.UpdateGroupParams) error {
	fails := make([]string, 0)

	if params.GroupId != "" {
		_, err := s.db.ExecContext(ctx, "UPDATE info SET gid = ?", params.GroupId)
		if err != nil {
			fails = append(fails, fmt.Sprintf("failed to set group ID: %v", err))
		}
	}

	if params.Name != "" {
		_, err := s.db.ExecContext(ctx, "UPDATE info SET name = ?", params.Name)
		if err != nil {
			fails = append(fails, fmt.Sprintf("failed to set group name: %v", err))
		}
	}

	if params.Description != "" {
		_, err := s.db.ExecContext(ctx, "UPDATE info SET desc = ?", params.Description)
		if err != nil {
			fails = append(fails, fmt.Sprintf("failed to set group description: %v", err))
		}
	}

	if len(fails) > 0 {
		return errors.New(strings.Join(fails, "; "))
	}

	return nil
}

func (s GroupService) ChangeGroupPassword(
	ctx context.Context, params model.ChangeGroupPasswordParams,
) error {
	passParams := model.GetGroupPassKeyParams{
		Password: params.OldPass,
	}
	passKey, err := s.GetGroupPassKey(ctx, passParams)
	if err != nil {
		return fmt.Errorf("failed to get group pass key: %v", err)
	}

	dataParams := model.GetGroupDataKeyParams{
		PassKey: passKey,
	}
	dataKey, err := s.GetGroupDataKey(ctx, dataParams)
	if err != nil {
		return fmt.Errorf("failed to get group data key: %v", err)
	}

	salt, err := crypto.NewSalt()
	if err != nil {
		return fmt.Errorf("failed to generate salt for new password: %v", err)
	}

	hash, err := crypto.HashPassword(params.NewPass)
	if err != nil {
		return fmt.Errorf("failed to generate hash for new password: %v", err)
	}

	passKey = crypto.NewDerivedKey(params.NewPass, salt)

	ekey, err := crypto.Encrypt(passKey, dataKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt data key with new password key: %v", err)
	}

	_, err = s.db.ExecContext(
		ctx,
		"UPDATE info SET psalt = ?, phash = ?, ekey = ?",
		salt, hash, ekey,
	)
	if err != nil {
		return fmt.Errorf("failed to store new password information in database: %v", err)
	}

	return nil
}

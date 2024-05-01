// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package sqlite

import (
	"context"
	"database/sql"
	"fmt"

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
) (model.Group, error) {
	g, err := model.NewGroup(params.GroupId, params.Name, params.Description)
	if err != nil {
		var g model.Group
		return g, fmt.Errorf("failed to create new group: %v", err)
	}

	// Generate a new random key that will encrypt all data for the group (data key)
	dkey, err := crypto.NewRandomKey()
	if err != nil {
		var g model.Group
		return g, fmt.Errorf("failed to create new data encryption key for group: %v", err)
	}

	// Generate a second key using the user password that will encrypt the data key (pass key)
	psalt, err := crypto.NewSalt()
	if err != nil {
		var g model.Group
		return g, fmt.Errorf("failed to create salt for new group: %v", err)
	}
	pkey := crypto.NewDerivedKey(params.Password, psalt)

	// Encyrypt the data key using the pass key before we store it in the database
	ekey, err := crypto.Encrypt(pkey, dkey)
	if err != nil {
		var g model.Group
		return g, fmt.Errorf("failed to encrypt data key: %v", err)
	}

	// Hash the group id so we can use it for authentication without leaking information
	ghash := crypto.HashData([]byte(g.GroupId))

	// Hash the password so we can use it for authentication
	phash, err := crypto.HashPassword(params.Password)
	if err != nil {
		var g model.Group
		return g, fmt.Errorf("failed to generated group password hash: %v", err)
	}

	// Encrypt the group information to protect privacy
	a := crypto.NewAgent[model.Group](dkey)
	data, err := a.Encrypt(g)
	if err != nil {
		var g model.Group
		return g, fmt.Errorf("failed to encrypt group data before storing in database: %v", err)
	}

	_, err = s.db.ExecContext(
		ctx,
		"INSERT INTO [group] VALUES (?, ?, ?, ?, ?, ?)",
		g.Id, ghash, psalt, phash, ekey, data,
	)
	if err != nil {
		var g model.Group
		return g, fmt.Errorf("failed to store group information in database: %v", err)
	}

	return g, nil
}

func (s GroupService) GetGroupName(
	ctx context.Context, p model.GetGroupNameParams,
) (string, error) {
	g, err := s.getDecryptedGroupData(ctx, p.PassKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt group data: %v", err)
	}
	return g.Name, nil
}

func (s GroupService) GetGroupDescription(
	ctx context.Context, p model.GetGroupDescriptionParams,
) (string, error) {
	g, err := s.getDecryptedGroupData(ctx, p.PassKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt group data: %v", err)
	}
	return g.Description, nil
}

func (s GroupService) GetGroupPassKey(
	ctx context.Context, params model.GetGroupPassKeyParams,
) (crypto.Key, error) {
	var salt crypto.Salt
	var hash crypto.PassHash
	err := s.db.QueryRowContext(ctx, "SELECT psalt, phash FROM [group]").Scan(&salt, &hash)
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
	err := s.db.QueryRowContext(ctx, "SELECT ekey FROM [group]").Scan(&ekey)
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
	ghash := crypto.HashData([]byte(params.GroupId))
	var phash crypto.PassHash
	err := s.db.QueryRowContext(
		ctx,
		"SELECT phash FROM [group] WHERE ghash = ?",
		ghash,
	).Scan(&phash)

	if err != nil {
		return fmt.Errorf("failed to get password hash from database: %v", err)
	}

	if !crypto.CheckPasswordHash(params.Password, phash) {
		return fmt.Errorf("password authentication failed for group: %v", err)
	}

	return nil
}

func (s GroupService) UpdateGroup(ctx context.Context, params model.UpdateGroupParams) error {
	g, err := s.getDecryptedGroupData(ctx, params.PassKey)
	if err != nil {
		return fmt.Errorf("failed to decrypt group data: %v", err)
	}

	if params.GroupId != nil {
		g.GroupId = *params.GroupId
	}
	if params.Name != nil {
		g.Name = *params.Name
	}
	if params.Description != nil {
		g.Description = *params.Description
	}

	return s.setEncryptedGroupData(ctx, params.PassKey, g)
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
		"UPDATE [group] SET psalt = ?, phash = ?, ekey = ?",
		salt, hash, ekey,
	)
	if err != nil {
		return fmt.Errorf("failed to store new password information in database: %v", err)
	}

	return nil
}

func (s GroupService) getDecryptedGroupData(
	ctx context.Context, pkey crypto.Key,
) (model.Group, error) {
	dkey, err := s.GetGroupDataKey(ctx, model.GetGroupDataKeyParams{PassKey: pkey})
	if err != nil {
		var g model.Group
		return g, fmt.Errorf("failed to get group data key: %v", err)
	}

	var data []byte
	err = s.db.QueryRowContext(ctx, "SELECT data FROM [group]").Scan(&data)
	if err != nil {
		var g model.Group
		return g, fmt.Errorf("failed to get group data from database: %v", err)
	}

	a := crypto.NewAgent[model.Group](dkey)
	return a.Decrypt(data)
}

func (s GroupService) setEncryptedGroupData(
	ctx context.Context, pkey crypto.Key, g model.Group,
) error {
	dkey, err := s.GetGroupDataKey(ctx, model.GetGroupDataKeyParams{PassKey: pkey})
	if err != nil {
		return fmt.Errorf("failed to get group data key for encryption: %v", err)
	}

	a := crypto.NewAgent[model.Group](dkey)
	data, err := a.Encrypt(g)
	if err != nil {
		return fmt.Errorf("failed to encrypt updated group data: %v", err)
	}

	_, err = s.db.ExecContext(ctx, "UPDATE [group] SET data = ?", data)
	if err != nil {
		return fmt.Errorf("failed to store updated group data in database: %v", err)
	}

	return nil
}

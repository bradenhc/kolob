// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/bradenhc/kolob/internal/crypto"
	"github.com/bradenhc/kolob/internal/fail"
	"github.com/bradenhc/kolob/internal/model"
)

type GroupService struct {
	db *sql.DB
}

func NewGroupService(db *sql.DB) GroupService {
	return GroupService{db}
}

func (s *GroupService) InitGroup(
	ctx context.Context, params model.InitGroupParams,
) (model.Group, error) {
	// Make sure there isn't already a set of group information in the database
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM [group]").Scan(&count)
	if err != nil {
		return fail.Zero[model.Group]("failed to check for existing group information", err)
	}
	if count != 0 {
		var g model.Group
		return g, errors.New("group has already been initialized")
	}

	g, err := model.NewGroup(params.GroupId, params.Name, params.Description)
	if err != nil {
		var g model.Group
		return g, fmt.Errorf("failed to initialize group: %v", err)
	}

	// Generate a new random key that will encrypt all data for the group (data key)
	slog.Info("Generating data encryption key", "group", g.Id)
	dkey, err := crypto.NewRandomKey()
	if err != nil {
		var g model.Group
		return g, fmt.Errorf("failed to create new data encryption key for group: %v", err)
	}

	// Generate a second key using the user password that will encrypt the data key (pass key)
	psalt, err := crypto.NewSalt()
	if err != nil {
		var g model.Group
		return g, fmt.Errorf("failed to create salt for group: %v", err)
	}
	slog.Info("Deriving password key", "group", g.Id)
	pkey := crypto.NewDerivedKey(params.Password, psalt)

	// Encyrypt the data key using the pass key before we store it in the database
	slog.Info("Encrypting data key", "group", g.Id)
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
	slog.Info("Encrypting group information", "group", g.Id)
	a := crypto.NewAgent[model.Group](dkey)
	data, err := a.Encrypt(g)
	if err != nil {
		var g model.Group
		return g, fmt.Errorf("failed to encrypt group data before storing in database: %v", err)
	}

	// Format the datetimes as RFC3339-compliant strings for storage in the DB
	created := g.CreatedAt.Format(time.RFC3339)
	updated := g.UpdatedAt.Format(time.RFC3339)

	slog.Info("Adding group information to database", "group", g.Id)
	_, err = s.db.ExecContext(
		ctx,
		"INSERT INTO [group] VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		g.Id, ghash[:], psalt, phash, ekey, created, updated, data,
	)
	if err != nil {
		var g model.Group
		return g, fmt.Errorf("failed to store group information in database: %v", err)
	}

	return g, nil
}

func (s *GroupService) GetGroupInfo(
	ctx context.Context, p model.GetGroupInfoParams,
) (model.Group, error) {
	var g model.Group
	eda, err := NewEncryptedDataAccessor[model.Group](ctx, s.db, "group", p.PassKey)
	if err != nil {
		return g, fmt.Errorf("failed to create encrypted data accessor: %v", err)
	}

	g, err = eda.Get(ctx, s.db, p.Id)
	if err != nil {
		var g model.Group
		return g, fmt.Errorf("failed to decrypt group data: %v", err)
	}

	return g, nil
}

func (s *GroupService) GetGroupPassKey(
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

func (s *GroupService) GetGroupDataKey(
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

func (s *GroupService) AuthenticateGroup(
	ctx context.Context, params model.AuthenticateGroupParams,
) error {
	ghash := crypto.HashData([]byte(params.GroupId))
	var phash crypto.PassHash
	err := s.db.QueryRowContext(
		ctx,
		"SELECT phash FROM [group] WHERE idhash = ?",
		ghash[:],
	).Scan(&phash)

	if err != nil {
		return fmt.Errorf("failed to get password hash from database: %v", err)
	}

	if !crypto.CheckPasswordHash(params.Password, phash) {
		return fmt.Errorf("password authentication failed")
	}

	return nil
}

func (s *GroupService) UpdateGroup(ctx context.Context, p model.UpdateGroupParams) error {
	eda, err := NewEncryptedDataAccessor[model.Group](ctx, s.db, "group", p.PassKey)
	if err != nil {
		return fmt.Errorf("failed to create encrypted data accessor: %v", err)
	}

	g, err := eda.Get(ctx, s.db, p.Id)
	if err != nil {
		return fmt.Errorf("failed to decrypt group data: %v", err)
	}

	if p.GroupId != nil {
		g.GroupId = *p.GroupId
	}
	if p.Name != nil {
		g.Name = *p.Name
	}
	if p.Description != nil {
		g.Description = *p.Description
	}

	g.UpdatedAt = time.Now()

	return eda.SetWithIdHash(ctx, s.db, p.Id, crypto.HashData([]byte(g.GroupId)), g)
}

func (s *GroupService) ChangeGroupPassword(
	ctx context.Context, params model.ChangeGroupPasswordParams,
) error {
	passParams := model.GetGroupPassKeyParams{
		Password: params.OldPass,
	}
	pkey, err := s.GetGroupPassKey(ctx, passParams)
	if err != nil {
		return fmt.Errorf("failed to get group pass key: %v", err)
	}

	dataParams := model.GetGroupDataKeyParams{
		PassKey: pkey,
	}
	dkey, err := s.GetGroupDataKey(ctx, dataParams)
	if err != nil {
		return fmt.Errorf("failed to get group data key: %v", err)
	}

	psalt, err := crypto.NewSalt()
	if err != nil {
		return fmt.Errorf("failed to generate salt for new password: %v", err)
	}

	phash, err := crypto.HashPassword(params.NewPass)
	if err != nil {
		return fmt.Errorf("failed to generate hash for new password: %v", err)
	}

	pkey = crypto.NewDerivedKey(params.NewPass, psalt)

	ekey, err := crypto.Encrypt(pkey, dkey)
	if err != nil {
		return fmt.Errorf("failed to encrypt data key with new password key: %v", err)
	}

	updated := time.Now().Format(time.RFC3339)

	_, err = s.db.ExecContext(
		ctx,
		"UPDATE [group] SET psalt = ?, phash = ?, ekey = ?, updated = ? WHERE id = ?",
		psalt, phash, ekey, updated, params.Id,
	)
	if err != nil {
		return fmt.Errorf("failed to store new password information in database: %v", err)
	}

	return nil
}

// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bradenhc/kolob/internal/crypto"
	"github.com/bradenhc/kolob/internal/model"
)

type MemberService struct {
	db *sql.DB
}

func NewMemberService(db *sql.DB) MemberService {
	return MemberService{db}
}

func (s *MemberService) CreateMember(
	ctx context.Context, p model.CreateMemberParams,
) (model.Member, error) {
	// Create the new member
	m, err := model.NewMember(p.Username, p.Name)
	if err != nil {
		return m, fmt.Errorf("failed to create new member: %v", err)
	}

	// Hash the username so that we can store it in the DB without leaking information
	uhash := crypto.HashData([]byte(p.Username))

	// Hash the password so we can use it for authentication
	phash, err := crypto.HashPassword(p.Password)
	if err != nil {
		var m model.Member
		return m, fmt.Errorf("failed to hash member password: %v", err)
	}

	// Format the datetimes as RFC3339-compliant strings for storage in the DB
	created := m.CreatedAt.Format(time.RFC3339)
	updated := m.UpdatedAt.Format(time.RFC3339)

	// Get the group data key in preparation for encrypting the member data
	gs := NewGroupService(s.db)
	dkey, err := gs.GetGroupDataKey(ctx, model.GetGroupDataKeyParams{PassKey: p.PassKey})
	if err != nil {
		var m model.Member
		return m, fmt.Errorf("failed to get group data key: %v", err)
	}

	// Encrypt the member data prior to storing it in the DB
	a := crypto.NewAgent[model.Member](dkey)
	data, err := a.Encrypt(m)
	if err != nil {
		var m model.Member
		return m, fmt.Errorf("failed to encrypt member data: %v", err)
	}

	// Store the new member in the DB
	_, err = s.db.ExecContext(
		ctx,
		"INSERT INTO member VALUES (?, ?, ?, ?, ?, ?)",
		m.Id, created, updated, uhash, phash, data,
	)
	if err != nil {
		var m model.Member
		return m, fmt.Errorf("failed to store new member in database: %v", err)
	}

	return m, nil
}

func (s *MemberService) UpdateMember(ctx context.Context, p model.UpdateMemberParams) error {
	return nil
}

func (s *MemberService) RemoveMember(ctx context.Context, p model.RemoveMemberParams) error {
	return nil
}

func (s *MemberService) ListMembers(
	ctx context.Context, p model.ListMembersParams,
) ([]model.Member, error) {
	return nil, nil
}

func (s *MemberService) FindMemberByUsername(
	ctx context.Context, p model.FindMemberByUsernameParams,
) (model.Member, error) {
	var m model.Member
	return m, nil
}

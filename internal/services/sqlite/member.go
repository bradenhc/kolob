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
	m, err := model.NewMember(p.Username, p.Name)
	if err != nil {
		return m, fmt.Errorf("failed to create new member: %v", err)
	}

	hash, err := crypto.HashPassword(p.Password)
	if err != nil {
		var m model.Member
		return m, fmt.Errorf("failed to hash member password: %v", err)
	}

	a := crypto.NewAgent[model.Member](p.Key)
	data, err := a.Encrypt(m)
	if err != nil {
		var m model.Member
		return m, fmt.Errorf("failed to encrypt member data: %v", err)
	}

	_, err = s.db.ExecContext(
		ctx,
		"INSERT INTO member VALUES (?, ?, ?, ?, ?)",
		m.Id, m.CreatedAt.Format(time.RFC3339), m.UpdatedAt.Format(time.RFC3339), hash, data,
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

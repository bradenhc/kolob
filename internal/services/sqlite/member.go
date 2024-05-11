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

	// Hash the username so that we can store it in the DB without leaking information and use it
	// for fast lookup later
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

	// Encrypt the member data prior to storing it in the DB
	eda, err := NewEncryptedDataAccessor[model.Member](ctx, s.db, "member", p.PassKey)
	if err != nil {
		var m model.Member
		return m, fmt.Errorf("failed to create encrypted data accessor: %v", err)
	}

	data, err := eda.Encrypt(m)
	if err != nil {
		var m model.Member
		return m, fmt.Errorf("failed to encrypt member data: %v", err)
	}

	// Store the new member in the DB
	_, err = s.db.ExecContext(
		ctx,
		"INSERT INTO member VALUES (?, ?, ?, ?, ?, ?)",
		m.Id, created, updated, uhash[:], phash, data,
	)
	if err != nil {
		var m model.Member
		return m, fmt.Errorf("failed to store new member in database: %v", err)
	}

	return m, nil
}

func (s *MemberService) AuthenticateMember(
	ctx context.Context, p model.AuthenticateMemberParams,
) error {
	idhash := crypto.HashData([]byte(p.Username))
	var phash crypto.PassHash
	err := s.db.QueryRowContext(
		ctx,
		"SELECT phash FROM member WHERE idhash = ?",
		idhash[:],
	).Scan(&phash)
	if err != nil {
		return fmt.Errorf("failed to get member auth from database: %v", err)
	}

	if !crypto.CheckPasswordHash(p.Password, phash) {
		return fmt.Errorf("password authentication failed")
	}

	return nil
}

func (s *MemberService) UpdateMember(ctx context.Context, p model.UpdateMemberParams) error {
	eda, err := NewEncryptedDataAccessor[model.Member](ctx, s.db, "member", p.PassKey)
	if err != nil {
		return fmt.Errorf("failed to create encrypted data accessor: %v", err)
	}

	m, err := eda.Get(ctx, s.db, p.Id)
	if err != nil {
		return fmt.Errorf("failed to get encrypted member data: %v", m)
	}

	if p.Username != nil {
		m.Username = *p.Username
	}
	if p.Name != nil {
		m.Name = *p.Name
	}

	m.UpdatedAt = time.Now()

	return eda.SetWithIdHash(ctx, s.db, m.Id, crypto.HashData([]byte(m.Username)), m)
}

func (s *MemberService) RemoveMember(ctx context.Context, p model.RemoveMemberParams) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM member WHERE id = ?", p.Id)
	if err != nil {
		return fmt.Errorf("failed to remove member from database: %v", err)
	}

	return nil
}

func (s *MemberService) ListMembers(
	ctx context.Context, p model.ListMembersParams,
) ([]model.Member, error) {
	eda, err := NewEncryptedDataAccessor[model.Member](ctx, s.db, "member", p.PassKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create encrypted data accessor: %v", err)
	}

	return eda.GetList(ctx, s.db)
}

func (s *MemberService) FindMemberByUsername(
	ctx context.Context, p model.FindMemberByUsernameParams,
) (model.Member, error) {
	eda, err := NewEncryptedDataAccessor[model.Member](ctx, s.db, "member", p.PassKey)
	if err != nil {
		var m model.Member
		return m, fmt.Errorf("failed to create encrypted data accessor: %v", err)
	}

	return eda.GetByIdHash(ctx, s.db, crypto.HashData([]byte(p.Username)))
}

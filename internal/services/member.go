// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package services

import (
	"context"
	"fmt"

	"github.com/bradenhc/kolob/internal/crypto"
	"github.com/bradenhc/kolob/internal/model"
	"github.com/bradenhc/kolob/internal/store"
)

type MemberService struct {
	store store.MemberStore
}

func NewMemberService(store store.MemberStore) MemberService {
	return MemberService{store}
}

func (s *MemberService) Create(
	ctx context.Context, req *MemberCreateRequest, key crypto.Key,
) (*model.Member, error) {
	// Create the new member
	m, err := model.NewMember(string(req.Username()), string(req.Name()))
	if err != nil {
		return m, fmt.Errorf("failed to create new member: %v", err)
	}

	var meta store.MemberMetadata
	meta.Id = model.Uuid(m.Id())

	// Hash the username so that we can store it in the DB without leaking information and use it
	// for fast lookup later
	meta.UsernameHash = crypto.HashData(req.Username())

	// Hash the password so we can use it for authentication
	meta.PassHash, err = crypto.HashPassword(crypto.Password(req.Password()))
	if err != nil {
		return nil, fmt.Errorf("failed to hash member password: %v", err)
	}

	// Format the datetimes as RFC3339-compliant strings for storage in the DB
	meta.CreatedAt = m.Created()
	meta.UpdatedAt = m.Updated()

	// Encrypt the member edata prior to storing it in the DB
	edata, err := crypto.Encrypt(key, m.Table().Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt member data before storing: %v", err)
	}

	// Store the new member in the DB
	err = s.store.AddMemberData(ctx, meta, edata)
	if err != nil {
		return nil, fmt.Errorf("failed to store member: %v", err)
	}

	return m, nil
}

func (s *MemberService) Authenticate(
	ctx context.Context, req *MemberAuthenticateRequest, key crypto.Key,
) (*model.Member, error) {
	uhash := crypto.HashData(req.Username())
	meta, edata, err := s.store.GetMemberDataByUname(ctx, uhash)
	if err != nil {
		return nil, fmt.Errorf("failed to get member: %v", err)
	}

	if !crypto.CheckPasswordHash(crypto.Password(req.Password()), meta.PassHash) {
		return nil, fmt.Errorf("password authentication failed")
	}

	data, err := crypto.Decrypt(key, edata)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt member data: %v", err)
	}

	return model.GetRootAsMember(data, 0), nil
}

func (s *MemberService) ChangePassword(
	ctx context.Context, req *MemberChangePasswordRequest,
) error {
	meta, _, err := s.store.GetMemberData(ctx, model.Uuid(req.Id()))
	if err != nil {
		return fmt.Errorf("failed to get member data: %v", err)
	}

	if !crypto.CheckPasswordHash(crypto.Password(req.OldPassword()), meta.PassHash) {
		return fmt.Errorf("invalid credentials")
	}

	npass, err := crypto.NewPassword(string(req.NewPassword()))
	if err != nil {
		return fmt.Errorf("invalid password: %v", err)
	}

	meta.PassHash, err = crypto.HashPassword(npass)
	if err != nil {
		return fmt.Errorf("failed to hash password for storage: %v", err)
	}

	return nil
}

func (s *MemberService) UpdateMember(
	ctx context.Context, req *MemberUpdateRequest, key crypto.Key,
) (*model.Member, error) {
	meta, edata, err := s.store.GetMemberData(ctx, model.Uuid(req.Id()))
	if err != nil {
		return nil, fmt.Errorf("failed to get member data: %v", err)
	}

	data, err := crypto.Decrypt(key, edata)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt member data: %v", err)
	}

	prev := model.GetRootAsMember(data, 0)
	m := model.CloneMemberWithUpdates(prev, req.Username(), req.Name())

	meta.UpdatedAt = m.Updated()
	meta.UsernameHash = crypto.HashData(m.Uname())

	edata, err = crypto.Encrypt(key, m.Table().Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt member data: %v", err)
	}

	err = s.store.UpdateMemberData(ctx, meta, edata)
	if err != nil {
		return nil, fmt.Errorf("failed to store updated member data: %v", err)
	}

	return m, nil
}

func (s *MemberService) RemoveMember(ctx context.Context, req *MemberRemoveRequest) error {
	err := s.store.RemoveMemberData(ctx, model.Uuid(req.Id()))
	if err != nil {
		return fmt.Errorf("failed to remove member data: %v", err)
	}
	return nil
}

func (s *MemberService) ListMembers(ctx context.Context, key crypto.Key) ([]*model.Member, error) {
	edatas, err := s.store.ListMemberData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get member list: %v", err)
	}

	ms := make([]*model.Member, 0)
	for i := range edatas {
		data, err := crypto.Decrypt(key, edatas[i])
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt member in list: %v", err)
		}
		ms = append(ms, model.GetRootAsMember(data, 0))
	}

	return ms, nil
}

func (s *MemberService) FindMemberByUsername(
	ctx context.Context, req *MemberFindByUsernameRequest, key crypto.Key,
) (*model.Member, error) {
	uhash := crypto.DataHash(req.Username())
	_, edata, err := s.store.GetMemberDataByUname(ctx, uhash)
	if err != nil {
		return nil, fmt.Errorf("failed to get member by username: %v", err)
	}

	data, err := crypto.Decrypt(key, edata)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt member data: %v", err)
	}

	return model.GetRootAsMember(data, 0), nil
}

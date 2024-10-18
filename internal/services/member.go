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
	// Verify the password
	upass, err := crypto.NewPassword(string(req.Password()))
	if err != nil {
		return nil, fmt.Errorf("password validation failed: %v", err)
	}

	// Create the new member
	m, err := model.NewMember(string(req.Username()), string(req.Name()))
	if err != nil {
		return nil, fmt.Errorf("failed to create new member: %v", err)
	}

	// Create the entity we will store in the database
	entity, err := store.NewMemberEntity(m, upass, key)
	if err != nil {
		return nil, fmt.Errorf("failed to create new member entity: %v", err)
	}

	// Store the new member in the DB
	err = s.store.AddMemberEntity(ctx, entity)
	if err != nil {
		return nil, fmt.Errorf("failed to store member: %v", err)
	}

	return m, nil
}

func (s *MemberService) Authenticate(
	ctx context.Context, req *MemberAuthenticateRequest, key crypto.Key,
) (*model.Member, error) {
	uhash := crypto.HashData(req.Username())
	entity, err := s.store.GetMemberEntityByUname(ctx, uhash)
	if err != nil {
		return nil, fmt.Errorf("failed to get member: %v", err)
	}

	if !crypto.CheckPasswordHash(crypto.Password(req.Password()), entity.PassHash) {
		return nil, fmt.Errorf("password authentication failed")
	}

	m, err := entity.Decrypt(key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt member: %v", err)
	}

	return m, nil
}

func (s *MemberService) ChangePassword(
	ctx context.Context, req *MemberChangePasswordRequest,
) error {
	entity, err := s.store.GetMemberEntity(ctx, model.Uuid(req.Id()))
	if err != nil {
		return fmt.Errorf("failed to get member data: %v", err)
	}

	if !crypto.CheckPasswordHash(crypto.Password(req.OldPassword()), entity.PassHash) {
		return fmt.Errorf("invalid credentials")
	}

	npass, err := crypto.NewPassword(string(req.NewPassword()))
	if err != nil {
		return fmt.Errorf("invalid password: %v", err)
	}

	entity.PassHash, err = crypto.HashPassword(npass)
	if err != nil {
		return fmt.Errorf("failed to hash password for storage: %v", err)
	}

	err = s.store.UpdateMemberEntity(ctx, entity)
	if err != nil {
		return fmt.Errorf("failed to store updated member: %v", err)
	}

	return nil
}

func (s *MemberService) UpdateMember(
	ctx context.Context, req *MemberUpdateRequest, key crypto.Key,
) (*model.Member, error) {
	entity, err := s.store.GetMemberEntity(ctx, model.Uuid(req.Id()))
	if err != nil {
		return nil, fmt.Errorf("failed to get member data: %v", err)
	}

	m, err := entity.Update(key, req.Username(), req.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to update member entity: %v", err)
	}

	err = s.store.UpdateMemberEntity(ctx, entity)
	if err != nil {
		return nil, fmt.Errorf("failed to store updated member data: %v", err)
	}

	return m, nil
}

func (s *MemberService) RemoveMember(ctx context.Context, req *MemberRemoveRequest) error {
	err := s.store.RemoveMemberEntity(ctx, model.Uuid(req.Id()))
	if err != nil {
		return fmt.Errorf("failed to remove member data: %v", err)
	}
	return nil
}

func (s *MemberService) ListMembers(ctx context.Context, key crypto.Key) ([]*model.Member, error) {
	entities, err := s.store.ListMemberEntities(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get member list: %v", err)
	}

	ms := make([]*model.Member, 0, len(entities))
	for _, e := range entities {
		m, err := e.Decrypt(key)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt member in list: %v", err)
		}
		ms = append(ms, m)
	}

	return ms, nil
}

func (s *MemberService) FindMemberByUsername(
	ctx context.Context, req *MemberFindByUsernameRequest, key crypto.Key,
) (*model.Member, error) {
	uhash := crypto.HashData(req.Username())
	entity, err := s.store.GetMemberEntityByUname(ctx, uhash)
	if err != nil {
		return nil, fmt.Errorf("failed to get member by username: %v", err)
	}

	m, err := entity.Decrypt(key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt member data: %v", err)
	}

	return m, nil
}

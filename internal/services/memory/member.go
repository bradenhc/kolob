// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package memory

import (
	"context"
	"fmt"
	"regexp"
	"slices"

	"github.com/bradenhc/kolob/internal/model"
)

type MemberService struct {
	store *SynchronizedStore
}

func (s *MemberService) CreateMember(
	ctx context.Context, p model.CreateMemberParams,
) (model.Member, error) {
	s.store.membersLock.Lock()
	defer s.store.membersLock.Unlock()

	for _, m := range s.store.members {
		if m.Username == p.Username {
			var m model.Member
			return m, fmt.Errorf("member with username %s already exists", p.Username)
		}
	}

	m, err := model.NewMember(p.Username, p.Name)
	if err != nil {
		return m, err
	}

	s.store.members = append(s.store.members, &m)
	return m, nil
}

func (s *MemberService) UpdateMember(ctx context.Context, p model.UpdateMemberParams) error {
	s.store.membersLock.Lock()
	defer s.store.membersLock.Unlock()

	for _, m := range s.store.members {
		if m.Id == p.Id {
			if p.Username != nil {
				m.Username = *p.Username
			}
			if p.Name != nil {
				m.Name = *p.Name
			}
			return nil
		}
	}

	return fmt.Errorf("member with ID %s does not exist", p.Id)
}

func (s *MemberService) RemoveMember(ctx context.Context, p model.RemoveMemberParams) error {
	s.store.membersLock.Lock()
	defer s.store.membersLock.Unlock()

	for i, m := range s.store.members {
		if m.Id == p.Id {
			s.store.members = slices.Delete(s.store.members, i, i+1)
			break
		}
	}

	return nil
}

func (s *MemberService) ListMembers(
	ctx context.Context, p model.ListMembersParams,
) ([]model.Member, error) {
	s.store.membersLock.RLock()
	defer s.store.membersLock.RUnlock()

	ret := make([]model.Member, 0)
	for _, m := range s.store.members {
		if p.NamePattern != nil {
			match, _ := regexp.MatchString(*p.NamePattern, m.Name)
			if !match {
				continue
			}
		}
		ret = append(ret, *m)
	}

	return ret, nil
}

func (s *MemberService) FindMemberByUsername(
	ctx context.Context, p model.FindMemberByUsernameParams,
) (model.Member, error) {
	s.store.membersLock.RLock()
	defer s.store.membersLock.RUnlock()

	for _, m := range s.store.members {
		if m.Username == p.Username {
			return *m, nil
		}
	}

	var m model.Member
	return m, fmt.Errorf("member with username %s does not exist", p.Username)
}

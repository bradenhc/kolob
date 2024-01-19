package memory

import (
	"context"
	"fmt"
	"regexp"

	"github.com/bradenhc/kolob/internal/model"
)

type MemberService struct {
	store *SynchronizedStore
}

func (s *MemberService) members(gid model.Uuid) ([]*model.Member, error) {
	ms := s.store.members[gid]
	if ms == nil {
		return nil, fmt.Errorf("group with ID %s does not exist", gid)
	}
	return ms, nil
}

func (s *MemberService) CreateMember(ctx context.Context, p model.CreateMemberParams) (model.Member, error) {
	s.store.mutex.Lock()
	defer s.store.mutex.Unlock()

	ms, err := s.members(p.GroupId)
	if err != nil {
		var m model.Member
		return m, err
	}

	for _, m := range ms {
		if m.Username == p.Username {
			var m model.Member
			return m, fmt.Errorf("member with username %s already exists", p.Username)
		}
	}

	m, err := model.NewMember(p.Username, p.Name)
	if err != nil {
		return m, err
	}

	s.store.members[p.GroupId] = append(s.store.members[p.GroupId], &m)
	return m, nil
}

func (s *MemberService) UpdateMember(ctx context.Context, p model.UpdateMemberParams) error {
	s.store.mutex.Lock()
	defer s.store.mutex.Unlock()

	ms, err := s.members(p.GroupId)
	if err != nil {
		return err
	}

	for _, m := range ms {
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
	s.store.mutex.Lock()
	defer s.store.mutex.Unlock()

	ms, err := s.members(p.GroupId)
	if err != nil {
		return err
	}

	for i, m := range ms {
		if m.Id == p.Id {
			nms := make([]*model.Member, len(ms)-1)
			nms = append(nms, ms[:i]...)
			if i < len(ms)-1 {
				nms = append(nms, ms[i+1:]...)
			}
			s.store.members[p.GroupId] = nms
			break
		}
	}

	return nil
}

func (s *MemberService) ListMembers(ctx context.Context, p model.ListMembersParams) ([]model.Member, error) {
	s.store.mutex.RLock()
	defer s.store.mutex.RUnlock()

	ms, err := s.members(p.GroupId)
	if err != nil {
		return nil, err
	}

	ret := make([]model.Member, 0)
	for _, m := range ms {
		match, _ := regexp.MatchString(*p.NamePattern, m.Name)
		if match {
			ret = append(ret, *m)
		}
	}

	return ret, nil
}

func (s *MemberService) FindMemberByUsername(ctx context.Context, p model.FindMemberByUsernameParams) (model.Member, error) {
	s.store.mutex.RLock()
	defer s.store.mutex.RUnlock()

	ms, err := s.members(p.GroupId)
	if err != nil {
		var m model.Member
		return m, err
	}

	for _, m := range ms {
		if m.Username == p.Username {
			return *m, nil
		}
	}

	var m model.Member
	return m, fmt.Errorf("member with username %s does not exist", p.Username)
}

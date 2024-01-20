package memory

import (
	"context"
	"fmt"
	"regexp"

	"github.com/bradenhc/kolob/internal/model"
)

type GroupService struct {
	store *SynchronizedStore
}

func (s *GroupService) CreateGroup(ctx context.Context, p model.CreateGroupParams) (model.Group, error) {
	s.store.lock.Lock()
	defer s.store.lock.Unlock()

	var g model.Group
	for _, v := range s.store.groups {
		if p.Name == v.Name {
			return g, fmt.Errorf("group with the name '%s' already exists", p.Name)
		}
	}

	g, err := model.NewGroup(p.Name)
	if err != nil {
		return g, err
	}

	m, err := model.NewMember(p.OwnerUsername, p.OwnerName)
	if err != nil {
		return g, err
	}
	m.IsOwner = true

	s.store.groups[g.Id] = &g
	s.store.members[g.Id] = append(s.store.members[g.Id], &m)

	return g, nil
}

func (s *GroupService) UpdateGroup(ctx context.Context, p model.UpdateGroupParams) error {
	s.store.lock.Lock()
	defer s.store.lock.Unlock()

	g := s.store.groups[p.Id]
	if g == nil {
		return fmt.Errorf("group with ID %s does not exist", p.Id)
	}

	if p.Name != nil {
		g.Name = *p.Name
	}

	return nil
}

func (s *GroupService) RemoveGroup(ctx context.Context, p model.RemoveGroupParams) error {
	s.store.lock.Lock()
	defer s.store.lock.Unlock()

	// This has to be a cascading delete from all the internal maps. Start with messages,
	// conversations, then members, then finally remove the group
	cs := s.store.conversations[p.Id]
	if cs != nil {
		for _, c := range cs {
			delete(s.store.messages, c.Id)
		}
		delete(s.store.conversations, p.Id)
	}

	delete(s.store.members, p.Id)
	delete(s.store.groups, p.Id)

	return nil
}

func (s *GroupService) ListGroups(ctx context.Context, p model.ListGroupsParams) ([]model.Group, error) {
	s.store.lock.RLock()
	defer s.store.lock.RUnlock()

	gs := make([]model.Group, len(s.store.groups))
	for _, g := range s.store.groups {
		if p.Pattern != nil {
			match, _ := regexp.MatchString(*p.Pattern, g.Name)
			if !match {
				continue
			}
		}
		gs = append(gs, *g)
	}

	return gs, nil
}

func (s *GroupService) FindGroupById(ctx context.Context, id model.Uuid) (model.Group, error) {
	s.store.lock.RLock()
	defer s.store.lock.RUnlock()

	g, ok := s.store.groups[id]
	if !ok {
		var g model.Group
		return g, fmt.Errorf("group with ID %s does not exist", id)
	}
	return *g, nil
}

func (s *GroupService) FindGroupByName(ctx context.Context, name string) (model.Group, error) {
	s.store.lock.RLock()
	defer s.store.lock.RUnlock()

	for _, g := range s.store.groups {
		if g.Name == name {
			return *g, nil
		}
	}

	var g model.Group
	return g, fmt.Errorf("group with name '%s' does not exist", name)
}

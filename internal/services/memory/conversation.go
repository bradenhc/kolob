package memory

import (
	"context"
	"fmt"
	"regexp"

	"github.com/bradenhc/kolob/internal/model"
)

type ConversationService struct {
	store *SynchronizedStore
}

func (s *ConversationService) conversations(gid model.Uuid) ([]*model.Conversation, error) {
	cs, ok := s.store.conversations[gid]
	if !ok {
		return nil, fmt.Errorf("group with ID %s does not exist", gid)
	}
	return cs, nil
}

func (s *ConversationService) CreateConversation(ctx context.Context, p model.CreateConversationParams) (model.Conversation, error) {
	s.store.lock.Lock()
	defer s.store.lock.Unlock()

	cs, err := s.conversations(p.GroupId)
	if err != nil {
		var c model.Conversation
		return c, err
	}

	for _, c := range cs {
		if c.Name == p.Name {
			var c model.Conversation
			return c, fmt.Errorf("conversation with name '%s' already exists", p.Name)
		}
	}

	c, err := model.NewConversation(p.Name, p.Moderator)
	if err != nil {
		var c model.Conversation
		return c, err
	}

	s.store.conversations[p.GroupId] = append(s.store.conversations[p.GroupId], &c)
	return c, nil
}

func (s *ConversationService) UpdateConversation(ctx context.Context, p model.UpdateConversationParams) error {
	s.store.lock.Lock()
	defer s.store.lock.Unlock()

	cs, err := s.conversations(p.GroupId)
	if err != nil {
		return err
	}

	for _, c := range cs {
		if c.Id == p.Id {
			if p.Name != nil {
				c.Name = *p.Name
			}
			if p.Moderator != nil {
				c.Moderator = *p.Moderator
			}
			return nil
		}
	}

	return fmt.Errorf("conversation with ID %s does not exist", p.Id)
}

func (s *ConversationService) RemoveConversation(ctx context.Context, p model.RemoveConversationParams) error {
	s.store.lock.Lock()
	defer s.store.lock.Unlock()

	cs, err := s.conversations(p.GroupId)
	if err != nil {
		return err
	}

	for i, c := range cs {
		if c.Id == p.Id {
			ncs := make([]*model.Conversation, len(cs)-1)
			ncs = append(ncs, cs[:i]...)
			if i < len(cs)-1 {
				ncs = append(ncs, cs[i+1:]...)
			}
			s.store.conversations[p.GroupId] = ncs
			break
		}
	}

	return nil

}

func (s *ConversationService) ListConversations(ctx context.Context, p model.ListConversationsParams) ([]model.Conversation, error) {
	s.store.lock.RLock()
	defer s.store.lock.RUnlock()

	cs, err := s.conversations(p.GroupId)
	if err != nil {
		return nil, err
	}

	ret := make([]model.Conversation, 0)
	for _, c := range cs {
		if p.Pattern != nil {
			match, _ := regexp.MatchString(*p.Pattern, c.Name)
			if !match {
				continue
			}
		}
		ret = append(ret, *c)
	}

	return ret, nil
}

func (s *ConversationService) FindConversationById(ctx context.Context, p model.FindConversationByIdParams) (model.Conversation, error) {
	s.store.lock.RLock()
	defer s.store.lock.RUnlock()

	cs, err := s.conversations(p.GroupId)
	if err != nil {
		var c model.Conversation
		return c, err
	}

	for _, c := range cs {
		if c.Id == p.Id {
			return *c, nil
		}
	}

	var c model.Conversation
	return c, fmt.Errorf("conversation with ID %s does not exist", p.Id)
}

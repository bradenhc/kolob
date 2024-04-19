// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package memory

import (
	"cmp"
	"context"
	"fmt"
	"regexp"
	"slices"

	"github.com/bradenhc/kolob/internal/model"
)

type ConversationService struct {
	store *SynchronizedStore
}

func (s *ConversationService) CreateConversation(ctx context.Context, p model.CreateConversationParams) (model.Conversation, error) {
	s.store.conversationsLock.Lock()
	defer s.store.conversationsLock.Unlock()

	for _, c := range s.store.conversations {
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

	s.store.conversations[c.Id] = &c
	return c, nil
}

func (s *ConversationService) UpdateConversation(ctx context.Context, p model.UpdateConversationParams) error {
	s.store.conversationsLock.Lock()
	defer s.store.conversationsLock.Unlock()

	c, ok := s.store.conversations[p.Id]
	if !ok {
		return fmt.Errorf("conversation with ID %s does not exist", p.Id)
	}

	if p.Name != nil {
		c.Name = *p.Name
	}
	if p.Moderator != nil {
		c.Moderator = *p.Moderator
	}

	return nil
}

func (s *ConversationService) RemoveConversation(ctx context.Context, p model.RemoveConversationParams) error {
	s.store.conversationsLock.Lock()
	defer s.store.conversationsLock.Unlock()

	delete(s.store.messages, p.Id)
	delete(s.store.conversations, p.Id)

	return nil
}

func (s *ConversationService) ListConversations(ctx context.Context, p model.ListConversationsParams) ([]model.Conversation, error) {
	s.store.conversationsLock.RLock()
	defer s.store.conversationsLock.RUnlock()

	ret := make([]model.Conversation, 0)
	for _, c := range s.store.conversations {
		if p.Pattern != nil {
			match, _ := regexp.MatchString(*p.Pattern, c.Name)
			if !match {
				continue
			}
		}
		ret = append(ret, *c)
	}

	slices.SortFunc(ret, func(a, b model.Conversation) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return ret, nil
}

func (s *ConversationService) FindConversationById(ctx context.Context, p model.FindConversationByIdParams) (model.Conversation, error) {
	s.store.conversationsLock.RLock()
	defer s.store.conversationsLock.RUnlock()

	c, ok := s.store.conversations[p.Id]
	if !ok {
		var c model.Conversation
		return c, fmt.Errorf("conversation with ID %s does not exist", p.Id)
	}

	return *c, nil
}

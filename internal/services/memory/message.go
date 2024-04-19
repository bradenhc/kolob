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

type MessageService struct {
	store *SynchronizedStore
}

func (s *MessageService) CreateMessage(ctx context.Context, p model.CreateMessageParams) (model.Message, error) {
	s.store.messagesLock.Lock()
	defer s.store.messagesLock.Unlock()

	if _, ok := s.store.messages[p.ConversationId]; !ok {
		var m model.Message
		return m, fmt.Errorf("conversation with ID %s does not exist", p.ConversationId)
	}

	m, err := model.NewMessage(p.Author, p.Content)
	if err != nil {
		return m, err
	}

	s.store.messages[p.ConversationId] = append(s.store.messages[p.ConversationId], &m)
	return m, nil
}

func (s *MessageService) UpdateMessage(ctx context.Context, p model.UpdateMessageParams) error {
	s.store.messagesLock.Lock()
	defer s.store.messagesLock.Unlock()

	ms, ok := s.store.messages[p.ConversationId]
	if !ok {
		return fmt.Errorf("conversation with ID %s does not exist", p.ConversationId)
	}

	for _, m := range ms {
		if m.Id == p.Id {
			if p.Content != nil {
				m.Content = *p.Content
			}
			return nil
		}
	}

	return fmt.Errorf("message with ID %s does not exist", p.Id)
}

func (s *MessageService) RemoveMessage(ctx context.Context, p model.RemoveMessageParams) error {
	s.store.messagesLock.Lock()
	defer s.store.messagesLock.Unlock()

	ms, ok := s.store.messages[p.ConversationId]
	if !ok {
		return fmt.Errorf("conversation with ID %s does not exist", p.ConversationId)
	}

	for i, m := range ms {
		if m.Id == p.Id {
			s.store.messages[p.ConversationId] = slices.Delete(ms, i, i+1)
			break
		}
	}

	return nil
}

func (s *MessageService) ListMessages(ctx context.Context, p model.ListMessagesParams) ([]model.Message, error) {
	s.store.messagesLock.RLock()
	defer s.store.messagesLock.RUnlock()

	ms, ok := s.store.messages[p.ConversationId]
	if !ok {
		return nil, fmt.Errorf("conversation with ID %s does not exist", p.ConversationId)
	}

	ret := make([]model.Message, 0)
	for _, m := range ms {
		if p.Pattern != nil {
			match, _ := regexp.MatchString(*p.Pattern, m.Content)
			if !match {
				continue
			}
		}
		ret = append(ret, *m)
	}

	return ret, nil
}

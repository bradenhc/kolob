package memory

import (
	"context"
	"fmt"
	"regexp"

	"github.com/bradenhc/kolob/internal/model"
)

type MessageService struct {
	store *SynchronizedStore
}

func (s *MessageService) messages(gid, cid model.Uuid) ([]*model.Message, error) {
	_, ok := s.store.groups[gid]
	if !ok {
		return nil, fmt.Errorf("group with ID %s does not exist", gid)
	}
	cs, ok := s.store.messages[cid]
	if !ok {
		return nil, fmt.Errorf("conversation with ID %s does not exist", cid)
	}
	return cs, nil
}

func (s *MessageService) CreateMessage(ctx context.Context, p model.CreateMessageParams) (model.Message, error) {
	s.store.lock.Lock()
	defer s.store.lock.Unlock()

	_, err := s.messages(p.GroupId, p.ConversationId)
	if err != nil {
		var m model.Message
		return m, err
	}

	m, err := model.NewMessage(p.Author, p.Content)
	if err != nil {
		return m, err
	}

	s.store.messages[p.ConversationId] = append(s.store.messages[p.ConversationId], &m)
	return m, nil
}

func (s *MessageService) UpdateMessage(ctx context.Context, p model.UpdateMessageParams) error {
	s.store.lock.Lock()
	defer s.store.lock.Unlock()

	ms, err := s.messages(p.GroupId, p.ConversationId)
	if err != nil {
		return err
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
	s.store.lock.Lock()
	defer s.store.lock.Unlock()

	ms, err := s.messages(p.GroupId, p.ConversationId)
	if err != nil {
		return err
	}

	for i, m := range ms {
		if m.Id == p.Id {
			nms := make([]*model.Message, len(ms)-1)
			nms = append(nms, ms[:i]...)
			if i < len(ms)-1 {
				nms = append(nms, ms[i+1:]...)
			}
			s.store.messages[p.ConversationId] = nms
			break

		}
	}

	return nil
}

func (s *MessageService) ListMessages(ctx context.Context, p model.ListMessagesParams) ([]model.Message, error) {
	s.store.lock.Lock()
	defer s.store.lock.Unlock()

	ms, err := s.messages(p.GroupId, p.ConversationId)
	if err != nil {
		return nil, err
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

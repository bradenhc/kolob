// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package services

import (
	"context"
	"fmt"
	"regexp"

	"github.com/bradenhc/kolob/internal/crypto"
	"github.com/bradenhc/kolob/internal/model"
	"github.com/bradenhc/kolob/internal/store"
)

type MessageService struct {
	store store.MessageStore
}

func NewMessageService(store store.MessageStore) MessageService {
	return MessageService{store}
}

func (s *MessageService) Add(
	ctx context.Context, req *MessageAddRequest, key crypto.Key,
) (*model.Message, error) {
	// Create the new message object
	m, err := model.NewMessage(
		model.Uuid(req.Author()), model.Uuid(req.Conversation()), string(req.Content()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create message object: %v", err)
	}

	entity, err := store.NewMessageEntity(m, key)
	if err != nil {
		return nil, fmt.Errorf("failed to create message entity: %v", err)
	}

	err = s.store.AddMessageEntity(ctx, entity)
	if err != nil {
		return nil, fmt.Errorf("failed to store message entity: %v", err)
	}

	return m, nil
}

func (s *MessageService) Get(
	ctx context.Context, req *MessageGetRequest, key crypto.Key,
) (*model.Message, error) {
	entity, err := s.store.GetMessageEntity(ctx, model.Uuid(req.Id()))
	if err != nil {
		return nil, fmt.Errorf("failed to get message entity from store: %v", err)
	}

	m, err := entity.Decrypt(key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt message entity: %v", err)
	}

	return m, nil
}

func (s *MessageService) Update(
	ctx context.Context, req *MessageUpdateRequest, key crypto.Key,
) (*model.Message, error) {
	entity, err := s.store.GetMessageEntity(ctx, model.Uuid(req.Id()))
	if err != nil {
		return nil, fmt.Errorf("failed to get message from store: %v", err)
	}

	next, err := entity.Update(key, req.Content())
	if err != nil {
		return nil, fmt.Errorf("failed to update message entity: %v", err)
	}

	err = s.store.UpdateMessageEntity(ctx, entity)
	if err != nil {
		return nil, fmt.Errorf("failed to store updated message: %v", err)
	}

	return next, nil
}

func (s *MessageService) Remove(ctx context.Context, req *MessageRemoveRequest) error {
	err := s.store.RemoveMessageEntity(ctx, model.Uuid(req.Id()))
	if err != nil {
		return fmt.Errorf("failed to remove message from store: %v", err)
	}
	return nil
}

func (s *MessageService) List(
	ctx context.Context, req *MessageListRequest, key crypto.Key,
) ([]*model.Message, error) {
	var query store.ListMessageDataQuery
	if req.Author() != nil {
		query.Author = new(model.Uuid)
		*query.Author = model.Uuid(req.Author())
	}
	if req.CreatedAfter() != 0 {
		query.CreatedAfter = new(int64)
		*query.CreatedAfter = req.CreatedAfter()
	}
	if req.CreatedBefore() != 0 {
		query.CreatedBefore = new(int64)
		*query.CreatedBefore = req.CreatedBefore()
	}

	entities, err := s.store.ListMessageEntities(ctx, model.Uuid(req.Conversation()), query)
	if err != nil {
		return nil, fmt.Errorf("failed to get message list from store: %v", err)
	}

	mlist := make([]*model.Message, 0, len(entities))
	for _, e := range entities {
		m, err := e.Decrypt(key)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt message in list: %v", err)
		}

		mlist = append(mlist, m)
	}

	if req.Pattern() == nil {
		return mlist, nil
	}

	// We have a content pattern, so we need to filter the list further
	r, err := regexp.Compile(string(req.Pattern()))
	if err != nil {
		return nil, fmt.Errorf("invalid content pattern: %v", err)
	}

	flist := make([]*model.Message, 0, len(mlist))
	for _, v := range mlist {
		if r.Match(v.Content()) {
			flist = append(flist, v)
		}
	}

	return flist, nil
}

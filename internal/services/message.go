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

	edata, err := crypto.Encrypt(key, m.Table().Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt message info: %v", err)
	}

	// Store the message in the database
	meta := store.MessageMetadata{
		Id:           model.Uuid(m.Id()),
		Conversation: model.Uuid(m.Conversation()),
		Author:       model.Uuid(m.Author()),
		CreatedAt:    m.Created(),
		UpdatedAt:    m.Updated(),
	}
	err = s.store.AddMessageData(ctx, meta, edata)
	if err != nil {
		return nil, fmt.Errorf("failed to store message: %v", err)
	}

	return m, nil
}

func (s *MessageService) Get(
	ctx context.Context, req *MessageGetRequest, key crypto.Key,
) (*model.Message, error) {
	_, edata, err := s.store.GetMessageData(ctx, model.Uuid(req.Id()))
	if err != nil {
		return nil, fmt.Errorf("failed to get message from store: %v", err)
	}

	data, err := crypto.Decrypt(key, edata)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt member data: %v", err)
	}

	return model.GetRootAsMessage(data, 0), nil
}

func (s *MessageService) Update(
	ctx context.Context, req *MessageUpdateRequest, key crypto.Key,
) error {
	meta, edata, err := s.store.GetMessageData(ctx, model.Uuid(req.Id()))
	if err != nil {
		return fmt.Errorf("failed to get message from store: %v", err)
	}

	data, err := crypto.Decrypt(key, edata)
	if err != nil {
		return fmt.Errorf("failed to decrypt member data: %v", err)
	}
	prev := model.GetRootAsMessage(data, 0)

	next := model.CloneMessageWithUpdates(prev, string(req.Content()))
	meta.UpdatedAt = next.Updated()

	edata, err = crypto.Encrypt(key, next.Table().Bytes)
	if err != nil {
		return fmt.Errorf("failed to encrypt updated message: %v", err)
	}

	err = s.store.UpdateMessageData(ctx, meta, edata)
	if err != nil {
		return fmt.Errorf("failed to store updated message: %v", err)
	}

	return nil
}

func (s *MessageService) Remove(ctx context.Context, req *MessageRemoveRequest) error {
	err := s.store.RemoveMessageData(ctx, model.Uuid(req.Id()))
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

	edatas, err := s.store.ListMessageData(ctx, model.Uuid(req.Conversation()), query)
	if err != nil {
		return nil, fmt.Errorf("failed to get message list form store: %v", err)
	}

	mlist := make([]*model.Message, len(edatas))
	for _, edata := range edatas {
		data, err := crypto.Decrypt(key, edata)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt message in list: %v", err)
		}

		mlist = append(mlist, model.GetRootAsMessage(data, 0))
	}

	if req.Pattern() == nil {
		return mlist, nil
	}

	// We have a content pattern, so we need to filter the list further
	r, err := regexp.Compile(string(req.Pattern()))
	if err != nil {
		return nil, fmt.Errorf("invalid content pattern: %v", err)
	}

	flist := make([]*model.Message, 0)
	for _, v := range mlist {
		if r.Match(v.Content()) {
			flist = append(flist, v)
		}
	}

	return flist, nil
}

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

type ConversationService struct {
	store store.ConversationStore
}

func NewConversationService(store store.ConversationStore) ConversationService {
	return ConversationService{store}
}

func (s *ConversationService) Add(
	ctx context.Context, req *ConversationAddRequest, key crypto.Key,
) (*model.Conversation, error) {
	// Create the conversation object
	mods := make([]model.Uuid, req.ModeratorsLength())
	for i := range req.ModeratorsLength() {
		mods = append(mods, model.Uuid(req.Moderators(i)))
	}
	c, err := model.NewConversation(string(req.Name()), string(req.Description()), mods)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation object: %v", err)
	}

	// Encrypt the conversation info before storing it in the database
	edata, err := crypto.Encrypt(key, c.Table().Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create encrypted data accessor: %v", err)
	}

	// Store the conversation in the database
	meta := store.ConversationMetadata{
		Id:        model.Uuid(c.Id()),
		CreatedAt: c.Created(),
		UpdatedAt: c.Updated(),
	}
	err = s.store.AddConversationData(ctx, meta, edata)
	if err != nil {
		return nil, fmt.Errorf("failed to store conversation in database: %v", err)
	}

	return c, nil
}

func (s *ConversationService) Get(
	ctx context.Context, req *ConversationGetRequest, key crypto.Key,
) (*model.Conversation, error) {
	_, edata, err := s.store.GetConversationData(ctx, model.Uuid(req.Id()))
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation from store: %v", err)
	}

	data, err := crypto.Decrypt(key, edata)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt conversation: %v", err)
	}

	return model.GetRootAsConversation(data, 0), nil
}

func (s *ConversationService) Update(
	ctx context.Context, req *ConversationUpdateRequest, key crypto.Key,
) error {
	meta, edata, err := s.store.GetConversationData(ctx, model.Uuid(req.Id()))
	if err != nil {
		return fmt.Errorf("failed to get conversation from store: %v", err)
	}

	data, err := crypto.Decrypt(key, edata)
	if err != nil {
		return fmt.Errorf("failed to decrypt conversation: %v", err)
	}

	prev := model.GetRootAsConversation(data, 0)
	c := model.CloneConversationWithUpdates(prev, req.Name(), req.Description(), nil)

	edata, err = crypto.Encrypt(key, c.Table().Bytes)
	if err != nil {
		return fmt.Errorf("failed to encrypt conversation: %v", err)
	}

	meta.UpdatedAt = c.Updated()

	err = s.store.UpdateConversationData(ctx, meta, edata)
	if err != nil {
		return fmt.Errorf("faled to store updated conversation: %v", err)
	}

	return nil
}

func (s *ConversationService) Remove(
	ctx context.Context, req *ConversationRemoveRequest,
) error {
	err := s.store.RemoveConversationData(ctx, model.Uuid(req.Id()))
	if err != nil {
		return fmt.Errorf("failed to remove member from database: %v", err)
	}

	return nil
}

func (s *ConversationService) ListAll(
	ctx context.Context, key crypto.Key,
) ([]*model.Conversation, error) {
	edatas, err := s.store.ListConversationData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation list: %v", err)
	}

	cs := make([]*model.Conversation, 0)
	for i := range edatas {
		data, err := crypto.Decrypt(key, edatas[i])
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt conversation in list: %v", err)
		}
		cs = append(cs, model.GetRootAsConversation(data, 0))
	}

	return cs, nil
}

func (s *ConversationService) AddMods(
	ctx context.Context, req *ConversationModsAddRequest, key crypto.Key,
) error {
	meta, edata, err := s.store.GetConversationData(ctx, model.Uuid(req.Id()))
	if err != nil {
		return fmt.Errorf("failed to get conversation from store: %v", err)
	}

	data, err := crypto.Decrypt(key, edata)
	if err != nil {
		return fmt.Errorf("failed to decrypt conversation: %v", err)
	}

	prev := model.GetRootAsConversation(data, 0)

	prevMods := make(map[string]bool, prev.ModsLength())
	newMods := make([][]byte, prev.ModsLength())
	for i := range prev.ModsLength() {
		prevMods[string(prev.Mods(i))] = true
		newMods = append(newMods, prev.Mods(i))
	}

	for i := range req.ModeratorsLength() {
		_, ok := prevMods[string(req.Moderators(i))]
		if !ok {
			newMods = append(newMods, req.Moderators(i))
		}
	}

	c := model.CloneConversationWithUpdates(prev, nil, nil, newMods)
	meta.UpdatedAt = c.Updated()

	edata, err = crypto.Encrypt(key, c.Table().Bytes)
	if err != nil {
		return fmt.Errorf("failed to encrypt updated conversation: %v", err)
	}

	err = s.store.UpdateConversationData(ctx, meta, edata)
	if err != nil {
		return fmt.Errorf("failed to store updated conversation: %v", err)
	}

	return nil
}

func (s *ConversationService) RemoveMods(
	ctx context.Context, req *ConversationModsRemoveRequest, key crypto.Key,
) error {
	meta, edata, err := s.store.GetConversationData(ctx, model.Uuid(req.Id()))
	if err != nil {
		return fmt.Errorf("failed to get conversation from store: %v", err)
	}

	data, err := crypto.Decrypt(key, edata)
	if err != nil {
		return fmt.Errorf("failed to decrypt conversation: %v", err)
	}

	prev := model.GetRootAsConversation(data, 0)

	modsToRemove := make(map[string]bool, req.ModeratorsLength())
	for i := range req.ModeratorsLength() {
		modsToRemove[string(req.Moderators(i))] = true
	}

	newMods := make([][]byte, prev.ModsLength())
	for i := range prev.ModsLength() {
		_, found := modsToRemove[string(prev.Mods(i))]
		if !found {
			newMods = append(newMods, prev.Mods(i))
		}
	}

	if len(newMods) == 0 {
		return fmt.Errorf("removing requested moderators would result in no moderators")
	}

	c := model.CloneConversationWithUpdates(prev, nil, nil, newMods)
	meta.UpdatedAt = c.Updated()

	edata, err = crypto.Encrypt(key, c.Table().Bytes)
	if err != nil {
		return fmt.Errorf("failed to encrypt updated conversation: %v", err)
	}

	err = s.store.UpdateConversationData(ctx, meta, edata)
	if err != nil {
		return fmt.Errorf("failed to store updated conversation: %v", err)
	}

	return nil
}

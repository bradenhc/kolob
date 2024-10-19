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
	mods := make([]model.Uuid, 0, req.ModeratorsLength())
	for i := range req.ModeratorsLength() {
		mods = append(mods, model.Uuid(req.Moderators(i)))
	}
	c, err := model.NewConversation(string(req.Name()), string(req.Description()), mods)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation object: %v", err)
	}

	// Create the entity we will store
	entity, err := store.NewConversationEntity(c, key)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation entity: %v", err)
	}

	// Store the entity
	err = s.store.AddConversationEntity(ctx, entity)
	if err != nil {
		return nil, fmt.Errorf("failed to store conversation entity: %v", err)
	}

	return c, nil
}

func (s *ConversationService) Get(
	ctx context.Context, req *ConversationGetRequest, key crypto.Key,
) (*model.Conversation, error) {
	entity, err := s.store.GetConversationEntity(ctx, model.Uuid(req.Id()))
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation from store: %v", err)
	}

	m, err := entity.Decrypt(key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt conversation: %v", err)
	}

	return m, nil
}

func (s *ConversationService) Update(
	ctx context.Context, req *ConversationUpdateRequest, key crypto.Key,
) (*model.Conversation, error) {
	entity, err := s.store.GetConversationEntity(ctx, model.Uuid(req.Id()))
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation from store: %v", err)
	}

	convo, err := entity.Update(key, req.Name(), req.Description(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to update conversation entity: %v", err)
	}

	err = s.store.UpdateConversationEntity(ctx, entity)
	if err != nil {
		return nil, fmt.Errorf("faled to store updated conversation entity: %v", err)
	}

	return convo, nil
}

func (s *ConversationService) Remove(
	ctx context.Context, req *ConversationRemoveRequest,
) error {
	err := s.store.RemoveConversationEntity(ctx, model.Uuid(req.Id()))
	if err != nil {
		return fmt.Errorf("failed to remove member from database: %v", err)
	}

	return nil
}

func (s *ConversationService) ListAll(
	ctx context.Context, key crypto.Key,
) ([]*model.Conversation, error) {
	edatas, err := s.store.ListConversationEntities(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation entity list: %v", err)
	}

	cs := make([]*model.Conversation, 0, len(edatas))
	for _, e := range edatas {
		c, err := e.Decrypt(key)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt conversation in list: %v", err)
		}
		cs = append(cs, c)
	}

	return cs, nil
}

func (s *ConversationService) AddMods(
	ctx context.Context, req *ConversationModsAddRequest, key crypto.Key,
) error {
	entity, err := s.store.GetConversationEntity(ctx, model.Uuid(req.Id()))
	if err != nil {
		return fmt.Errorf("failed to get conversation from store: %v", err)
	}

	prev, err := entity.Decrypt(key)
	if err != nil {
		return fmt.Errorf("failed to descrypt conversation: %v", err)
	}

	// First, copy the existing moderator entries into the new moderator list. While doing so, keep
	// track of which moderator ids are already in the list so that we don't add duplicates.
	prevMods := make(map[string]bool, prev.ModsLength())
	newMods := make([][]byte, 0, prev.ModsLength())
	for i := range prev.ModsLength() {
		prevMods[string(prev.Mods(i))] = true
		newMods = append(newMods, prev.Mods(i))
	}

	// Add new moderators to the list if their id isn't already there.
	for i := range req.ModeratorsLength() {
		_, ok := prevMods[string(req.Moderators(i))]
		if !ok {
			newMods = append(newMods, req.Moderators(i))
		}
	}

	c := model.CloneConversationWithUpdates(prev, nil, nil, newMods)
	entity, err = store.NewConversationEntity(c, key)
	if err != nil {
		return fmt.Errorf("failed to create updated conversation entity: %v", err)
	}

	err = s.store.UpdateConversationEntity(ctx, entity)
	if err != nil {
		return fmt.Errorf("failed to store updated conversation: %v", err)
	}

	return nil
}

func (s *ConversationService) RemoveMods(
	ctx context.Context, req *ConversationModsRemoveRequest, key crypto.Key,
) error {
	entity, err := s.store.GetConversationEntity(ctx, model.Uuid(req.Id()))
	if err != nil {
		return fmt.Errorf("failed to get conversation from store: %v", err)
	}

	prev, err := entity.Decrypt(key)
	if err != nil {
		return fmt.Errorf("failed to decrypt conversation: %v", err)
	}

	// First, convert the list of mods to remove to a map so that we can easily test for ids
	modsToRemove := make(map[string]bool, req.ModeratorsLength())
	for i := range req.ModeratorsLength() {
		modsToRemove[string(req.Moderators(i))] = true
	}

	// Iterate over the existing mods. If the id exists in the mapping of mods to remove, then
	// remove it.
	newMods := make([][]byte, 0, prev.ModsLength())
	for i := range prev.ModsLength() {
		_, found := modsToRemove[string(prev.Mods(i))]
		if !found {
			newMods = append(newMods, prev.Mods(i))
		}
	}

	if len(newMods) == 0 {
		return fmt.Errorf("removing requested moderators would result in no moderators")
	}

	// Create the updated entity with the new moderator list
	c := model.CloneConversationWithUpdates(prev, nil, nil, newMods)
	entity, err = store.NewConversationEntity(c, key)
	if err != nil {
		return fmt.Errorf("failed to create update conversation entity: %v", err)
	}

	// Store the updated entity
	err = s.store.UpdateConversationEntity(ctx, entity)
	if err != nil {
		return fmt.Errorf("failed to store updated conversation: %v", err)
	}

	return nil
}

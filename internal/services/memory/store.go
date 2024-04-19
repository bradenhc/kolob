// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package memory

import (
	"sync"

	"github.com/bradenhc/kolob/internal/model"
)

type SynchronizedStore struct {
	members           []*model.Member
	membersLock       sync.RWMutex
	conversations     map[model.Uuid]*model.Conversation
	conversationsLock sync.RWMutex
	messages          map[model.Uuid][]*model.Message
	messagesLock      sync.RWMutex
}

func NewSynchrnoizedStore() *SynchronizedStore {
	return &SynchronizedStore{
		members:       make([]*model.Member, 0),
		conversations: make(map[model.Uuid]*model.Conversation),
		messages:      make(map[model.Uuid][]*model.Message),
	}
}

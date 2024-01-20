package memory

import (
	"sync"

	"github.com/bradenhc/kolob/internal/model"
)

type SynchronizedStore struct {
	lock          sync.RWMutex
	groups        map[model.Uuid]*model.Group
	members       map[model.Uuid][]*model.Member
	conversations map[model.Uuid][]*model.Conversation
	messages      map[model.Uuid][]*model.Message
}

func NewSynchrnoizedStore() *SynchronizedStore {
	return &SynchronizedStore{
		groups:        make(map[model.Uuid]*model.Group),
		members:       make(map[model.Uuid][]*model.Member),
		conversations: make(map[model.Uuid][]*model.Conversation),
		messages:      make(map[model.Uuid][]*model.Message),
	}
}

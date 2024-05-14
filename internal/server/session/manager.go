// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package session

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/bradenhc/kolob/internal/crypto"
	"github.com/bradenhc/kolob/internal/model"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
)

type session struct {
	pkey crypto.Key
	last time.Time
}

type Manager struct {
	sessions   map[model.Uuid]session
	sessionsmx sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		sessions: make(map[model.Uuid]session),
	}
}

func (m *Manager) Add(k crypto.Key) (model.Uuid, error) {
	id, err := model.NewUuid()
	if err != nil {
		return "", fmt.Errorf("failed to generate session ID: %v", err)
	}
	m.sessionsmx.Lock()
	defer m.sessionsmx.Unlock()
	m.sessions[id] = session{pkey: k, last: time.Now()}
	return id, nil
}

func (m *Manager) Get(id model.Uuid) (crypto.Key, error) {
	m.sessionsmx.Lock()
	defer m.sessionsmx.Unlock()

	s, ok := m.sessions[id]
	if !ok {
		return nil, ErrSessionNotFound
	}

	if time.Now().After(s.last.Add(15 * time.Minute)) {
		delete(m.sessions, id)
		return nil, ErrSessionExpired
	}

	s.last = time.Now()

	return s.pkey, nil
}

func (m *Manager) Remove(id model.Uuid) {
	m.sessionsmx.Lock()
	defer m.sessionsmx.Unlock()

	delete(m.sessions, id)
}

// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/bradenhc/kolob/internal/crypto"
	"github.com/bradenhc/kolob/internal/model"
)

const (
	SessionCookie  = "sessionid"
	SessionPassKey = ContextKey("pkey")
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
)

type session struct {
	pkey crypto.Key
	last time.Time
}

type SessionManager struct {
	sessions map[model.Uuid]session
}

func NewSessionManager() SessionManager {
	return SessionManager{
		sessions: make(map[model.Uuid]session),
	}
}

func (m *SessionManager) Add(k crypto.Key) (model.Uuid, error) {
	id, err := model.NewUuid()
	if err != nil {
		return "", fmt.Errorf("failed to generate session ID: %v", err)
	}
	m.sessions[id] = session{pkey: k, last: time.Now()}
	return id, nil
}

func (m *SessionManager) Get(id model.Uuid) (crypto.Key, error) {
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

func (m *SessionManager) Remove(id model.Uuid) {
	delete(m.sessions, id)
}

func (s *SessionManager) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(SessionCookie)
		if err != nil {
			switch err {
			case http.ErrNoCookie:
				w.WriteHeader(http.StatusUnauthorized)
			default:
				w.WriteHeader(http.StatusBadRequest)
			}
			return
		}

		pkey, err := s.Get(model.Uuid(c.Value))
		if err != nil {
			switch err {
			case ErrSessionNotFound:
				w.WriteHeader(http.StatusBadRequest)
			case ErrSessionExpired:
				w.WriteHeader(http.StatusUnauthorized)
			}
			return
		}

		next(w, r.WithContext(context.WithValue(r.Context(), SessionPassKey, pkey)))
	}
}

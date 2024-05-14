// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package session

import (
	"net/http"

	"github.com/bradenhc/kolob/internal/model"
)

const (
	cookieName = "sessionid"
)

func (s *Manager) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(cookieName)
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

		next(w, r.WithContext(NewContext(r.Context(), pkey)))
	}
}

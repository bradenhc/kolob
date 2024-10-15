// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package server

import (
	"net/http"

	"github.com/bradenhc/kolob/internal/server/session"
	"github.com/bradenhc/kolob/internal/services"
)

type GroupHandler struct {
	groups services.GroupService
}

func NewGroupHandler(gs services.GroupService) GroupHandler {
	return GroupHandler{gs}
}

func (h *GroupHandler) InitGroup(w http.ResponseWriter, r *http.Request) {
	body := make([]byte, 1024)
	r.Body.Read(body)
	req := services.GetRootAsGroupInitRequest(body, 0)
	g, err := h.groups.Create(r.Context(), req)
	if err != nil {
		WriteJsonErr(w, http.StatusInternalServerError, err)
		return
	}

	WriteJson(w, http.StatusOK, g)
}

func (h *GroupHandler) GetGroupInfo(w http.ResponseWriter, r *http.Request) {
	_, err := session.FromContext(r.Context())
	if err != nil {
		WriteJsonErr(w, http.StatusInternalServerError, err)
		return
	}
}

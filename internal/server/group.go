// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package server

import (
	"encoding/json"
	"net/http"

	"github.com/bradenhc/kolob/internal/model"
	"github.com/bradenhc/kolob/internal/server/session"
)

type GroupHandler struct {
	groups model.GroupService
}

func NewGroupHandler(gs model.GroupService) GroupHandler {
	return GroupHandler{gs}
}

func (h *GroupHandler) InitGroup(w http.ResponseWriter, r *http.Request) {
	var p model.InitGroupParams
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		WriteJsonErr(w, http.StatusBadRequest, err)
		return
	}

	g, err := h.groups.InitGroup(r.Context(), p)
	if err != nil {
		WriteJsonErr(w, http.StatusInternalServerError, err)
		return
	}

	WriteJson(w, http.StatusOK, g)
}

func (h *GroupHandler) GetGroupInfo(w http.ResponseWriter, r *http.Request) {
	pkey, err := session.FromContext(r.Context())
	if err != nil {
		WriteJsonErr(w, http.StatusInternalServerError, err)
		return
	}

	g, err := h.groups.GetGroupInfo(r.Context(), model.GetGroupInfoParams{
		Id:      model.Uuid(r.PathValue("id")),
		PassKey: pkey,
	})
	if err != nil {
		WriteJsonErr(w, http.StatusInternalServerError, err)
	}

	WriteJson(w, http.StatusOK, g)
}

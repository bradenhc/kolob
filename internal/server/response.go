// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type eresponse struct {
	Error string `json:"error"`
}

func WriteJsonErr(w http.ResponseWriter, status int, err error) {
	WriteJson(w, status, eresponse{err.Error()})
}

func WriteJson(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("Failed to write JSON response", "err", err.Error())
	}
}

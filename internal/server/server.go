// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bradenhc/kolob/internal/services/sqlite"
)

type ContextKey string

type Server struct {
	sessions     SessionManager
	db           *sql.DB
	groupHandler GroupHandler
	httpServer   *http.Server
}

func NewServer(c Config) (*Server, error) {
	db, err := sqlite.Open(c.DatabaseFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	groupService := sqlite.NewGroupService(db)
	groupHandler := NewGroupHandler(&groupService)

	sessions := NewSessionManager()

	middlware := NewMiddlewareChain(&sessions)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/group", middlware.Finish(groupHandler.InitGroup))

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%d", c.Port),
		Handler: mux,
	}

	server := &Server{
		sessions, db, groupHandler, &httpServer,
	}

	return server, nil
}

func (s *Server) Start() {
	go func() {
		if err := s.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server error", "err", err.Error())
		}
		slog.Info("Stopped serving new connections")
	}()

	slog.Info("HTTP server started: waiting for terminating signal")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, release := context.WithTimeout(context.Background(), 10*time.Second)
	defer release()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		slog.Error("HTTP shutdown error", "err", err.Error())
	}

	slog.Info("Kolob server shut down successfully")
}

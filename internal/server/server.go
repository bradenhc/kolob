// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package server

import (
	"context"
	"crypto/tls"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bradenhc/kolob/internal/appfs"
	"github.com/bradenhc/kolob/internal/crypto"
	"github.com/bradenhc/kolob/internal/services"
	"github.com/bradenhc/kolob/internal/session"
	"github.com/bradenhc/kolob/internal/store/sqlite"
)

type ContextKey string

type Server struct {
	sessions     *session.Manager
	db           *sql.DB
	groupHandler GroupHandler
	httpServer   *http.Server
}

func NewServer(c Config) (*Server, error) {
	// First, we need to create a self-signed TLS configuration to use with our server. This will
	// guarantee that we force all connections to be encrypted and can use HTTP/2. The user can
	// override the self-signed certificate with their own if they choose using the provided
	// configuration object.
	slog.Info("Generating self-signed TLS configuration")
	tlsConfig, err := createSelfSignedTlsConfig()
	if err != nil {
		return nil, err
	}

	slog.Info("Openning database")
	db, err := sqlite.Open(c.DatabaseFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	slog.Info("Creating database stores")
	groupStore, err := sqlite.NewGroupStore(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create group store: %v", err)
	}
	groupService := services.NewGroupService(groupStore)
	groupHandler := NewGroupHandler(groupService)

	sessions := session.NewManager()

	middlware := NewMiddlewareChain(sessions)

	slog.Info("Registering routes")
	mux := http.NewServeMux()

	// Setup the static filer server for the web UI
	fsys := fs.FS(appfs.AppFS)
	webui, _ := fs.Sub(fsys, "webui")
	mux.Handle("GET /", http.FileServer(http.FS(webui)))

	// Setup the routes for the API
	mux.HandleFunc("POST /api/v1/group", groupHandler.InitGroup)
	mux.HandleFunc("GET /api/v1/group", middlware.Finish(groupHandler.GetGroupInfo))

	slog.Info("Creating HTTP server")
	httpServer := http.Server{
		Addr:      fmt.Sprintf(":%d", c.Port),
		Handler:   mux,
		TLSConfig: tlsConfig,
	}

	server := &Server{
		sessions, db, groupHandler, &httpServer,
	}

	return server, nil
}

func (s *Server) Start() {
	go func() {
		if err := s.httpServer.ListenAndServeTLS("", ""); !errors.Is(err, http.ErrServerClosed) {
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

func createSelfSignedTlsConfig() (*tls.Config, error) {
	crt, key, err := crypto.GenerateSelfSignedCert()
	if err != nil {
		return nil, fmt.Errorf("failed to generate self-signed certificate: %v", err)
	}

	cert, err := tls.X509KeyPair(crt, key)
	if err != nil {
		return nil, fmt.Errorf("failed to load X509 key pair: %v", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "localhost",
	}, nil
}

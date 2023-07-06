package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/hiring-seedtag/mihai-lupoiu-go-backend-test/internal/core/services"
)

type ServerHTTP struct {
	svc *services.EndorService
	srv *http.Server
	h   *HandlerHTTP
}

func NewHTTPServer(endorService *services.EndorService, validate *validator.Validate) *ServerHTTP {
	handler := &HandlerHTTP{
		svc: endorService,
		v:   validate,
		r:   chi.NewRouter(),
	}

	handler.configureRoutes()

	server := &http.Server{
		Addr:         ":3000",           // configure the bind address with default port
		Handler:      handler.r,         // set the default handler
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	httpServer := &ServerHTTP{
		svc: endorService,
		srv: server,
		h:   handler,
	}

	return httpServer
}

// ListenAndServe starts the HTTP server and listens on the specified port.
func (s *ServerHTTP) ListenAndServe(port string) error {
	if len(port) != 0 {
		s.srv.Addr = port
	}
	return s.srv.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server.
func (s *ServerHTTP) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

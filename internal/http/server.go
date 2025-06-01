package http

import (
	"goat/internal"

	_ "goat/docs"

	"context"
	"net/http"
	"strconv"
	"time"
)

// ShutdownTimeout is the time given for outstanding requests to finish before shutdown.
const ShutdownTimeout = 1 * time.Second

// Server represents an HTTP server. It is meant to wrap all HTTP functionality
// used by the application so that dependent packages (such as cmd/main) do not
// need to reference the "net/http" package at all.
type Server struct {
	server *http.Server

	// HTTP Mux for handling HTTP communication.
	Router *http.ServeMux

	// Bind address & domain for the server's listener.
	// If domain is specified, server is run on TLS using acme/autocert.
	Port int

	// Services used by the various HTTP routes.
	UserService internal.UserService

	BlacklistedToken map[string]bool
}

// NewServer returns a new instance of Server.
func NewServer(port int) *Server {
	// Initiate Router
	router := http.NewServeMux()

	// Create a new server that wraps the net/http server & add a gorilla router.
	server := &Server{
		Port: port,
		server: &http.Server{
			Addr:    ":" + strconv.Itoa(port),
			Handler: router,
		},
		BlacklistedToken: make(map[string]bool),
	}

	// ?
	server.Router = router

	// Setup Base Routes
	server.registerAuthRoutes(router)
	server.registerUserRoutes(router)

	// Load Swagger Doc
	server.loadSwagger()

	return server
}

// Close gracefully shuts down the server.
func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}

// ListenAndServeTLSRedirect runs an HTTP server on port 80 to redirect users
// to the TLS-enabled port 443 server.
func (s *Server) ListenAndServe() error {
	internal.Debug("Server::ListenAndServe", "Listening on port "+strconv.Itoa(s.Port))
	return s.server.ListenAndServe()
}

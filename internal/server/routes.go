package server

// API route wiring lives here so `server.go` stays the single mux owner.

import (
	"net/http"

	"github.com/Mauricio-AFLadeira/taskforge-api-go/internal/auth"
	"github.com/Mauricio-AFLadeira/taskforge-api-go/internal/organizations"
	"github.com/Mauricio-AFLadeira/taskforge-api-go/internal/middleware"
)

func (s *Server) mountAPIRoutes(mux *http.ServeMux) {
	h := auth.NewHandler(s.pool, s.cfg)
	h.RegisterAuth(mux)
	mux.Handle("/me", middleware.RequireAuth(s.cfg.JWTSecret, http.HandlerFunc(h.Me)))
	orgs := organizations.NewHandler(s.pool, s.cfg)
	orgs.RegisterRoutes(mux, s.cfg.JWTSecret)
}

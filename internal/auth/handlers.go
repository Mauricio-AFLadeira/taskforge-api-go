package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Mauricio-AFLadeira/taskforge-api-go/internal/config"
	"github.com/Mauricio-AFLadeira/taskforge-api-go/internal/shared"
	"github.com/Mauricio-AFLadeira/taskforge-api-go/internal/users"
)

// Handler exposes HTTP endpoints for authentication.
type Handler struct {
	svc *Service
}

// NewHandler wires repositories and token issuance from configuration.
func NewHandler(pool *pgxpool.Pool, cfg config.Config) *Handler {
	userRepo := users.NewRepository(pool)
	refreshRepo := NewRefreshRepository(pool)
	tokens := NewTokenIssuer(cfg.JWTSecret, cfg.JWTAccessDuration)
	svc := NewService(userRepo, refreshRepo, tokens, cfg.RefreshTokenDuration)
	return &Handler{svc: svc}
}

// RegisterAuth mounts POST /auth/register|login|refresh|logout on mux.
func (h *Handler) RegisterAuth(mux *http.ServeMux) {
	mux.HandleFunc("/auth/register", h.register)
	mux.HandleFunc("/auth/login", h.login)
	mux.HandleFunc("/auth/refresh", h.refresh)
	mux.HandleFunc("/auth/logout", h.logout)
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w, http.MethodPost)
		return
	}
	var req RegisterRequest
	if err := decodeJSON(w, r, &req); err != nil {
		return
	}
	out, err := h.svc.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		writeSvcErr(w, err)
		return
	}
	if err := shared.WriteJSON(w, http.StatusCreated, out); err != nil {
		return
	}
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w, http.MethodPost)
		return
	}
	var req LoginRequest
	if err := decodeJSON(w, r, &req); err != nil {
		return
	}
	out, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		writeSvcErr(w, err)
		return
	}
	_ = shared.WriteJSON(w, http.StatusOK, out)
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w, http.MethodPost)
		return
	}
	var req RefreshRequest
	if err := decodeJSON(w, r, &req); err != nil {
		return
	}
	out, err := h.svc.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		writeSvcErr(w, err)
		return
	}
	_ = shared.WriteJSON(w, http.StatusOK, out)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w, http.MethodPost)
		return
	}
	var req LogoutRequest
	if err := decodeJSON(w, r, &req); err != nil {
		return
	}
	if err := h.svc.Logout(r.Context(), req.RefreshToken); err != nil {
		writeSvcErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Me returns the authenticated user (requires JWT middleware).
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, http.MethodGet)
		return
	}
	userID, ok := shared.UserIDFromContext(r.Context())
	if !ok {
		_ = shared.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	up, err := h.svc.Me(r.Context(), userID)
	if err != nil {
		if errors.Is(err, users.ErrNotFound) {
			_ = shared.WriteError(w, http.StatusUnauthorized, "user not found")
			return
		}
		_ = shared.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}
	_ = shared.WriteJSON(w, http.StatusOK, up)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		_ = shared.WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return err
	}
	return nil
}

func writeMethodNotAllowed(w http.ResponseWriter, allow string) {
	w.Header().Set("Allow", allow)
	_ = shared.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
}

func writeSvcErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrValidation):
		_ = shared.WriteError(w, http.StatusBadRequest, "validation failed")
	case errors.Is(err, users.ErrEmailTaken):
		_ = shared.WriteError(w, http.StatusConflict, "email already registered")
	case errors.Is(err, ErrInvalidCredentials):
		_ = shared.WriteError(w, http.StatusUnauthorized, "invalid credentials")
	case errors.Is(err, ErrInvalidRefresh):
		_ = shared.WriteError(w, http.StatusUnauthorized, "invalid refresh token")
	default:
		_ = shared.WriteError(w, http.StatusInternalServerError, "internal error")
	}
}

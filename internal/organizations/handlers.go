package organizations

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mauricio-reportei/taskforge-api-go/internal/config"
	"github.com/mauricio-reportei/taskforge-api-go/internal/middleware"
	"github.com/mauricio-reportei/taskforge-api-go/internal/shared"
	"github.com/mauricio-reportei/taskforge-api-go/internal/users"
)

type Handler struct {
	svc *Service
}

func NewHandler(pool *pgxpool.Pool, cfg config.Config) *Handler {
	userRepo := users.NewRepository(pool)
	repo := NewRepository(pool)
	svc := NewService(userRepo, repo)
	_ = cfg
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, jwtSecret string) {
	mux.Handle("/organizations", middleware.RequireAuth(jwtSecret, http.HandlerFunc(h.organizationsCollection)))
	mux.Handle("/organizations/", middleware.RequireAuth(jwtSecret, http.HandlerFunc(h.organizationsSubtree)))
}

func (h *Handler) organizationsCollection(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/organizations" {
		notFound(w)
		return
	}
	userID, ok := shared.UserIDFromContext(r.Context())
	if !ok {
		_ = shared.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	switch r.Method {
	case http.MethodPost:
		var req CreateOrganizationRequest
		if err := decodeJSON(w, r, &req); err != nil {
			return
		}
		out, err := h.svc.CreateOrganization(r.Context(), userID, req.Name)
		if err != nil {
			writeErr(w, err)
			return
		}
		_ = shared.WriteJSON(w, http.StatusCreated, out)
	case http.MethodGet:
		out, err := h.svc.ListOrganizations(r.Context(), userID)
		if err != nil {
			writeErr(w, err)
			return
		}
		_ = shared.WriteJSON(w, http.StatusOK, out)
	default:
		writeMethodNotAllowed(w, http.MethodGet, http.MethodPost)
	}
}

func (h *Handler) organizationsSubtree(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/organizations/")
	if path == "" {
		notFound(w)
		return
	}
	parts := strings.Split(path, "/")
	if len(parts) == 0 || parts[0] == "" {
		notFound(w)
		return
	}
	orgID := parts[0]
	userID, ok := shared.UserIDFromContext(r.Context())
	if !ok {
		_ = shared.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if len(parts) == 1 {
		h.organizationItem(w, r, userID, orgID)
		return
	}

	if parts[1] != "members" {
		notFound(w)
		return
	}

	switch {
	case len(parts) == 2:
		h.organizationMembersCollection(w, r, userID, orgID)
	case len(parts) == 3:
		h.organizationMemberItem(w, r, userID, orgID, parts[2])
	case len(parts) == 4 && parts[3] == "role":
		h.organizationMemberRole(w, r, userID, orgID, parts[2])
	default:
		notFound(w)
	}
}

func (h *Handler) organizationItem(w http.ResponseWriter, r *http.Request, userID, orgID string) {
	switch r.Method {
	case http.MethodGet:
		out, err := h.svc.GetOrganization(r.Context(), userID, orgID)
		if err != nil {
			writeErr(w, err)
			return
		}
		_ = shared.WriteJSON(w, http.StatusOK, out)
	case http.MethodPatch:
		var req UpdateOrganizationRequest
		if err := decodeJSON(w, r, &req); err != nil {
			return
		}
		out, err := h.svc.UpdateOrganization(r.Context(), userID, orgID, req.Name)
		if err != nil {
			writeErr(w, err)
			return
		}
		_ = shared.WriteJSON(w, http.StatusOK, out)
	case http.MethodDelete:
		if err := h.svc.DeleteOrganization(r.Context(), userID, orgID); err != nil {
			writeErr(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		writeMethodNotAllowed(w, http.MethodDelete, http.MethodGet, http.MethodPatch)
	}
}

func (h *Handler) organizationMembersCollection(w http.ResponseWriter, r *http.Request, userID, orgID string) {
	switch r.Method {
	case http.MethodGet:
		out, err := h.svc.ListMembers(r.Context(), userID, orgID)
		if err != nil {
			writeErr(w, err)
			return
		}
		_ = shared.WriteJSON(w, http.StatusOK, out)
	case http.MethodPost:
		var req AddMemberRequest
		if err := decodeJSON(w, r, &req); err != nil {
			return
		}
		out, err := h.svc.AddMember(r.Context(), userID, orgID, req.Email, req.Role)
		if err != nil {
			writeErr(w, err)
			return
		}
		_ = shared.WriteJSON(w, http.StatusCreated, out)
	default:
		writeMethodNotAllowed(w, http.MethodGet, http.MethodPost)
	}
}

func (h *Handler) organizationMemberItem(w http.ResponseWriter, r *http.Request, userID, orgID, memberID string) {
	if r.Method != http.MethodDelete {
		writeMethodNotAllowed(w, http.MethodDelete)
		return
	}
	if err := h.svc.RemoveMember(r.Context(), userID, orgID, memberID); err != nil {
		writeErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) organizationMemberRole(w http.ResponseWriter, r *http.Request, userID, orgID, memberID string) {
	if r.Method != http.MethodPatch {
		writeMethodNotAllowed(w, http.MethodPatch)
		return
	}
	var req ChangeMemberRoleRequest
	if err := decodeJSON(w, r, &req); err != nil {
		return
	}
	out, err := h.svc.ChangeMemberRole(r.Context(), userID, orgID, memberID, req.Role)
	if err != nil {
		writeErr(w, err)
		return
	}
	_ = shared.WriteJSON(w, http.StatusOK, out)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		_ = shared.WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return err
	}
	return nil
}

func writeMethodNotAllowed(w http.ResponseWriter, allow ...string) {
	w.Header().Set("Allow", strings.Join(allow, ", "))
	_ = shared.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
}

func notFound(w http.ResponseWriter) {
	_ = shared.WriteError(w, http.StatusNotFound, "not found")
}

func writeErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrValidation):
		_ = shared.WriteError(w, http.StatusBadRequest, "validation failed")
	case errors.Is(err, ErrNotFound):
		_ = shared.WriteError(w, http.StatusNotFound, "not found")
	case errors.Is(err, ErrForbidden):
		_ = shared.WriteError(w, http.StatusForbidden, "forbidden")
	case errors.Is(err, ErrConflict):
		_ = shared.WriteError(w, http.StatusConflict, "conflict")
	default:
		_ = shared.WriteError(w, http.StatusInternalServerError, "internal error")
	}
}

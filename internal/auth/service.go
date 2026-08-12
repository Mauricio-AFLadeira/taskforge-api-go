package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"golang.org/x/crypto/bcrypt"

	"github.com/Mauricio-AFLadeira/taskforge-api-go/internal/users"
)

const (
	minPasswordLen = 8
	minNameLen     = 1
)

// Service orchestrates registration, login, refresh, and logout.
type Service struct {
	usersRepo   *users.Repository
	refreshRepo *RefreshRepository
	tokens      *TokenIssuer
	refreshTTL  time.Duration
}

// NewService constructs the auth service.
func NewService(usersRepo *users.Repository, refreshRepo *RefreshRepository, tokens *TokenIssuer, refreshTTL time.Duration) *Service {
	return &Service{
		usersRepo:   usersRepo,
		refreshRepo: refreshRepo,
		tokens:      tokens,
		refreshTTL:  refreshTTL,
	}
}

func publicUser(u users.User) UserPublic {
	return UserPublic{ID: u.ID, Name: u.Name, Email: u.Email}
}

func (s *Service) issueSession(ctx context.Context, u users.User) (*AuthResponse, error) {
	rawRefresh, refreshHash, err := NewRefreshOpaque()
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(s.refreshTTL)
	if err := s.refreshRepo.Insert(ctx, u.ID, refreshHash, expiresAt); err != nil {
		return nil, err
	}
	access, err := s.tokens.IssueAccess(u.ID, u.Email)
	if err != nil {
		return nil, err
	}
	return &AuthResponse{
		User:         publicUser(u),
		AccessToken:  access,
		RefreshToken: rawRefresh,
	}, nil
}

// Register validates input, hashes the password, and creates a user session.
func (s *Service) Register(ctx context.Context, name, email, password string) (*AuthResponse, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(strings.ToLower(email))
	if len(name) < minNameLen || email == "" || len(password) < minPasswordLen {
		return nil, ErrValidation
	}
	if !strings.Contains(email, "@") {
		return nil, ErrValidation
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u, err := s.usersRepo.Create(ctx, name, email, string(hash))
	if err != nil {
		if errors.Is(err, users.ErrEmailTaken) {
			return nil, users.ErrEmailTaken
		}
		return nil, err
	}

	return s.issueSession(ctx, u)
}

// Login verifies credentials and returns a new session.
func (s *Service) Login(ctx context.Context, email, password string) (*AuthResponse, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || password == "" {
		return nil, ErrInvalidCredentials
	}

	u, err := s.usersRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, users.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.issueSession(ctx, u)
}

// Refresh rotates the refresh token and returns new tokens (MVP rotation).
func (s *Service) Refresh(ctx context.Context, rawRefresh string) (*AuthResponse, error) {
	rawRefresh = strings.TrimSpace(rawRefresh)
	if rawRefresh == "" {
		return nil, ErrInvalidRefresh
	}

	hash := HashRefreshToken(rawRefresh)
	row, err := s.refreshRepo.FindByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidRefresh
		}
		return nil, err
	}

	now := time.Now()
	if row.RevokedAt != nil || !now.Before(row.ExpiresAt) {
		return nil, ErrInvalidRefresh
	}

	u, err := s.usersRepo.GetByID(ctx, row.UserID)
	if err != nil {
		return nil, err
	}

	newRaw, newHash, err := NewRefreshOpaque()
	if err != nil {
		return nil, err
	}
	newExpires := now.Add(s.refreshTTL)

	if err := s.refreshRepo.Rotate(ctx, row.ID, row.UserID, newHash, newExpires); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidRefresh
		}
		return nil, err
	}

	access, err := s.tokens.IssueAccess(u.ID, u.Email)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:         publicUser(u),
		AccessToken:  access,
		RefreshToken: newRaw,
	}, nil
}

// Logout revokes a refresh token (idempotent).
func (s *Service) Logout(ctx context.Context, rawRefresh string) error {
	rawRefresh = strings.TrimSpace(rawRefresh)
	if rawRefresh == "" {
		return ErrValidation
	}
	hash := HashRefreshToken(rawRefresh)
	return s.refreshRepo.RevokeByHash(ctx, hash)
}

// Me returns the current user without tokens.
func (s *Service) Me(ctx context.Context, userID string) (UserPublic, error) {
	u, err := s.usersRepo.GetByID(ctx, userID)
	if err != nil {
		return UserPublic{}, err
	}
	return publicUser(u), nil
}

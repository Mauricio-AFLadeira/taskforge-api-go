package auth

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RefreshRepository persists refresh token hashes.
type RefreshRepository struct {
	pool *pgxpool.Pool
}

// NewRefreshRepository constructs a repository backed by pgxpool.
func NewRefreshRepository(pool *pgxpool.Pool) *RefreshRepository {
	return &RefreshRepository{pool: pool}
}

// RefreshRow is one row from refresh_tokens.
type RefreshRow struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
	RevokedAt *time.Time
}

// Insert creates a refresh token row.
func (r *RefreshRepository) Insert(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`, userID, tokenHash, expiresAt)
	return err
}

// FindByHash loads a refresh token by hash (may be expired or revoked).
func (r *RefreshRepository) FindByHash(ctx context.Context, tokenHash string) (RefreshRow, error) {
	var row RefreshRow
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, expires_at, revoked_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`, tokenHash).Scan(&row.ID, &row.UserID, &row.ExpiresAt, &row.RevokedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return RefreshRow{}, pgx.ErrNoRows
		}
		return RefreshRow{}, err
	}
	return row, nil
}

// Rotate revokes the previous refresh token row and inserts a new one (same user).
func (r *RefreshRepository) Rotate(ctx context.Context, oldRowID, userID, newHash string, expiresAt time.Time) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	tag, err := tx.Exec(ctx, `
		UPDATE refresh_tokens
		SET revoked_at = now()
		WHERE id = $1 AND revoked_at IS NULL
	`, oldRowID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`, userID, newHash, expiresAt)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// RevokeByHash marks an active refresh token as revoked (no-op if unknown or already revoked).
func (r *RefreshRepository) RevokeByHash(ctx context.Context, tokenHash string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE refresh_tokens
		SET revoked_at = now()
		WHERE token_hash = $1 AND revoked_at IS NULL
	`, tokenHash)
	return err
}

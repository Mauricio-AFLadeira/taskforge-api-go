package organizations

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateOrganization(ctx context.Context, ownerID, name string) (Organization, error) {
	name = strings.TrimSpace(name)
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return Organization{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var org Organization
	err = tx.QueryRow(ctx, `
		INSERT INTO organizations (name, owner_id)
		VALUES ($1, $2)
		RETURNING id, name, owner_id, created_at, updated_at
	`, name, ownerID).Scan(&org.ID, &org.Name, &org.OwnerID, &org.CreatedAt, &org.UpdatedAt)
	if err != nil {
		return Organization{}, err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO organization_members (organization_id, user_id, role)
		VALUES ($1, $2, $3)
	`, org.ID, ownerID, RoleOwner)
	if err != nil {
		return Organization{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return Organization{}, err
	}
	return org, nil
}

func (r *Repository) ListOrganizationsForUser(ctx context.Context, userID string) ([]Organization, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT DISTINCT o.id, o.name, o.owner_id, o.created_at, o.updated_at
		FROM organizations o
		JOIN organization_members om ON om.organization_id = o.id
		WHERE om.user_id = $1
		ORDER BY o.created_at DESC, o.id DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orgs := make([]Organization, 0)
	for rows.Next() {
		var org Organization
		if err := rows.Scan(&org.ID, &org.Name, &org.OwnerID, &org.CreatedAt, &org.UpdatedAt); err != nil {
			return nil, err
		}
		orgs = append(orgs, org)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orgs, nil
}

func (r *Repository) GetOrganization(ctx context.Context, orgID string) (Organization, error) {
	var org Organization
	err := r.pool.QueryRow(ctx, `
		SELECT id, name, owner_id, created_at, updated_at
		FROM organizations
		WHERE id = $1
	`, orgID).Scan(&org.ID, &org.Name, &org.OwnerID, &org.CreatedAt, &org.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Organization{}, ErrNotFound
		}
		return Organization{}, err
	}
	return org, nil
}

func (r *Repository) UpdateOrganization(ctx context.Context, orgID, name string) (Organization, error) {
	name = strings.TrimSpace(name)
	var org Organization
	err := r.pool.QueryRow(ctx, `
		UPDATE organizations
		SET name = $2,
		    updated_at = now()
		WHERE id = $1
		RETURNING id, name, owner_id, created_at, updated_at
	`, orgID, name).Scan(&org.ID, &org.Name, &org.OwnerID, &org.CreatedAt, &org.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Organization{}, ErrNotFound
		}
		return Organization{}, err
	}
	return org, nil
}

func (r *Repository) DeleteOrganization(ctx context.Context, orgID string) error {
	tag, err := r.pool.Exec(ctx, `
		DELETE FROM organizations
		WHERE id = $1
	`, orgID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) GetMember(ctx context.Context, orgID, userID string) (Member, error) {
	var m Member
	err := r.pool.QueryRow(ctx, `
		SELECT u.id, u.name, u.email, om.role, om.created_at
		FROM organization_members om
		JOIN users u ON u.id = om.user_id
		WHERE om.organization_id = $1 AND om.user_id = $2
	`, orgID, userID).Scan(&m.User.ID, &m.User.Name, &m.User.Email, &m.Role, &m.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Member{}, ErrNotFound
		}
		return Member{}, err
	}
	return m, nil
}

func (r *Repository) GetMemberRole(ctx context.Context, orgID, userID string) (string, error) {
	var role string
	err := r.pool.QueryRow(ctx, `
		SELECT role
		FROM organization_members
		WHERE organization_id = $1 AND user_id = $2
	`, orgID, userID).Scan(&role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", err
	}
	return role, nil
}

func (r *Repository) ListMembers(ctx context.Context, orgID string) ([]Member, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT u.id, u.name, u.email, om.role, om.created_at
		FROM organization_members om
		JOIN users u ON u.id = om.user_id
		WHERE om.organization_id = $1
		ORDER BY om.created_at ASC, u.id ASC
	`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := make([]Member, 0)
	for rows.Next() {
		var m Member
		if err := rows.Scan(&m.User.ID, &m.User.Name, &m.User.Email, &m.Role, &m.CreatedAt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return members, nil
}

func (r *Repository) AddMember(ctx context.Context, orgID, userID, role string) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO organization_members (organization_id, user_id, role)
		VALUES ($1, $2, $3)
	`, orgID, userID, role)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrConflict
		}
		return err
	}
	return nil
}

func (r *Repository) UpdateMemberRole(ctx context.Context, orgID, userID, role string) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE organization_members
		SET role = $3
		WHERE organization_id = $1 AND user_id = $2
	`, orgID, userID, role)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) RemoveMember(ctx context.Context, orgID, userID string) error {
	tag, err := r.pool.Exec(ctx, `
		DELETE FROM organization_members
		WHERE organization_id = $1 AND user_id = $2
	`, orgID, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

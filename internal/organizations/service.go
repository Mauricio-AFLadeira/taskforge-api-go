package organizations

import (
	"context"
	"errors"
	"strings"

	"github.com/mauricio-reportei/taskforge-api-go/internal/users"
)

type userStore interface {
	GetByID(context.Context, string) (users.User, error)
	GetByEmail(context.Context, string) (users.User, error)
}

type organizationStore interface {
	CreateOrganization(context.Context, string, string) (Organization, error)
	ListOrganizationsForUser(context.Context, string) ([]Organization, error)
	GetOrganization(context.Context, string) (Organization, error)
	UpdateOrganization(context.Context, string, string) (Organization, error)
	DeleteOrganization(context.Context, string) error
	GetMember(context.Context, string, string) (Member, error)
	GetMemberRole(context.Context, string, string) (string, error)
	ListMembers(context.Context, string) ([]Member, error)
	AddMember(context.Context, string, string, string) error
	UpdateMemberRole(context.Context, string, string, string) error
	RemoveMember(context.Context, string, string) error
}

type Service struct {
	usersRepo userStore
	repo      organizationStore
}

func NewService(usersRepo userStore, repo organizationStore) *Service {
	return &Service{usersRepo: usersRepo, repo: repo}
}

func publicOrganization(org Organization) OrganizationPublic {
	return OrganizationPublic{
		ID:        org.ID,
		Name:      org.Name,
		OwnerID:   org.OwnerID,
		CreatedAt: org.CreatedAt,
		UpdatedAt: org.UpdatedAt,
	}
}

func publicMember(member Member) MemberPublic {
	return MemberPublic{
		User: UserPublic{
			ID:    member.User.ID,
			Name:  member.User.Name,
			Email: member.User.Email,
		},
		Role:      member.Role,
		CreatedAt: member.CreatedAt,
	}
}

func toPublicOrganizations(orgs []Organization) []OrganizationPublic {
	result := make([]OrganizationPublic, 0, len(orgs))
	for _, org := range orgs {
		result = append(result, publicOrganization(org))
	}
	return result
}

func toPublicMembers(members []Member) []MemberPublic {
	result := make([]MemberPublic, 0, len(members))
	for _, member := range members {
		result = append(result, publicMember(member))
	}
	return result
}

func (s *Service) loadOrganizationAccess(ctx context.Context, orgID, actorID string) (Organization, string, error) {
	org, err := s.repo.GetOrganization(ctx, orgID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return Organization{}, "", ErrNotFound
		}
		return Organization{}, "", err
	}
	role, err := s.repo.GetMemberRole(ctx, orgID, actorID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return Organization{}, "", ErrForbidden
		}
		return Organization{}, "", err
	}
	return org, role, nil
}

func validateOrganizationName(name string) string {
	return strings.TrimSpace(name)
}

func validateMemberRole(role string) string {
	return normalizeRole(strings.TrimSpace(strings.ToLower(role)))
}

func (s *Service) CreateOrganization(ctx context.Context, actorID, name string) (OrganizationPublic, error) {
	name = validateOrganizationName(name)
	if name == "" {
		return OrganizationPublic{}, ErrValidation
	}
	org, err := s.repo.CreateOrganization(ctx, actorID, name)
	if err != nil {
		return OrganizationPublic{}, err
	}
	return publicOrganization(org), nil
}

func (s *Service) ListOrganizations(ctx context.Context, actorID string) ([]OrganizationPublic, error) {
	orgs, err := s.repo.ListOrganizationsForUser(ctx, actorID)
	if err != nil {
		return nil, err
	}
	return toPublicOrganizations(orgs), nil
}

func (s *Service) GetOrganization(ctx context.Context, actorID, orgID string) (OrganizationPublic, error) {
	org, _, err := s.loadOrganizationAccess(ctx, orgID, actorID)
	if err != nil {
		return OrganizationPublic{}, err
	}
	return publicOrganization(org), nil
}

func (s *Service) UpdateOrganization(ctx context.Context, actorID, orgID, name string) (OrganizationPublic, error) {
	name = validateOrganizationName(name)
	if name == "" {
		return OrganizationPublic{}, ErrValidation
	}
	org, role, err := s.loadOrganizationAccess(ctx, orgID, actorID)
	if err != nil {
		return OrganizationPublic{}, err
	}
	if !canUpdateOrganization(role) {
		return OrganizationPublic{}, ErrForbidden
	}
	updated, err := s.repo.UpdateOrganization(ctx, org.ID, name)
	if err != nil {
		return OrganizationPublic{}, err
	}
	return publicOrganization(updated), nil
}

func (s *Service) DeleteOrganization(ctx context.Context, actorID, orgID string) error {
	org, role, err := s.loadOrganizationAccess(ctx, orgID, actorID)
	if err != nil {
		return err
	}
	if !canDeleteOrganization(role) {
		return ErrForbidden
	}
	return s.repo.DeleteOrganization(ctx, org.ID)
}

func (s *Service) AddMember(ctx context.Context, actorID, orgID, email, role string) (MemberPublic, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return MemberPublic{}, ErrValidation
	}
	role = validateMemberRole(role)
	if role == "" || role == RoleOwner {
		return MemberPublic{}, ErrValidation
	}
	_, actorRole, err := s.loadOrganizationAccess(ctx, orgID, actorID)
	if err != nil {
		return MemberPublic{}, err
	}
	if !canManageMembers(actorRole) {
		return MemberPublic{}, ErrForbidden
	}
	target, err := s.usersRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, users.ErrNotFound) {
			return MemberPublic{}, ErrNotFound
		}
		return MemberPublic{}, err
	}
	if err := s.repo.AddMember(ctx, orgID, target.ID, role); err != nil {
		if errors.Is(err, ErrConflict) {
			return MemberPublic{}, ErrConflict
		}
		return MemberPublic{}, err
	}
	member, err := s.repo.GetMember(ctx, orgID, target.ID)
	if err != nil {
		return MemberPublic{}, err
	}
	return publicMember(member), nil
}

func (s *Service) ListMembers(ctx context.Context, actorID, orgID string) ([]MemberPublic, error) {
	_, _, err := s.loadOrganizationAccess(ctx, orgID, actorID)
	if err != nil {
		return nil, err
	}
	members, err := s.repo.ListMembers(ctx, orgID)
	if err != nil {
		return nil, err
	}
	return toPublicMembers(members), nil
}

func (s *Service) ChangeMemberRole(ctx context.Context, actorID, orgID, targetUserID, role string) (MemberPublic, error) {
	role = validateMemberRole(role)
	if role == "" || role == RoleOwner {
		return MemberPublic{}, ErrValidation
	}
	_, actorRole, err := s.loadOrganizationAccess(ctx, orgID, actorID)
	if err != nil {
		return MemberPublic{}, err
	}
	if !canManageMembers(actorRole) {
		return MemberPublic{}, ErrForbidden
	}
	current, err := s.repo.GetMember(ctx, orgID, targetUserID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return MemberPublic{}, ErrNotFound
		}
		return MemberPublic{}, err
	}
	if current.Role == RoleOwner {
		return MemberPublic{}, ErrForbidden
	}
	if err := s.repo.UpdateMemberRole(ctx, orgID, targetUserID, role); err != nil {
		return MemberPublic{}, err
	}
	updated, err := s.repo.GetMember(ctx, orgID, targetUserID)
	if err != nil {
		return MemberPublic{}, err
	}
	return publicMember(updated), nil
}

func (s *Service) RemoveMember(ctx context.Context, actorID, orgID, targetUserID string) error {
	_, actorRole, err := s.loadOrganizationAccess(ctx, orgID, actorID)
	if err != nil {
		return err
	}
	if !canManageMembers(actorRole) {
		return ErrForbidden
	}
	target, err := s.repo.GetMember(ctx, orgID, targetUserID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		}
		return err
	}
	if target.Role == RoleOwner {
		return ErrForbidden
	}
	return s.repo.RemoveMember(ctx, orgID, targetUserID)
}

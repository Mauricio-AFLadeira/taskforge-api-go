package organizations

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mauricio-reportei/taskforge-api-go/internal/users"
)

type fakeUsersRepo struct {
	byEmail map[string]users.User
	byID    map[string]users.User
}

func (f fakeUsersRepo) GetByID(_ context.Context, id string) (users.User, error) {
	if u, ok := f.byID[id]; ok {
		return u, nil
	}
	return users.User{}, users.ErrNotFound
}

func (f fakeUsersRepo) GetByEmail(_ context.Context, email string) (users.User, error) {
	if u, ok := f.byEmail[email]; ok {
		return u, nil
	}
	return users.User{}, users.ErrNotFound
}

type fakeOrgRepo struct {
	orgs          map[string]Organization
	members       map[string]map[string]Member
	createdOwnerID string
	lastCreated    string
	lastAddedRole  string
	lastChanged    string
	deletedOrgID   string
	nextID         int
}

func newFakeOrgRepo() *fakeOrgRepo {
	return &fakeOrgRepo{
		orgs:    map[string]Organization{},
		members: map[string]map[string]Member{},
		nextID:  1,
	}
}

func (f *fakeOrgRepo) mkID(prefix string) string {
	id := prefix + "-" + string(rune('0'+f.nextID))
	f.nextID++
	return id
}

func (f *fakeOrgRepo) CreateOrganization(_ context.Context, ownerID, name string) (Organization, error) {
	org := Organization{ID: f.mkID("org"), Name: name, OwnerID: ownerID, CreatedAt: time.Unix(1, 0), UpdatedAt: time.Unix(1, 0)}
	f.orgs[org.ID] = org
	f.members[org.ID] = map[string]Member{
		ownerID: {User: UserPublic{ID: ownerID, Name: "Owner", Email: "owner@example.com"}, Role: RoleOwner, CreatedAt: time.Unix(1, 0)},
	}
	f.createdOwnerID = ownerID
	f.lastCreated = name
	return org, nil
}

func (f *fakeOrgRepo) ListOrganizationsForUser(_ context.Context, userID string) ([]Organization, error) {
	orgs := make([]Organization, 0)
	for orgID := range f.members {
		if _, ok := f.members[orgID][userID]; ok {
			orgs = append(orgs, f.orgs[orgID])
		}
	}
	return orgs, nil
}

func (f *fakeOrgRepo) GetOrganization(_ context.Context, orgID string) (Organization, error) {
	org, ok := f.orgs[orgID]
	if !ok {
		return Organization{}, ErrNotFound
	}
	return org, nil
}

func (f *fakeOrgRepo) UpdateOrganization(_ context.Context, orgID, name string) (Organization, error) {
	org, ok := f.orgs[orgID]
	if !ok {
		return Organization{}, ErrNotFound
	}
	org.Name = name
	org.UpdatedAt = time.Unix(2, 0)
	f.orgs[orgID] = org
	return org, nil
}

func (f *fakeOrgRepo) DeleteOrganization(_ context.Context, orgID string) error {
	if _, ok := f.orgs[orgID]; !ok {
		return ErrNotFound
	}
	delete(f.orgs, orgID)
	delete(f.members, orgID)
	f.deletedOrgID = orgID
	return nil
}

func (f *fakeOrgRepo) GetMember(_ context.Context, orgID, userID string) (Member, error) {
	group, ok := f.members[orgID]
	if !ok {
		return Member{}, ErrNotFound
	}
	m, ok := group[userID]
	if !ok {
		return Member{}, ErrNotFound
	}
	return m, nil
}

func (f *fakeOrgRepo) GetMemberRole(ctx context.Context, orgID, userID string) (string, error) {
	m, err := f.GetMember(ctx, orgID, userID)
	if err != nil {
		return "", err
	}
	return m.Role, nil
}

func (f *fakeOrgRepo) ListMembers(_ context.Context, orgID string) ([]Member, error) {
	group, ok := f.members[orgID]
	if !ok {
		return nil, ErrNotFound
	}
	members := make([]Member, 0, len(group))
	for _, member := range group {
		members = append(members, member)
	}
	return members, nil
}

func (f *fakeOrgRepo) AddMember(_ context.Context, orgID, userID, role string) error {
	group, ok := f.members[orgID]
	if !ok {
		return ErrNotFound
	}
	if _, exists := group[userID]; exists {
		return ErrConflict
	}
	group[userID] = Member{User: UserPublic{ID: userID, Name: "Member", Email: "member@example.com"}, Role: role, CreatedAt: time.Unix(3, 0)}
	f.lastAddedRole = role
	return nil
}

func (f *fakeOrgRepo) UpdateMemberRole(_ context.Context, orgID, userID, role string) error {
	group, ok := f.members[orgID]
	if !ok {
		return ErrNotFound
	}
	m, exists := group[userID]
	if !exists {
		return ErrNotFound
	}
	m.Role = role
	group[userID] = m
	f.lastChanged = role
	return nil
}

func (f *fakeOrgRepo) RemoveMember(_ context.Context, orgID, userID string) error {
	group, ok := f.members[orgID]
	if !ok {
		return ErrNotFound
	}
	if _, exists := group[userID]; !exists {
		return ErrNotFound
	}
	delete(group, userID)
	return nil
}

func (f *fakeOrgRepo) Exists(_ context.Context, orgID string) (bool, error) {
	_, ok := f.orgs[orgID]
	return ok, nil
}

func TestCreateOrganizationMakesCreatorOwner(t *testing.T) {
	t.Parallel()
	usersRepo := fakeUsersRepo{}
	repo := newFakeOrgRepo()
	svc := NewService(usersRepo, repo)

	out, err := svc.CreateOrganization(context.Background(), "user-1", "Acme")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if out.Name != "Acme" || out.OwnerID != "user-1" {
		t.Fatalf("unexpected org: %#v", out)
	}
	if repo.createdOwnerID != "user-1" {
		t.Fatalf("owner mismatch: %q", repo.createdOwnerID)
	}
	member, err := repo.GetMember(context.Background(), out.ID, "user-1")
	if err != nil {
		t.Fatalf("owner membership missing: %v", err)
	}
	if member.Role != RoleOwner {
		t.Fatalf("expected owner role, got %q", member.Role)
	}
}

func TestMemberCannotManageMembers(t *testing.T) {
	t.Parallel()
	usersRepo := fakeUsersRepo{byEmail: map[string]users.User{
		"member@example.com": {ID: "user-2", Name: "Member", Email: "member@example.com"},
	}}
	repo := newFakeOrgRepo()
	svc := NewService(usersRepo, repo)
	org, err := svc.CreateOrganization(context.Background(), "user-1", "Acme")
	if err != nil {
		t.Fatal(err)
	}
	group := repo.members[org.ID]
	group["user-2"] = Member{User: UserPublic{ID: "user-2", Name: "Member", Email: "member@example.com"}, Role: RoleMember, CreatedAt: time.Unix(2, 0)}

	_, err = svc.AddMember(context.Background(), "user-2", org.ID, "member@example.com", RoleMember)
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}
	if repo.lastAddedRole != "" {
		t.Fatalf("member should not have been added, got %q", repo.lastAddedRole)
	}
}

func TestOwnerCanAddMembers(t *testing.T) {
	t.Parallel()
	usersRepo := fakeUsersRepo{byEmail: map[string]users.User{
		"member@example.com": {ID: "user-2", Name: "Member", Email: "member@example.com"},
	}}
	repo := newFakeOrgRepo()
	svc := NewService(usersRepo, repo)
	org, err := svc.CreateOrganization(context.Background(), "user-1", "Acme")
	if err != nil {
		t.Fatal(err)
	}

	out, err := svc.AddMember(context.Background(), "user-1", org.ID, "member@example.com", RoleMember)
	if err != nil {
		t.Fatalf("add member: %v", err)
	}
	if out.User.ID != "user-2" || out.Role != RoleMember {
		t.Fatalf("unexpected member: %#v", out)
	}
	if repo.lastAddedRole != RoleMember {
		t.Fatalf("role mismatch: %q", repo.lastAddedRole)
	}
}

func TestOwnerCannotBeRemoved(t *testing.T) {
	t.Parallel()
	usersRepo := fakeUsersRepo{}
	repo := newFakeOrgRepo()
	svc := NewService(usersRepo, repo)
	org, err := svc.CreateOrganization(context.Background(), "user-1", "Acme")
	if err != nil {
		t.Fatal(err)
	}

	err = svc.RemoveMember(context.Background(), "user-1", org.ID, "user-1")
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}
}

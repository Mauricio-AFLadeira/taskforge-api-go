package organizations

import "testing"

func TestPermissionHelpers(t *testing.T) {
	t.Parallel()
	if !canManageMembers(RoleOwner) || !canManageMembers(RoleAdmin) || canManageMembers(RoleMember) {
		t.Fatal("manage members helper mismatch")
	}
	if !canUpdateOrganization(RoleOwner) || canUpdateOrganization(RoleAdmin) || canUpdateOrganization(RoleMember) {
		t.Fatal("update helper mismatch")
	}
	if !canDeleteOrganization(RoleOwner) || canDeleteOrganization(RoleAdmin) || canDeleteOrganization(RoleMember) {
		t.Fatal("delete helper mismatch")
	}
	if normalizeRole("owner") != RoleOwner || normalizeRole("bad") != "" {
		t.Fatal("normalize role mismatch")
	}
}

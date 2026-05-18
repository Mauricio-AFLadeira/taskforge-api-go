package organizations

func canManageMembers(role string) bool {
	return role == RoleOwner || role == RoleAdmin
}

func canUpdateOrganization(role string) bool {
	return role == RoleOwner
}

func canDeleteOrganization(role string) bool {
	return role == RoleOwner
}

func normalizeRole(role string) string {
	switch role {
	case RoleOwner, RoleAdmin, RoleMember:
		return role
	default:
		return ""
	}
}

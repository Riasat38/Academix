package permissions

var permissions map[string]map[string][]string

func init() {
	/*
		role:{
			resource: {permissionList}
		}
	*/
	permissions = map[string]map[string][]string{
		"student": {
			"course":  {"viewAll", "viewOwn", "enroll", "view"},
			"profile": {"view", "edit"},
		},
		"teacher": {
			"course":  {"viewAll", "viewOwn", "update", "view"},
			"profile": {"view", "edit"},
		},
		"admin": {
			"course":  {"viewAll", "viewOwn", "modify", "create", "delete", "addUser", "view"},
			"profile": {"view", "edit", "delete"},
			"user":    {"view", "edit", "delete"},
		},
	}
}

func ValidatePermission(role string, resource string, expectedPermission string) bool {
	for _, action := range permissions[role][resource] {
		if action == expectedPermission {
			return true
		}
	}
	return false

}

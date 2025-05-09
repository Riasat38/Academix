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
			"course":     {"viewAll", "viewOwn", "enroll", "view"},
			"profile":    {"view", "edit"},
			"assignment": {"view", "viewAll"},
			"submission": {"view", "post", "getMarks:Feedback"},
		},
		"teacher": {
			"course":     {"viewAll", "viewOwn", "update", "view"},
			"profile":    {"view", "edit"},
			"assignment": {"view", "create", "edit", "post", "delete"},
			"submission": {"viewAll", "post", "viewMarks", "postMarks:Feedback", "delete"},
		},
		"admin": {
			"course":     {"viewAll", "viewOwn", "modify", "create", "delete", "addUser", "view"},
			"profile":    {"view", "edit", "delete"},
			"user":       {"view", "edit", "delete"},
			"assignment": {"view", "viewAll", "create", "edit", "post", "delete"},
			"submission": {"view", "viewAll", "post", "getMarks:Feedback"},
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

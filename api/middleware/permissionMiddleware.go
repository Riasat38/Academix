package middleware

const (
	RolePermission = map[string]string{
		"student": {"enroll_course"},
		"teacher": {"update_course"},
	}
)

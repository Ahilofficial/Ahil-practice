package constants

const (
	AdminRole     = "admin"
	PrincipalRole = "principal"
	FacultyRole   = "faculty"
	StudentRole   = "student"
	UserRole      = "user"
)

var (
	rolePermissions = map[string][]string{
		AdminRole: {
			"CREATE_INSTITUTION",
			"VIEW_INSTITUTIONS",
			"UPDATE_INSTITUTION",
			"DELETE_INSTITUTION",
		},

		PrincipalRole: {
			"VIEW_INSTITUTIONS",
			"UPDATE_INSTITUTION",
		},

		FacultyRole: {
			"VIEW_FACULTIES",
			"UPDATE_FACULTIES",
		},

		StudentRole: {
			"VIEW_STUDENTS",
			"UPDATE_STUDENTS",
		},
		UserRole: {
			"VIEW_INSTITUTIONS",
			"VIEW_DEPARTMENTS",
			"VIEW_FACULTIES",
			"VIEW_STUDENTS",
		},
	}
)

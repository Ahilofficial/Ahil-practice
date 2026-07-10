package constants

// import "google.golang.org/grpc/admin"

const (
	AdminRole     = "admin"
	PrincipalRole = "principal"
	FacultyRole   = "faculty"
	StudentRole   = "student"
	UserRole      = "user"
)

const (
	PermissionCreateInstitutes = "CREATE_INSTITUTION"
	PermissionViewInstitutes   = "VIEW_INSTITUTIONS"
	PermissionUpdateInstitutes = "UPDATE_INSTITUTION"
	PermissionDeleteInstitutes = "DELETE_INSTITUTION"

	PermissionCreateDepartments = "CREATE_DEPARTMENT"
	PermissionViewDepartments   = "VIEW_DEPARTMENTS"
	PermissionUpdateDepartments = "UPDATE_DEPARTMENT"
	PermissionDeleteDepartments = "DELETE_DEPARTMENT"

	PermissionCreateFaculties = "CREATE_FACULTY"
	PermissionViewFaculties   = "VIEW_FACULTIES"
	PermissionUpdateFaculties = "UPDATE_FACULTY"
	PermissionDeleteFaculties = "DELETE_FACULTY"

	PermissionCreateStudents = "CREATE_STUDENT"
	PermissionViewStudents   = "VIEW_STUDENTS"
	PermissionUpdateStudents = "UPDATE_STUDENT"
	PermissionDeleteStudents = "DELETE_STUDENT"

	PermissionCreateFees = "CREATE_FEE"
	PermissionViewFees   = "VIEW_FEES"
	PermissionUpdateFees = "UPDATE_FEE"
	PermissionDeleteFees = "DELETE_FEE"

	PermissionAssignRoles = "ASSIGN_ROLE"
)

var PermissionGroups = map[string][]string{

	"institute": {
		PermissionCreateInstitutes,
		PermissionViewInstitutes,
		PermissionUpdateInstitutes,
		PermissionDeleteInstitutes,
	},

	"department": {
		PermissionCreateDepartments,
		PermissionViewDepartments,
		PermissionUpdateDepartments,
		PermissionDeleteDepartments,
	},
	"faculty": {
		PermissionCreateFaculties,
		PermissionViewFaculties,
		PermissionUpdateFaculties,
		PermissionDeleteFaculties,
	},
	"student": {
		PermissionCreateStudents,
		PermissionViewStudents,
		PermissionUpdateStudents,
		PermissionDeleteStudents,
	},
	
	"other": {
		PermissionAssignRoles,
	},
	
}

var Admin []string
func init() {
	for _, roles := range PermissionGroups {
		Admin = append(Admin, roles...)
	}
}
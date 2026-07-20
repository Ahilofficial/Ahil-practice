package constants

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

// var PermissionGroups = map[string][]string{
// 	"admin": {
// 		PermissionCreateInstitutes,
// 		PermissionViewInstitutes,
// 		PermissionUpdateInstitutes,
// 		PermissionDeleteInstitutes,

// 		PermissionCreateDepartments,
// 		PermissionViewDepartments,
// 		PermissionUpdateDepartments,
// 		PermissionDeleteDepartments,

// 		PermissionCreateFaculties,
// 		PermissionViewFaculties,
// 		PermissionUpdateFaculties,
// 		PermissionDeleteFaculties,

// 		PermissionCreateStudents,
// 		PermissionViewStudents,
// 		PermissionUpdateStudents,
// 		PermissionDeleteStudents,

// 		PermissionCreateFees,
// 		PermissionViewFees,
// 		PermissionUpdateFees,
// 		PermissionDeleteFees,

// 		PermissionAssignRoles,
// 	},

// 	"principal": {
// 		PermissionCreateDepartments,
// 		PermissionViewDepartments,
// 		PermissionUpdateDepartments,
// 		PermissionDeleteDepartments,

// 		PermissionCreateFaculties,
// 		PermissionViewFaculties,
// 		PermissionUpdateFaculties,
// 		PermissionDeleteFaculties,

// 		PermissionCreateStudents,
// 		PermissionViewStudents,
// 		PermissionUpdateStudents,
// 		PermissionDeleteStudents,
// 	},

// 	"faculty": {
// 		PermissionCreateFaculties,
// 		PermissionViewFaculties,
// 		PermissionUpdateFaculties,
// 		PermissionDeleteFaculties,
// 		PermissionCreateStudents,
// 		PermissionViewStudents,
// 		PermissionUpdateStudents,
// 		PermissionDeleteStudents,
// 	},

// 	"student": {
// 		PermissionCreateStudents,
// 		PermissionViewStudents,
// 		PermissionUpdateStudents,
// 		PermissionDeleteStudents,
// 	},

// 	"user": {},
// }

var (
	instPerms        = []string{PermissionCreateInstitutes, PermissionViewInstitutes, PermissionUpdateInstitutes, PermissionDeleteInstitutes}
	deptPerms        = []string{PermissionCreateDepartments, PermissionViewDepartments, PermissionUpdateDepartments, PermissionDeleteDepartments}
	facPerms         = []string{PermissionCreateFaculties, PermissionViewFaculties, PermissionUpdateFaculties, PermissionDeleteFaculties}
	studPerms        = []string{PermissionCreateStudents, PermissionViewStudents, PermissionUpdateStudents, PermissionDeleteStudents}
	feePerms         = []string{PermissionCreateFees, PermissionViewFees, PermissionUpdateFees, PermissionDeleteFees}
	PermissionGroups = map[string][]string{}
)

func init() {
	PermissionGroups["student"] = studPerms
	PermissionGroups["faculty"] = append(append([]string{}, facPerms...), studPerms...)
	PermissionGroups["principal"] = append(append([]string{}, deptPerms...), PermissionGroups["faculty"]...)
	PermissionGroups["admin"] = append(append(append([]string{}, instPerms...), feePerms...), append(PermissionGroups["principal"], PermissionAssignRoles)...)
	PermissionGroups["user"] = []string{}
}

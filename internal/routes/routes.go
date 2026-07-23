package routes

import (
	"backend_institutions/internal/constants"
	"backend_institutions/internal/controller"
	"backend_institutions/internal/middleware"
	"fmt"

	"github.com/gofiber/fiber/v3"
)

func NewApp(
	userController *controller.UserController,
	instituteController *controller.InstituteController,
	departmentController *controller.DepartmentController,
	facultyController *controller.FacultyController,
	studentController *controller.StudentController,
	feesController *controller.FeesController,
	roleController *controller.RoleController,
	permissionController *controller.PermissionController,
) *fiber.App {
	app := fiber.New()
	RegisterRoutes(
		app,
		userController,
		instituteController,
		departmentController,
		facultyController,
		studentController,
		feesController,
		roleController,
		permissionController,
	)
	return app
}

func RegisterRoutes(
	app *fiber.App,
	userController *controller.UserController,
	instituteController *controller.InstituteController,
	departmentController *controller.DepartmentController,
	facultyController *controller.FacultyController,
	studentController *controller.StudentController,
	feesController *controller.FeesController,
	roleController *controller.RoleController,
	permissionController *controller.PermissionController,
) {
	app.Get("/auth/verify", userController.VerifyEmail)
	app.Post("/auth/forgot-password", userController.ForgotPassword)
	app.Post("/auth/reset-password", userController.ResetPassword)
	
	
	app.Use(middleware.RequestResponseLogger())

	app.Post("/signup", userController.SignUpController)
	app.Post("/signin", userController.SignInController)
	app.Post("/logout", userController.Logout)

	protected := app.Group("", middleware.AuthRequired())

	protected.Post("/users/assign-role", middleware.RequirePermission(constants.PermissionAssignRoles), userController.AssignRoleController)

	roleRoute := protected.Group("/roles", middleware.RequirePermission(constants.PermissionAssignRoles))
	roleRoute.Post("", roleController.CreateRoleController)


	
	roleRoute.Post("/:id/permissions", roleController.AssignPermissionsController)
	roleRoute.Get("/:id/permissions", roleController.GetRolePermissionsController)
	roleRoute.Delete("/:id/permissions/:permissionId", roleController.RemovePermissionController)

	// Permission Read API
	protected.Get("/permissions", middleware.RequirePermission(constants.PermissionAssignRoles), permissionController.GetAllPermissionsController)

	InstituteRoute := protected.Group("/institutes")
	InstituteRoute.Post("", middleware.RequirePermission(constants.PermissionCreateInstitutes), instituteController.CreateInstituteController)
	InstituteRoute.Get("", middleware.RequirePermission(constants.PermissionViewInstitutes), instituteController.GetAllInstitutesController)
	InstituteRoute.Get("/all", middleware.RequirePermission(constants.PermissionViewInstitutes), instituteController.FetchAllInstitutesController)
	InstituteRoute.Get("/active", middleware.RequirePermission(constants.PermissionViewInstitutes), instituteController.GetActiveInstituteController)
	InstituteRoute.Get("/inactive", middleware.RequirePermission(constants.PermissionViewInstitutes), instituteController.GetInactiveInstituteController)
	InstituteRoute.Get("/deleted", middleware.RequirePermission(constants.PermissionViewInstitutes), instituteController.GetDeletedInstitutesController)
	InstituteRoute.Get("/:id", middleware.RequirePermission(constants.PermissionViewInstitutes), instituteController.GetInstituteByIDController)
	InstituteRoute.Put("/:id", middleware.RequirePermission(constants.PermissionUpdateInstitutes), instituteController.UpdateInstituteController)
	InstituteRoute.Delete("/:id", middleware.RequirePermission(constants.PermissionDeleteInstitutes), instituteController.DeleteInstituteController)

	DepartmentRoute := protected.Group("/departments")
	DepartmentRoute.Post("", middleware.RequirePermission(constants.PermissionCreateDepartments), departmentController.CreateDepartmentController)
	DepartmentRoute.Get("", middleware.RequirePermission(constants.PermissionViewDepartments), departmentController.GetAllDepartmentsController)
	DepartmentRoute.Get("/all", middleware.RequirePermission(constants.PermissionViewDepartments), departmentController.FetchAllDepartmentsController)
	DepartmentRoute.Get("/active", middleware.RequirePermission(constants.PermissionViewDepartments), departmentController.GetActiveDepartmentController)
	DepartmentRoute.Get("/inactive", middleware.RequirePermission(constants.PermissionViewDepartments), departmentController.GetInactiveDepartmentController)
	DepartmentRoute.Get("/deleted", middleware.RequirePermission(constants.PermissionViewDepartments), departmentController.GetDeletedDepartmentsController)
	DepartmentRoute.Get("/:id", middleware.RequirePermission(constants.PermissionViewDepartments), departmentController.GetDepartmentByIDController)
	DepartmentRoute.Put("/:id", middleware.RequirePermission(constants.PermissionUpdateDepartments), departmentController.UpdateDepartmentController)
	DepartmentRoute.Delete("/:id", middleware.RequirePermission(constants.PermissionDeleteDepartments), departmentController.DeleteDepartmentController)

	FacultyRoute := protected.Group("/faculties")
	FacultyRoute.Post("", middleware.RequirePermission(constants.PermissionCreateFaculties), facultyController.CreateFacultyController)
	FacultyRoute.Get("", middleware.RequirePermission(constants.PermissionViewFaculties), facultyController.GetAllFacultiesController)
	FacultyRoute.Get("/all", middleware.RequirePermission(constants.PermissionViewFaculties), facultyController.FetchAllFacultiesController)
	FacultyRoute.Get("/active", middleware.RequirePermission(constants.PermissionViewFaculties), facultyController.GetActiveFacultyController)
	FacultyRoute.Get("/inactive", middleware.RequirePermission(constants.PermissionViewFaculties), facultyController.GetInactiveFacultyController)
	FacultyRoute.Get("/deleted", middleware.RequirePermission(constants.PermissionViewFaculties), facultyController.GetDeletedFacultiesController)
	FacultyRoute.Get("/:id", middleware.RequirePermission(constants.PermissionViewFaculties), facultyController.GetFacultyByIDController)
	FacultyRoute.Put("/:id", middleware.RequirePermission(constants.PermissionUpdateFaculties), facultyController.UpdateFacultyController)
	FacultyRoute.Delete("/:id", middleware.RequirePermission(constants.PermissionDeleteFaculties), facultyController.DeleteFacultyController)

	StudentRoute := protected.Group("/students")
	StudentRoute.Post("", middleware.RequirePermission(constants.PermissionCreateStudents), studentController.CreateStudentControllers)
	StudentRoute.Get("/all", middleware.RequirePermission(constants.PermissionViewStudents), studentController.FetchAllStudentsControllers)
	StudentRoute.Get("/active", middleware.RequirePermission(constants.PermissionViewStudents), studentController.GetActiveStudentController)
	StudentRoute.Get("/inactive", middleware.RequirePermission(constants.PermissionViewStudents), studentController.GetInactiveStudentController)
	StudentRoute.Get("/:id", middleware.RequirePermission(constants.PermissionViewStudents), studentController.GetStudentByIDControllers)
	StudentRoute.Put("/:id", middleware.RequirePermission(constants.PermissionUpdateStudents), studentController.UpdateStudentControllers)
	StudentRoute.Delete("/:id", middleware.RequirePermission(constants.PermissionDeleteStudents), studentController.DeleteStudentControllers)

	FeesRoute := protected.Group("/fees")
	FeesRoute.Post("", middleware.RequirePermission(constants.PermissionCreateFees), feesController.CreateFeesController)
	FeesRoute.Get("", middleware.RequirePermission(constants.PermissionViewFees), feesController.GetAllFeesController)
	FeesRoute.Get("/all", middleware.RequirePermission(constants.PermissionViewFees), feesController.FetchAllFeesController)
	FeesRoute.Get("/inactive", middleware.RequirePermission(constants.PermissionViewFees), feesController.GetInactiveFeesController)
	FeesRoute.Get("/:id", middleware.RequirePermission(constants.PermissionViewFees), feesController.GetFeesByIDController)
	FeesRoute.Put("/:id", middleware.RequirePermission(constants.PermissionUpdateFees), feesController.UpdateFeesController)
	FeesRoute.Delete("/:id", middleware.RequirePermission(constants.PermissionDeleteFees), feesController.DeleteFeesController)

	userRoute := protected.Group("/users")
	userRoute.Post("/assign-role", userController.AssignRoleController)
	userRoute.Delete("/:id", userController.DeleteUserController)
	

	fmt.Println("All routes registered successfully")
}

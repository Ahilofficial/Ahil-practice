package routes

import (
	"backend_institutions/controller"
	"backend_institutions/middleware"
	"fmt"

	"github.com/gofiber/fiber/v3"
)

func SetUpRoutes(app *fiber.App) {
	fmt.Println("Routes Loaded")

	//public
	app.Post("/signup", controller.SignUpController)
	app.Post("/signin", controller.SignInController)

	//private

	protected := app.Group("", middleware.AuthRequired())

	protected.Post("/users/assign-role", middleware.RequirePermission("assign_roles"), controller.AssignRoleController)

	instituteRoute := protected.Group("/institutes")
	instituteRoute.Post("", middleware.RequirePermission("create_institutes"), controller.CreateInstituteController)
	instituteRoute.Get("", middleware.RequirePermission("view_institutes"), controller.GetAllInstitutesController)
	instituteRoute.Get("/all", middleware.RequirePermission("view_institutes"), controller.FetchAllInstitutesController)
	instituteRoute.Get("/active", middleware.RequirePermission("view_institutes"), controller.GetActiveInstituteController)
	instituteRoute.Get("/inactive", middleware.RequirePermission("view_institutes"), controller.GetInactiveInstituteController)
	instituteRoute.Get("/deleted", middleware.RequirePermission("view_institutes"), controller.GetDeletedInstitutesController)
	instituteRoute.Get("/:id", middleware.RequirePermission("view_institutes"), controller.GetInstituteByIDController)
	instituteRoute.Put("/:id", middleware.RequirePermission("update_institutes"), controller.UpdateInstituteController)
	instituteRoute.Delete("/:id", middleware.RequirePermission("delete_institutes"), controller.DeleteInstituteController)

	departmentRoute := protected.Group("/departments")
	departmentRoute.Post("", middleware.RequirePermission("create_departments"), controller.CreateDepartmentController)
	departmentRoute.Get("", middleware.RequirePermission("view_departments"), controller.GetAllDepartmentsController)
	departmentRoute.Get("/all", middleware.RequirePermission("view_departments"), controller.FetchAllDepartmentsController)
	departmentRoute.Get("/active", middleware.RequirePermission("view_departments"), controller.GetActiveDepartmentController)
	departmentRoute.Get("/inactive", middleware.RequirePermission("view_departments"), controller.GetInactiveDepartmentController)
	departmentRoute.Get("/deleted", middleware.RequirePermission("view_departments"), controller.GetDeletedDepartmentsController)
	departmentRoute.Get("/:id", middleware.RequirePermission("view_departments"), controller.GetDepartmentByIDController)
	departmentRoute.Put("/:id", middleware.RequirePermission("update_departments"), controller.UpdateDepartmentController)
	departmentRoute.Delete("/:id", middleware.RequirePermission("delete_departments"), controller.DeleteDepartmentController)

	facultyRoute := protected.Group("/faculties")
	facultyRoute.Post("", middleware.RequirePermission("create_faculties"), controller.CreateFacultyController)
	facultyRoute.Get("", middleware.RequirePermission("view_faculties"), controller.GetAllFacultiesController)
	facultyRoute.Get("/all", middleware.RequirePermission("view_faculties"), controller.FetchAllFacultiesController)
	facultyRoute.Get("/active", middleware.RequirePermission("view_faculties"), controller.GetActiveFacultyController)
	facultyRoute.Get("/inactive", middleware.RequirePermission("view_faculties"), controller.GetInactiveFacultyController)
	facultyRoute.Get("/deleted", middleware.RequirePermission("view_faculties"), controller.GetDeletedFacultiesController)
	facultyRoute.Get("/:id", middleware.RequirePermission("view_faculties"), controller.GetFacultyByIDController)
	facultyRoute.Put("/:id", middleware.RequirePermission("update_faculties"), controller.UpdateFacultyController)
	facultyRoute.Delete("/:id", middleware.RequirePermission("delete_faculties"), controller.DeleteFacultyController)

	studentRoute := protected.Group("/students")
	studentRoute.Post("", middleware.RequirePermission("create_students"), controller.CreateStudentControllers)
	studentRoute.Get("", middleware.RequirePermission("view_students"), controller.GetAllStudentsControllers)
	studentRoute.Get("/all", middleware.RequirePermission("view_students"), controller.FetchAllStudentsControllers)
	studentRoute.Get("/active", middleware.RequirePermission("view_students"), controller.GetActiveStudentController)
	studentRoute.Get("/inactive", middleware.RequirePermission("view_students"), controller.GetInactiveStudentController)
	studentRoute.Get("/deleted", middleware.RequirePermission("view_students"), controller.GetDeletedStudentsController)
	studentRoute.Get("/:id", middleware.RequirePermission("view_students"), controller.GetStudentByIDControllers)
	studentRoute.Put("/:id", middleware.RequirePermission("update_students"), controller.UpdateStudentControllers)
	studentRoute.Delete("/:id", middleware.RequirePermission("delete_students"), controller.DeleteStudentControllers)

	feesRoute := protected.Group("/fees")
	feesRoute.Post("", middleware.RequirePermission("create_fees"), controller.CreateFeesController)
	feesRoute.Get("", middleware.RequirePermission("view_fees"), controller.GetAllFeesController)
	feesRoute.Get("/all", middleware.RequirePermission("view_fees"), controller.FetchAllFeesController)
	feesRoute.Get("/inactive", middleware.RequirePermission("view_fees"), controller.GetInactiveFeesController)
	feesRoute.Get("/:id", middleware.RequirePermission("view_fees"), controller.GetFeesByIDController)
	feesRoute.Put("/:id", middleware.RequirePermission("update_fees"), controller.UpdateFeesController)
	feesRoute.Delete("/:id", middleware.RequirePermission("delete_fees"), controller.DeleteFeesController)

	userRoute := protected.Group("/users")
	userRoute.Post("/assign-role", controller.AssignRoleController)

	fmt.Println("All routes registered successfully")
}

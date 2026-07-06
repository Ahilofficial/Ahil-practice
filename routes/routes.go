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

	InstituteRoute := protected.Group("/institutes")
	InstituteRoute.Post("", middleware.RequirePermission("create_institutes"), controller.CreateInstituteController)
	InstituteRoute.Get("", middleware.RequirePermission("view_institutes"), controller.GetAllInstitutesController)
	InstituteRoute.Get("/all", middleware.RequirePermission("view_institutes"), controller.FetchAllInstitutesController)
	InstituteRoute.Get("/active", middleware.RequirePermission("view_institutes"), controller.GetActiveInstituteController)
	InstituteRoute.Get("/inactive", middleware.RequirePermission("view_institutes"), controller.GetInactiveInstituteController)
	InstituteRoute.Get("/deleted", middleware.RequirePermission("view_institutes"), controller.GetDeletedInstitutesController)
	InstituteRoute.Get("/:id", middleware.RequirePermission("view_institutes"), controller.GetInstituteByIDController)
	InstituteRoute.Put("/:id", middleware.RequirePermission("update_institutes"), controller.UpdateInstituteController)
	InstituteRoute.Delete("/:id", middleware.RequirePermission("delete_institutes"), controller.DeleteInstituteController)

	DepartmentRoute := protected.Group("/departments")
	DepartmentRoute.Post("", middleware.RequirePermission("create_departments"), controller.CreateDepartmentController)
	DepartmentRoute.Get("", middleware.RequirePermission("view_departments"), controller.GetAllDepartmentsController)
	DepartmentRoute.Get("/all", middleware.RequirePermission("view_departments"), controller.FetchAllDepartmentsController)
	DepartmentRoute.Get("/active", middleware.RequirePermission("view_departments"), controller.GetActiveDepartmentController)
	DepartmentRoute.Get("/inactive", middleware.RequirePermission("view_departments"), controller.GetInactiveDepartmentController)
	DepartmentRoute.Get("/deleted", middleware.RequirePermission("view_departments"), controller.GetDeletedDepartmentsController)
	DepartmentRoute.Get("/:id", middleware.RequirePermission("view_departments"), controller.GetDepartmentByIDController)
	DepartmentRoute.Put("/:id", middleware.RequirePermission("update_departments"), controller.UpdateDepartmentController)
	DepartmentRoute.Delete("/:id", middleware.RequirePermission("delete_departments"), controller.DeleteDepartmentController)

	FacultyRoute := protected.Group("/faculties")
	FacultyRoute.Post("", middleware.RequirePermission("create_faculties"), controller.CreateFacultyController)
	FacultyRoute.Get("", middleware.RequirePermission("view_faculties"), controller.GetAllFacultiesController)
	FacultyRoute.Get("/all", middleware.RequirePermission("view_faculties"), controller.FetchAllFacultiesController)
	FacultyRoute.Get("/active", middleware.RequirePermission("view_faculties"), controller.GetActiveFacultyController)
	FacultyRoute.Get("/inactive", middleware.RequirePermission("view_faculties"), controller.GetInactiveFacultyController)
	FacultyRoute.Get("/deleted", middleware.RequirePermission("view_faculties"), controller.GetDeletedFacultiesController)
	FacultyRoute.Get("/:id", middleware.RequirePermission("view_faculties"), controller.GetFacultyByIDController)
	FacultyRoute.Put("/:id", middleware.RequirePermission("update_faculties"), controller.UpdateFacultyController)
	FacultyRoute.Delete("/:id", middleware.RequirePermission("delete_faculties"), controller.DeleteFacultyController)

	StudentRoute := protected.Group("/students")
	StudentRoute.Post("", middleware.RequirePermission("create_students"), controller.CreateStudentControllers)
	StudentRoute.Get("", middleware.RequirePermission("view_students"), controller.GetAllStudentsControllers)
	StudentRoute.Get("/all", middleware.RequirePermission("view_students"), controller.FetchAllStudentsControllers)
	StudentRoute.Get("/active", middleware.RequirePermission("view_students"), controller.GetActiveStudentController)
	StudentRoute.Get("/inactive", middleware.RequirePermission("view_students"), controller.GetInactiveStudentController)
	StudentRoute.Get("/deleted", middleware.RequirePermission("view_students"), controller.GetDeletedStudentsController)
	StudentRoute.Get("/:id", middleware.RequirePermission("view_students"), controller.GetStudentByIDControllers)
	StudentRoute.Put("/:id", middleware.RequirePermission("update_students"), controller.UpdateStudentControllers)
	StudentRoute.Delete("/:id", middleware.RequirePermission("delete_students"), controller.DeleteStudentControllers)

	FeesRoute := protected.Group("/fees")
	FeesRoute.Post("", middleware.RequirePermission("create_fees"), controller.CreateFeesController)
	FeesRoute.Get("", middleware.RequirePermission("view_fees"), controller.GetAllFeesController)
	FeesRoute.Get("/all", middleware.RequirePermission("view_fees"), controller.FetchAllFeesController)
	FeesRoute.Get("/inactive", middleware.RequirePermission("view_fees"), controller.GetInactiveFeesController)
	FeesRoute.Get("/:id", middleware.RequirePermission("view_fees"), controller.GetFeesByIDController)
	FeesRoute.Put("/:id", middleware.RequirePermission("update_fees"), controller.UpdateFeesController)
	FeesRoute.Delete("/:id", middleware.RequirePermission("delete_fees"), controller.DeleteFeesController)

	userRoute := protected.Group("/users")
	userRoute.Post("/assign-role", controller.AssignRoleController)

	fmt.Println("All routes registered successfully")
}

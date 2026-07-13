package controller

import (
	"backend_institutions/internal/dto"
	"backend_institutions/internal/helper"
	"backend_institutions/internal/model"
	"backend_institutions/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type StudentController struct {
	studentService *services.StudentService
}

func NewStudentController(studentService *services.StudentService) *StudentController {
	return &StudentController{studentService: studentService}
}

func (cl *StudentController) CreateStudentControllers(c fiber.Ctx) error {
	var student model.Student
	if err := c.Bind().Body(&student); err != nil {
		return helper.Error(c, 400, "invalid request body: "+err.Error())
	}

	if student.Name == "" {
		return helper.Error(c, 400, "name is required")
	}
	if student.Email == "" {
		return helper.Error(c, 400, "email is required")
	}
	if student.FacultyID == 0 {
		return helper.Error(c, 400, "faculty_id is required")
	}

	createdStudent, err := cl.studentService.CreateStudentService(&student)
	if err != nil {
		return helper.Error(c, 400, err.Error())
	}

	return helper.Success(
		c,
		"Student created successfully",
		dto.ToStudentResponseDTO(&createdStudent),
	)
}

func (cl *StudentController) GetActiveStudentController(c fiber.Ctx) error {
	student, err := cl.studentService.GetActiveStudentService()
	if err != nil {
		return helper.Error(c, 404, err.Error())
	}

	return helper.Success(
		c,
		"Active student fetched successfully",
		dto.ToStudentResponseDTO(&student),
	)
}

func (cl *StudentController) GetInactiveStudentController(c fiber.Ctx) error {
	student, err := cl.studentService.GetInactiveStudentService()
	if err != nil {
		return helper.Error(c, 404, err.Error())
	}

	return helper.Success(
		c,
		"Inactive student fetched successfully",
		dto.ToStudentResponseDTO(&student),
	)
}

func (cl *StudentController) GetStudentByIDControllers(c fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return helper.Error(c, 400, "invalid student id")
	}

	student, err := cl.studentService.GetStudentServiceById(uint(id))
	if err != nil {
		return helper.Error(c, 404, err.Error())
	}

	return helper.Success(
		c,
		"Student fetched successfully",
		dto.ToStudentResponseDTO(&student),
	)
}

func (cl *StudentController) UpdateStudentControllers(c fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return helper.Error(c, 400, "invalid student id")
	}

	var body dto.UpdateStudentDTO
	body.Sanitize()

	if err := c.Bind().Body(&body); err != nil {
		return helper.Error(c, 400, "invalid request body")
	}

	if err := body.Validate(); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	if err := cl.studentService.UpdateStudentService(
		uint(id),
		&body,
	); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	updated, err := cl.studentService.GetStudentServiceById(uint(id))
	if err != nil {
		return helper.Error(c, 500, err.Error())
	}

	return helper.Success(
		c,
		"Student updated successfully",
		dto.ToStudentResponseDTO(&updated),
	)
}

func (cl *StudentController) DeleteStudentControllers(c fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return helper.Error(c, 400, "invalid student id")
	}

	if err := cl.studentService.DeleteStudentService(uint(id)); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	return helper.Success(
		c,
		"Student deleted successfully",
		nil,
	)
}

func (cl *StudentController) FetchAllStudentsControllers(c fiber.Ctx) error {
	students, err := cl.studentService.GetStudentService()
	if err != nil {
		return helper.Error(c, 500, err.Error())
	}

	return helper.Success(
		c,
		"All students fetched successfully",
		dto.ToStudentResponseListDTO(students),
	)
}
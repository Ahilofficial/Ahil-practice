package controller

import (
"backend_institutions/helper"
"backend_institutions/model"
"backend_institutions/services"
"strconv"

"github.com/gofiber/fiber/v3"
)

type FacultyController struct {
facultyService *services.FacultyService
}

func NewFacultyController() *FacultyController {
return &FacultyController{
facultyService: services.NewFacultyService(),
}
}

func (cl *FacultyController) CreateFacultyController(c fiber.Ctx) error {
var faculty model.Faculty

if err := c.Bind().Body(&faculty); err != nil {
return helper.Error(c, 400, "invalid request body")
}

if err := cl.facultyService.CreateFacultyService(&faculty); err != nil {
return helper.Error(c, 400, err.Error())
}

return helper.Success(c, "Faculty created successfully", faculty)
}

func (cl *FacultyController) GetAllFacultiesController(c fiber.Ctx) error {
faculties, err := cl.facultyService.GetFacultyService()
if err != nil {
return helper.Error(c, 500, err.Error())
}

return helper.Success(c, "Faculties fetched successfully", faculties)
}

func (cl *FacultyController) GetFacultyByIDController(c fiber.Ctx) error {
idStr := c.Params("id")
id, err := strconv.ParseUint(idStr, 10, 32)
if err != nil {
return helper.Error(c, 400, "invalid faculty ID")
}

faculty, err := cl.facultyService.GetFacultyServiceById(uint(id))
if err != nil {
return helper.Error(c, 404, err.Error())
}

return helper.Success(c, "Faculty fetched successfully", faculty)
}

func (cl *FacultyController) UpdateFacultyController(c fiber.Ctx) error {
idStr := c.Params("id")
id, err := strconv.ParseUint(idStr, 10, 32)
if err != nil {
return helper.Error(c, 400, "invalid id")
}

var faculty model.Faculty
if err := c.Bind().Body(&faculty); err != nil {
return helper.Error(c, 400, "invalid request body")
}

if err := cl.facultyService.UpdateFacultyService(uint(id), faculty.Name, faculty.Gender); err != nil {
return helper.Error(c, 400, err.Error())
}

return helper.Success(c, "Faculty updated successfully", faculty)
}

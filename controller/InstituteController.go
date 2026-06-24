package controller

import (
"backend_institutions/helper"
"backend_institutions/model"
"backend_institutions/services"
"strconv"

"github.com/gofiber/fiber/v3"
)

type InstituteController struct {
instituteService *services.InstituteService
}

func NewInstituteController() *InstituteController {
return &InstituteController{
instituteService: services.NewInstituteService(),
}
}

func (cl *InstituteController) CreateInstituteController(c fiber.Ctx) error {
var institute model.Institutions

if err := c.Bind().Body(&institute); err != nil {
return helper.Error(c, 400, "invalid request body")
}

if err := cl.instituteService.CreateInsituteService(&institute); err != nil {
return helper.Error(c, 400, err.Error())
}

return helper.Success(c, "Institution created successfully", institute)
}

func (cl *InstituteController) GetAllInstitutesController(c fiber.Ctx) error {
institutes, err := cl.instituteService.GetInstituteService()
if err != nil {
return helper.Error(c, 500, err.Error())
}

return helper.Success(c, "Institutes fetched successfully", institutes)
}

func (cl *InstituteController) GetInstituteByIDController(c fiber.Ctx) error {
idStr := c.Params("id")
id, err := strconv.ParseUint(idStr, 10, 32)
if err != nil {
return helper.Error(c, 400, "Invalid institute ID")
}

institute, err := cl.instituteService.GetInstituteServiceById(uint(id))
if err != nil {
return helper.Error(c, 404, err.Error())
}

return helper.Success(c, "Institute fetched successfully", institute)
}

func (cl *InstituteController) UpdateInstituteController(c fiber.Ctx) error {
idStr := c.Params("id")
id, err := strconv.ParseUint(idStr, 10, 32)
if err != nil {
return helper.Error(c, 400, "invalid id")
}

var institute model.Institutions
if err := c.Bind().Body(&institute); err != nil {
return helper.Error(c, 400, "invalid request body")
}

if err := cl.instituteService.UpdateInstitutionService(uint(id), institute.Name, institute.Institution_code, institute.State); err != nil {
return helper.Error(c, 400, err.Error())
}

return helper.Success(c, "Institution updated successfully", institute)
}

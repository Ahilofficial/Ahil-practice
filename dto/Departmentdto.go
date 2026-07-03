package dto

type CreateDepartmentDTO struct {
	DepartmentName string `json:"department_name"`
	InstitutionID  uint   `json:"institution_id"`
}

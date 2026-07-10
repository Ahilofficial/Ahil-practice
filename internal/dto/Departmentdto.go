package dto

import (
	"backend_institutions/internal/model"
	"errors"
	"strings"
	"time"
)

type CreateDepartmentDTO struct {
	DepartmentName string `json:"department_name"`
	InstitutionID  uint   `json:"institution_id"`
}

func (dto *CreateDepartmentDTO) Sanitize() {
	dto.DepartmentName = strings.TrimSpace(dto.DepartmentName)
}

func (dto *CreateDepartmentDTO) Validate() error {
	

	if dto.DepartmentName == "" {
		return errors.New("department name is required")
	}
	if dto.InstitutionID == 0 {
		return errors.New("institution id is required")
	}
	return nil
}

type UpdateDepartmentDTO struct {
	DepartmentName string `json:"department_name"`
}

func (dto *UpdateDepartmentDTO) Sanitize() {
	dto.DepartmentName = strings.TrimSpace(dto.DepartmentName)
}

func (dto *UpdateDepartmentDTO) Validate() error {
	

	if dto.DepartmentName == "" {
		return errors.New("department name is required")
	}
	return nil
}

type DepartmentResponseDTO struct {
	ID             uint                 `json:"id"`
	DepartmentName string               `json:"department_name"`
	InstitutionID  uint                 `json:"institution_id"`
	IsActive       bool                 `json:"isactive"`
	Faculties      []FacultyResponseDTO `json:"faculties"`
}

func ToDepartmentResponseDTO(dept *model.Department) DepartmentResponseDTO {
	faculties := make([]FacultyResponseDTO, len(dept.Faculties))

	for i, faculty := range dept.Faculties {
		faculties[i] = ToFacultyResponseDTO(&faculty)
	}

	return DepartmentResponseDTO{
		ID:             dept.ID,
		DepartmentName: dept.DepartmentName,
		InstitutionID:  dept.InstitutionID,
		IsActive:       dept.IsActive,
		Faculties:      faculties,
	}
}

func ToDepartmentResponseListDTO(depts []model.Department) []DepartmentResponseDTO {
	list := make([]DepartmentResponseDTO, len(depts))

	for i := range depts {
		list[i] = ToDepartmentResponseDTO(&depts[i])
	}

	return list
}

type DepartmentFlatRow struct {
	DeptID             uint
	DepartmentName     string
	InstitutionID      uint
	DeptActive         bool

	FacID          *uint
	FacName        *string
	FacGender      *string
	FacJoiningDate *time.Time
	FacActive      *bool

	StudID     *uint
	StudName   *string
	StudEmail  *string
	StudGender *string
	StudActive *bool

	FeeID          *uint
	FeePaymentMode *string
	FeeAmount      *float64
	FeeActive      *bool
}
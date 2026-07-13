package dto

import (
	"backend_institutions/internal/model"
	"errors"
	"github.com/jinzhu/copier"
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
	var dto DepartmentResponseDTO
	copier.Copy(&dto, dept)

	dto.Faculties = make([]FacultyResponseDTO, len(dept.Faculties))
	for i := range dept.Faculties {
		dto.Faculties[i] = ToFacultyResponseDTO(&dept.Faculties[i])
	}

	return dto
}

func ToDepartmentResponseListDTO(depts []model.Department) []DepartmentResponseDTO {
	list := make([]DepartmentResponseDTO, len(depts))

	for i := range depts {
		list[i] = ToDepartmentResponseDTO(&depts[i])
	}

	return list
}

type DepartmentFlatRow struct {
	DeptID         uint   `gorm:"column:dept_id"`
	DepartmentName string `gorm:"column:department_name"`
	InstitutionID  uint   `gorm:"column:institution_id"`
	DeptActive     bool   `gorm:"column:dept_active"`

	FacID          *uint      `gorm:"column:fac_id"`
	FacName        *string    `gorm:"column:fac_name"`
	FacGender      *string    `gorm:"column:fac_gender"`
	FacJoiningDate *time.Time `gorm:"column:fac_joining_date"`
	FacActive      *bool      `gorm:"column:fac_active"`

	StudID     *uint   `gorm:"column:stud_id"`
	StudName   *string `gorm:"column:stud_name"`
	StudEmail  *string `gorm:"column:stud_email"`
	StudGender *string `gorm:"column:stud_gender"`
	StudActive *bool   `gorm:"column:stud_active"`

	FeeID          *uint    `gorm:"column:fee_id"`
	FeePaymentMode *string  `gorm:"column:fee_payment_mode"`
	FeeAmount      *float64 `gorm:"column:fee_amount"`
	FeeActive      *bool    `gorm:"column:fee_active"`
}

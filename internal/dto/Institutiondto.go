package dto

import (
	"backend_institutions/internal/model"
	"errors"
	"strings"
	"time"
	"github.com/jinzhu/copier"
)

type CreateInstitutionDTO struct {
	Name            string `json:"name"`
	InstitutionCode string `json:"institution_code"`
	State           string `json:"state"`
}

func (dto *CreateInstitutionDTO) Sanitize() {
	dto.Name = strings.TrimSpace(dto.Name)
	dto.InstitutionCode = strings.TrimSpace(dto.InstitutionCode)
	dto.State = strings.TrimSpace(dto.State)
}

func (dto *CreateInstitutionDTO) Validate() error {


	if dto.Name == "" {
		return errors.New("name is required")
	}
	if dto.InstitutionCode == "" {
		return errors.New("institution code is required")
	}
	if dto.State == "" {
		return errors.New("state is required")
	}
	return nil
}

type UpdateInstitutionDTO struct {
	Name            string `json:"name"`
	InstitutionCode string `json:"institution_code"`
	State           string `json:"state"`
}

func (dto *UpdateInstitutionDTO) Sanitize() {
	dto.Name = strings.TrimSpace(dto.Name)
	dto.InstitutionCode = strings.TrimSpace(dto.InstitutionCode)
	dto.State = strings.TrimSpace(dto.State)
}

func (dto *UpdateInstitutionDTO) Validate() error {
	

	if dto.Name == "" {
		return errors.New("name is required")
	}
	if dto.InstitutionCode == "" {
		return errors.New("institution code is required")
	}
	if dto.State == "" {
		return errors.New("state is required")
	}
	return nil
}

type InstitutionResponseDTO struct {
	ID              uint                    `json:"id"`
	Name            string                  `json:"name"`
	InstitutionCode string                  `json:"institution_code"`
	State           string                  `json:"state"`
	IsActive        bool                    `json:"isactive"`
	Departments     []DepartmentResponseDTO `json:"departments"`
}

func ToInstitutionResponseDTO(inst *model.Institutions) InstitutionResponseDTO {
	var dto InstitutionResponseDTO
	copier.Copy(&dto, inst)

	dto.Departments = make([]DepartmentResponseDTO, len(inst.Departments))
	for i := range inst.Departments {
		dto.Departments[i] = ToDepartmentResponseDTO(&inst.Departments[i])
	}

	return dto
}

func ToInstitutionResponseListDTO(insts []model.Institutions) []InstitutionResponseDTO {
	list := make([]InstitutionResponseDTO, len(insts))

	for i := range insts {
		list[i] = ToInstitutionResponseDTO(&insts[i])
	}

	return list
}

type InstitutionFlatRow struct {
	InstID          uint   `gorm:"column:inst_id"`
	InstName        string `gorm:"column:inst_name"`
	InstitutionCode string `gorm:"column:institution_code"`
	InstState       string `gorm:"column:inst_state"`
	InstActive      bool   `gorm:"column:inst_active"`

	DeptID         *uint   `gorm:"column:dept_id"`
	DepartmentName *string `gorm:"column:department_name"`
	DeptActive     *bool   `gorm:"column:dept_active"`

	FacID          *uint      `gorm:"column:fac_id"`
	FacName        *string    `gorm:"column:fac_name"`
	FacGender      *string    `gorm:"column:fac_gender"`
	FacJoiningDate *time.Time `gorm:"column:fac_joining_date"`
	FacActive      *bool      `gorm:"column:fac_active"`

	StudID         *uint   `gorm:"column:stud_id"`
	StudName       *string `gorm:"column:stud_name"`
	StudEmail      *string `gorm:"column:stud_email"`
	StudGender     *string `gorm:"column:stud_gender"`
	StudActive     *bool   `gorm:"column:stud_active"`

	FeeID          *uint    `gorm:"column:fee_id"`
	FeePaymentMode *string  `gorm:"column:fee_payment_mode"`
	FeeAmount      *float64 `gorm:"column:fee_amount"`
	FeeActive      *bool    `gorm:"column:fee_active"`
}
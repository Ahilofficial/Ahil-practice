package dto

import (
	"backend_institutions/internal/model"
	"errors"
	"strings"
	"time"
	"github.com/jinzhu/copier"
)

type CreateFacultyDTO struct {
	Name         string    `json:"name"`
	Gender       string    `json:"gender"`
	JoiningDate  time.Time `json:"joining_date"`
	DepartmentID uint      `json:"department_id"`
}

func (dto *CreateFacultyDTO) Sanitize() {
	dto.Name = strings.TrimSpace(dto.Name)
	dto.Gender = strings.TrimSpace(strings.ToLower(dto.Gender))
}

func (dto *CreateFacultyDTO) Validate() error {
	

	if dto.Name == "" {
		return errors.New("name is required")
	}
	if dto.Gender == "" {
		return errors.New("gender is required")
	}
	if dto.JoiningDate.IsZero() {
		return errors.New("joining date is required")
	}
	if dto.DepartmentID == 0 {
		return errors.New("department id is required")
	}
	return nil
}

type UpdateFacultyDTO struct {
	Name   string `json:"name"`
	Gender string `json:"gender"`
}

func (dto *UpdateFacultyDTO) Sanitize() {
	dto.Name = strings.TrimSpace(dto.Name)
	dto.Gender = strings.TrimSpace(strings.ToLower(dto.Gender))
}

func (dto *UpdateFacultyDTO) Validate() error {
	
	if dto.Name == "" {
		return errors.New("name is required")
	}
	if dto.Gender == "" {
		return errors.New("gender is required")
	}
	return nil
}

type FacultyResponseDTO struct {
	ID           uint                 `json:"id"`
	Name         string               `json:"name"`
	Gender       string               `json:"gender"`
	JoiningDate  time.Time            `json:"joining_date"`
	DepartmentID uint                 `json:"department_id"`
	IsActive     bool                 `json:"isactive"`
	Students     []StudentResponseDTO `json:"students"`
}

func ToFacultyResponseDTO(fac *model.Faculty) FacultyResponseDTO {
	var dto FacultyResponseDTO
	copier.Copy(&dto, fac)

	dto.Students = make([]StudentResponseDTO, len(fac.Students))
	for i := range fac.Students {
		dto.Students[i] = ToStudentResponseDTO(&fac.Students[i])
	}

	return dto
}

func ToFacultyResponseListDTO(facs []model.Faculty) []FacultyResponseDTO {
	list := make([]FacultyResponseDTO, len(facs))

	for i := range facs {
		list[i] = ToFacultyResponseDTO(&facs[i])
	}

	return list
}


type FacultyFlatRow struct {
	FacultyID uint      `gorm:"column:faculty_id"`
	FacultyName string  `gorm:"column:faculty_name"`
	FacultyGender string `gorm:"column:faculty_gender"`
	JoiningDate time.Time `gorm:"column:joining_date"`
	DepartmentID uint `gorm:"column:department_id"`
	FacultyActive bool `gorm:"column:is_active"`

	StudentID *uint `gorm:"column:student_id"`
	StudentName *string `gorm:"column:student_name"`
	StudentEmail *string `gorm:"column:student_email"`
	StudentGender *string `gorm:"column:student_gender"`
	StudentActive *bool `gorm:"column:student_active"`

	FeeID *uint `gorm:"column:fee_id"`
	FeePaymentMode *string `gorm:"column:payment_mode"`
	FeeAmount *float64 `gorm:"column:amount"`
	FeeActive *bool `gorm:"column:fee_active"`
}
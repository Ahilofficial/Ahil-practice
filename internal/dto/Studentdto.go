package dto

import (


	"backend_institutions/internal/model"
	"errors"
	"strings"
)

type CreateStudentDTO struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
	FacultyID uint   `json:"faculty_id"`
}

func (dto *CreateStudentDTO) Sanitize() {
	dto.Name = strings.TrimSpace(dto.Name)
	dto.Email = strings.TrimSpace(strings.ToLower(dto.Email))
	dto.Gender = strings.TrimSpace(strings.ToLower(dto.Gender))
}

func (dto *CreateStudentDTO) Validate() error {


	if dto.Name == "" {
		return errors.New("name is required")
	}
	if dto.Email == "" {
		return errors.New("email is required")
	}
	if !emailRegex.MatchString(dto.Email) {
		return errors.New("invalid email format")
	}
	if dto.Gender == "" {
		return errors.New("gender is required")
	}
	if dto.FacultyID == 0 {
		return errors.New("faculty id is required")
	}
	return nil
}

type UpdateStudentDTO struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Gender string `json:"gender"`
}

func (dto *UpdateStudentDTO) Sanitize() {
	dto.Name = strings.TrimSpace(dto.Name)
	dto.Email = strings.TrimSpace(strings.ToLower(dto.Email))
	dto.Gender = strings.TrimSpace(strings.ToLower(dto.Gender))
}

func (dto *UpdateStudentDTO) Validate() error {
	

	if dto.Name == "" {
		return errors.New("name is required")
	}
	if dto.Email == "" {
		return errors.New("email is required")
	}
	if !emailRegex.MatchString(dto.Email) {
		return errors.New("invalid email format")
	}
	if dto.Gender == "" {
		return errors.New("gender is required")
	}
	return nil
}

type StudentResponseDTO struct {
	ID         uint             `json:"id"`
	Name       string           `json:"name"`
	Email      string           `json:"email"`
	Gender     string           `json:"gender"`
	FacultyID  uint             `json:"faculty_id"`
	IsActive   bool             `json:"isactive"`
	Fees        []FeesResponseDTO `json:"fees"`
}

func ToStudentResponseDTO(stud *model.Student) StudentResponseDTO {
	fees := make([]FeesResponseDTO, len(stud.Fees))

	for i, fee := range stud.Fees {
		fees[i] = ToFeesResponseDTO(&fee)
	}

	return StudentResponseDTO{
		ID:        stud.ID,
		Name:      stud.Name,
		Email:     stud.Email,
		Gender:    stud.Gender,
		FacultyID: stud.FacultyID,
		IsActive:  stud.IsActive,
		Fees:      fees,
	}
}

func ToStudentResponseListDTO(studs []StudentFlatRow) []StudentResponseDTO {
	list := make([]StudentResponseDTO, len(studs))

	for i := range studs {
		list[i] = ToStudentFlatRowResponseDTO(&studs[i])
	}

	return list
}
func ToStudentFlatRowResponseDTO(stud *StudentFlatRow) StudentResponseDTO {
	fees := []FeesResponseDTO{}

	if stud.FeeID != nil {
		fees = append(fees, FeesResponseDTO{
			ID:          *stud.FeeID,
			PaymentMode: *stud.FeePaymentMode,
			Amount:      *stud.FeeAmount,
			IsActive:    *stud.FeeActive,
			
		})
	}

	return StudentResponseDTO{
		ID:        stud.StudID,
		Name:      stud.StudName,
		Email:     stud.StudEmail,
		Gender:    stud.StudGender,
		FacultyID: stud.FacultyID,
		IsActive:  stud.StudActive,
		Fees:      fees,
	}
}

type StudentFlatRow struct {
	StudID     uint
	StudName   string
	StudEmail  string
	StudGender string
	FacultyID  uint
	StudActive bool

	FeeID          *uint
	FeePaymentMode *string
	FeeAmount      *float64
	FeeActive      *bool
}

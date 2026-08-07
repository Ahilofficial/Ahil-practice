package dto

type AssignInstitutionAdminDTO struct {
	UserID        uint `json:"user_id" validate:"required"`
	InstitutionID uint `json:"institution_id" validate:"required"`
}
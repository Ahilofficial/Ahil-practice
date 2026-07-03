package dto

type CreateInstitutionDTO struct {
	Name            string `json:"name" `
	InstitutionCode string `json:"institution_code" `
	State           string `json:"state"`
}

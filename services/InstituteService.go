package services

import (
	"backend_institutions/model"
	"backend_institutions/repository"
)

func CreateInsituteService(institute *model.Institutions) error {
	return repository.CreateInstitution(institute)
}

func GetInstituteService() ([]model.Institutions, error) {
	return repository.FetchInstitution()
}

func GetInstituteServicePaginated(page, limit int) ([]model.Institutions, int64, error) {
	return repository.FetchInstitutionPaginated(page, limit)
}

func GetInstituteServiceById(id uint) (model.Institutions, error) {
	return repository.FetchInstitutionById(id)
}

func GetInstituteServiceDeleted() ([]model.Institutions, error) {
	return repository.FetchInstitutionDeleted()
}

func DeleteInstitutionService(id uint) error {
	return repository.DeleteInstitution(id)
}

func GetActiveInstitute() (model.Institutions, error) {
	return repository.GetActiveInstitute()
}

func GetInactiveInstitute() (model.Institutions, error) {
	return repository.GetInactiveInstitute()
}

func UpdateInstitutionService(
	id uint,
	name string,
	institutionCode string,
	state string,
) error {
	institute, err := repository.FetchInstitutionById(id)
	if err != nil {
		return err
	}

	institute.Name = name
	institute.InstitutionCode = institutionCode
	institute.State = state

	return repository.UpdateInstitution(&institute)
}

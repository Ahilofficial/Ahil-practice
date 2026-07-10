package services

import (
	"backend_institutions/internal/dto"
	"backend_institutions/internal/model"
	"backend_institutions/internal/repository"
)

type InstituteService struct {
	instituterepo *repository.InstitutionRepository
}

func NewInstituteService(instituterepo *repository.InstitutionRepository) *InstituteService {
	return &InstituteService{instituterepo: instituterepo}
}

func (s *InstituteService) CreateInsituteService(dto *dto.CreateInstitutionDTO) (model.Institutions, error) {
	institute := model.Institutions{
		Name:            dto.Name,
		InstitutionCode: dto.InstitutionCode,
		State:           dto.State,
	}

	err := s.instituterepo.CreateInstitution(&institute)
	return institute, err
}

func (s *InstituteService) GetInstituteService() ([]model.Institutions, error) {
	return s.instituterepo.FetchInstitution()
}

func (s *InstituteService) GetInstituteServicePaginated(page, limit int) ([]model.Institutions, int64, error) {
	return s.instituterepo.FetchInstitutionPaginated(page, limit)
}

func (s *InstituteService) GetInstituteServiceById(id uint) (model.Institutions, error) {
	return s.instituterepo.FetchInstitutionById(id)
}

func (s *InstituteService) GetInstituteServiceDeleted() ([]model.Institutions, error) {
	return s.instituterepo.FetchInstitutionDeleted()
}

func (s *InstituteService) DeleteInstitutionService(id uint) error {
	return s.instituterepo.DeleteInstitution(id)
}

func (s *InstituteService) GetActiveInstitute() (model.Institutions, error) {
	return s.instituterepo.GetActiveInstitute()
}

func (s *InstituteService) GetInactiveInstitute() (model.Institutions, error) {
	return s.instituterepo.GetInactiveInstitute()
}

func (s *InstituteService) UpdateInstitutionService(id uint, dto *dto.UpdateInstitutionDTO) error {
	institute, err := s.instituterepo.FetchInstitutionById(id)
	if err != nil {
		return err
	}

	institute.Name = dto.Name
	institute.InstitutionCode = dto.InstitutionCode
	institute.State = dto.State

	return s.instituterepo.UpdateInstitution(&institute)
}

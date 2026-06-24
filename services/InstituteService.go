package services

import (
"backend_institutions/model"
"backend_institutions/repository"
)

type InstituteService struct {
instituterepo *repository.InstitutionRepository
}

func NewInstituteService() *InstituteService {
return &InstituteService{
instituterepo: repository.NewInstituteRepository(),
}
}

func (s *InstituteService) CreateInsituteService(institute *model.Institutions) error {
return s.instituterepo.CreateInstitution(institute)
}

func (s *InstituteService) GetInstituteService() ([]model.Institutions, error) {
return s.instituterepo.FetchInstitution()
}

func (s *InstituteService) GetInstituteServiceById(id uint) (model.Institutions, error) {
return s.instituterepo.FetchInstitutionById(id)
}

func (s *InstituteService) DeleteInstitutionService(id uint) error {
return s.instituterepo.DeleteInstitution(id)
}

func (s *InstituteService) UpdateInstitutionService(id uint, name string, institutionCode string, state string) error {
institute, err := s.instituterepo.FetchInstitutionById(id)
if err != nil {
return err
}

institute.Name = name
institute.Institution_code = institutionCode
institute.State = state

return s.instituterepo.UpdateInstitution(&institute)
}

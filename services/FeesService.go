package services

import (
"backend_institutions/model"
"backend_institutions/repository"
)

type FeesService struct {
feesrepo *repository.FeesRepository
}

func NewFeesService() *FeesService {
return &FeesService{
feesrepo: repository.NewFeesRepository(),
}
}

func (s *FeesService) CreateFeesService(fees *model.Fees) error {
return s.feesrepo.CreateFees(fees)
}

func (s *FeesService) GetFeesService() ([]model.Fees, error) {
return s.feesrepo.FetchFees()
}

func (s *FeesService) GetFeesServiceById(id uint) (model.Fees, error) {
return s.feesrepo.FetchFeesById(id)
}

func (s *FeesService) UpdateFeesService(id uint, payment_mode string, amount uint) error {
fees, err := s.feesrepo.FetchFeesById(id)
if err != nil {
return err
}

fees.Payment_mode = payment_mode
fees.Amount = amount

return s.feesrepo.UpdateFeesById(&fees)
}

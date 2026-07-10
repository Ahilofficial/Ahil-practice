package services

import (
	"backend_institutions/internal/model"
	"backend_institutions/internal/repository"
	"backend_institutions/internal/dto"
)

type FeesService struct {
	feesrepo *repository.FeesRepository
}

func NewFeesService(feesrepo *repository.FeesRepository) *FeesService {
	return &FeesService{feesrepo: feesrepo}
}

func (s *FeesService) CreateFeesService(dto *dto.CreateFeesDTO) (model.Fees, error) {
	fees := model.Fees{
		PaymentMode: dto.PaymentMode,
		Amount:      dto.Amount,
		StudentID:   dto.StudentID,
	}
	
	err := s.feesrepo.CreateFees(&fees)
	return fees, err
}

func (s *FeesService) GetFeesService() ([]model.Fees, error) {
	return s.feesrepo.FetchFees()
}

func (s *FeesService) GetFeesServicePaginated(page, limit int) ([]model.Fees, int64, error) {
	return s.feesrepo.FetchFeesPaginated(page, limit)
}

func (s *FeesService) GetFeesServiceById(id uint) (model.Fees, error) {
	return s.feesrepo.FetchFeesById(id)
}

func (s *FeesService) DeleteFeesService(id uint) error {
	return s.feesrepo.DeleteFees(id)
}

func (s *FeesService) GetInactiveFeesService() ([]model.Fees, error) {
	return s.feesrepo.FetchInactiveFees()
}

func (s *FeesService) UpdateFeesService(id uint, dto *dto.UpdateFeesDTO) error {
	fees, err := s.feesrepo.FetchFeesById(id)
	if err != nil {
		return err
	}
	fees.PaymentMode = dto.PaymentMode
	fees.Amount = dto.Amount
	return s.feesrepo.UpdateFeesById(&fees)
}

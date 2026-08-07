package services

import (
	"backend_institutions/internal/dto"
	"backend_institutions/internal/model"
	"backend_institutions/internal/repository"
	"errors"
	
)

type FeesService struct {
	feesrepo *repository.FeesRepository
}

func NewFeesService(feesrepo *repository.FeesRepository) *FeesService {
	return &FeesService{feesrepo: feesrepo}
}

func (s *FeesService) CreateFeesService(dto *dto.CreateFeesDTO) (model.Fees, error) {
	fees := model.Fees{
		PaymentMode:   dto.PaymentMode,
		TotalAmount:   dto.TotalAmount,
		TotalPaid:     0,
		PendingAmount: dto.TotalAmount,
		StudentID:     dto.StudentID,
	}

	err := s.feesrepo.CreateFees(&fees)
	if err != nil {
		return model.Fees{}, err
	}

	return fees, nil
}


func (s *FeesService) GetFeesService() ([]model.Fees, error) {
	return s.feesrepo.FetchFees()
}

func (s *FeesService) GetFeesServicePaginated(search string,page, limit int) ([]model.Fees, int64, error) {
	return s.feesrepo.FetchFeesPaginated(search,page, limit)
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
	fees.TotalAmount = dto.Amount

	return s.feesrepo.UpdateFeesById(&fees)
}

func (s *FeesService) GetPaymentByIDService(id uint) (model.Payment, error) {
	return s.feesrepo.FetchPaymentByID(id)
}

func (s *FeesService) GetPaymentByFeeIDService(feeID uint) ([]model.Payment, error) {
	return s.feesrepo.FetchPaymentByFeeID(feeID)
}

func (s *FeesService) CreatePaymentService(dto *dto.CreatePaymentDTO) (model.Payment, error) {

	
	fee, err := s.feesrepo.FetchFeesById(dto.FeeID)
	if err != nil {
		return model.Payment{}, err
	}

	
	if dto.AmountPaid <= 0 {
		return model.Payment{}, errors.New("payment amount must be greater than zero")
	}

	if dto.AmountPaid > fee.PendingAmount {
		return model.Payment{}, errors.New("payment amount exceeds pending amount")
	}

	payment := model.Payment{
		Month:       dto.Month,
		AmountPaid:  dto.AmountPaid,
		PaymentMode: dto.PaymentMode,
		FeeID:       dto.FeeID,
	}

	if err := s.feesrepo.CreatePayment(&payment); err != nil {
		return model.Payment{}, err
	}

	
	fee.TotalPaid += dto.AmountPaid
	fee.PendingAmount = fee.TotalAmount - fee.TotalPaid

	
	if err := s.feesrepo.UpdateFeesById(&fee); err != nil {
		return model.Payment{}, err
	}

	return payment, nil
}


func (s *FeesService) FetchFeesByStudentID(studentID uint) (*model.Fees, error) {
	return s.feesrepo.FetchFeesByStudentID(studentID)
}
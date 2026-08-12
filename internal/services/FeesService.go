package services

import (
	"errors"

	"backend_institutions/internal/dto"
	"backend_institutions/internal/model"
	"backend_institutions/internal/repository"
)

type FeesService struct {
	feesRepo    *repository.FeesRepository
	studentRepo *repository.StudentRepository
	userRepo    *repository.UserRepository
}

func NewFeesService(
	feesRepo *repository.FeesRepository,
	studentRepo *repository.StudentRepository,
	userRepo *repository.UserRepository,
) *FeesService {
	return &FeesService{
		feesRepo:    feesRepo,
		studentRepo: studentRepo,
		userRepo:    userRepo,
	}
}

func (s *FeesService) checkInstitutionAccess(
	userID uint,
	institutionID uint,
) error {

	hasAccess, err := s.userRepo.HasInstitutionAccess(
		userID,
		institutionID,
	)
	if err != nil {
		return err
	}

	if !hasAccess {
		return errors.New(
			"user does not have access to this institution",
		)
	}

	return nil
}

func (s *FeesService) CreateFeesService(
	userID uint,
	req *dto.CreateFeesDTO,
) (model.Fees, error) {

	institutionID, err := s.studentRepo.GetInstitutionByStudentID(
		req.StudentID,
	)
	if err != nil {
		return model.Fees{}, err
	}

	if err := s.checkInstitutionAccess(
		userID,
		institutionID,
	); err != nil {
		return model.Fees{}, err
	}

	fees := model.Fees{
		PaymentMode:   req.PaymentMode,
		TotalAmount:   req.TotalAmount,
		TotalPaid:     0,
		PendingAmount: req.TotalAmount,
		StudentID:     req.StudentID,
	}

	if err := s.feesRepo.CreateFees(&fees); err != nil {
		return model.Fees{}, err
	}

	return fees, nil
}

func (s *FeesService) GetFeesService() ([]model.Fees, error) {
	return s.feesRepo.FetchFees()
}

func (s *FeesService) GetFeesServicePaginated(
	search string,
	page int,
	limit int,
) ([]model.Fees, int64, error) {
	return s.feesRepo.FetchFeesPaginated(
		search,
		page,
		limit,
	)
}

func (s *FeesService) GetFeesServiceById(
	userID uint,
	id uint,
) (model.Fees, error) {

	fees, err := s.feesRepo.FetchFeesById(id)
	if err != nil {
		return model.Fees{}, err
	}

	institutionID, err := s.studentRepo.GetInstitutionByStudentID(
		fees.StudentID,
	)
	if err != nil {
		return model.Fees{}, err
	}

	if err := s.checkInstitutionAccess(
		userID,
		institutionID,
	); err != nil {
		return model.Fees{}, err
	}

	return fees, nil
}

func (s *FeesService) DeleteFeesService(
	userID uint,
	id uint,
) error {

	fees, err := s.feesRepo.FetchFeesById(id)
	if err != nil {
		return err
	}

	institutionID, err := s.studentRepo.GetInstitutionByStudentID(
		fees.StudentID,
	)
	if err != nil {
		return err
	}

	if err := s.checkInstitutionAccess(
		userID,
		institutionID,
	); err != nil {
		return err
	}

	return s.feesRepo.DeleteFees(id)
}

func (s *FeesService) GetInactiveFeesService() ([]model.Fees, error) {
	return s.feesRepo.FetchInactiveFees()
}

func (s *FeesService) UpdateFeesService(
	userID uint,
	id uint,
	req *dto.UpdateFeesDTO,
) error {

	fees, err := s.feesRepo.FetchFeesById(id)
	if err != nil {
		return err
	}

	institutionID, err := s.studentRepo.GetInstitutionByStudentID(
		fees.StudentID,
	)
	if err != nil {
		return err
	}

	if err := s.checkInstitutionAccess(
		userID,
		institutionID,
	); err != nil {
		return err
	}

	fees.PaymentMode = req.PaymentMode
	fees.TotalAmount = req.Amount
	fees.PendingAmount = fees.TotalAmount - fees.TotalPaid

	if fees.PendingAmount < 0 {
		return errors.New(
			"total amount cannot be less than total paid",
		)
	}

	return s.feesRepo.UpdateFeesById(&fees)
}

func (s *FeesService) GetPaymentByIDService(
	id uint,
) (model.Payment, error) {
	return s.feesRepo.FetchPaymentByID(id)
}

func (s *FeesService) GetPaymentByFeeIDService(
	feeID uint,
) ([]model.Payment, error) {
	return s.feesRepo.FetchPaymentByFeeID(feeID)
}

func (s *FeesService) CreatePaymentService(
	userID uint,
	req *dto.CreatePaymentDTO,
) (model.Payment, error) {

	fee, err := s.feesRepo.FetchFeesById(req.FeeID)
	if err != nil {
		return model.Payment{}, err
	}

	institutionID, err := s.studentRepo.GetInstitutionByStudentID(
		fee.StudentID,
	)
	if err != nil {
		return model.Payment{}, err
	}

	if err := s.checkInstitutionAccess(
		userID,
		institutionID,
	); err != nil {
		return model.Payment{}, err
	}

	if req.AmountPaid <= 0 {
		return model.Payment{}, errors.New(
			"payment amount must be greater than zero",
		)
	}

	if req.AmountPaid > fee.PendingAmount {
		return model.Payment{}, errors.New(
			"payment amount exceeds pending amount",
		)
	}

	payment := model.Payment{
		Month:       req.Month,
		AmountPaid:  req.AmountPaid,
		PaymentMode: req.PaymentMode,
		FeeID:       req.FeeID,
	}

	if err := s.feesRepo.CreatePayment(&payment); err != nil {
		return model.Payment{}, err
	}

	fee.TotalPaid += req.AmountPaid
	fee.PendingAmount = fee.TotalAmount - fee.TotalPaid

	if err := s.feesRepo.UpdateFeesById(&fee); err != nil {
		return model.Payment{}, err
	}

	return payment, nil
}

func (s *FeesService) FetchFeesByStudentID(
	studentID uint,
) (*model.Fees, error) {
	return s.feesRepo.FetchFeesByStudentID(studentID)
}

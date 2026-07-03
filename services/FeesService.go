package services

import (
	"backend_institutions/model"
	"backend_institutions/repository"
)



func  CreateFeesService(fees *model.Fees) error {
	return repository.CreateFees(fees)
}

func  GetFeesService() ([]model.Fees, error) {
	return repository.FetchFees()
}

func  GetFeesServicePaginated(page, limit int) ([]model.Fees, int64, error) {
	return repository.FetchFeesPaginated(page, limit)
}

func  GetFeesServiceById(id uint) (model.Fees, error) {
	return repository.FetchFeesById(id)
}



func  DeleteFeesService(id uint) error {
	return repository.DeleteFees(id)
}



func  GetInactiveFeesService() ([]model.Fees, error) {
	return repository.FetchInactiveFees()
}

func  UpdateFeesService(id uint, payment_mode string, amount uint) error {
	fees, err := repository.FetchFeesById(id)
	if err != nil {
		return err
	}
	fees.PaymentMode = payment_mode
	fees.Amount = amount
	return repository.UpdateFeesById(&fees)
}

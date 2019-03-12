package service

import (
	"github.com/taufanmahaputra/forex/pkg/repository"
	"log"
)

type RateService struct {
	rateRepository repository.RateRepositoryItf
}

func InitRateService(rateRepository repository.RateRepositoryItf) RateService {
	return RateService{
		rateRepository: rateRepository,
	}
}

func (rs RateService) CreateExchangeRate(rate *repository.ExchangeRate) error {
	err := rs.rateRepository.InsertExchangeRate(rate)

	if err != nil {
		log.Printf("[RateService - CreateExchangeRate] : %s", err)
		return err
	}

	return nil
}

func (rs RateService) DeleteExchangeRate(rate *repository.ExchangeRate) error {
	err := rs.rateRepository.DeleteExchangeRateById(rate)

	if err != nil {
		log.Printf("[RateService - DeleteExchangeRate] : %s", err)
		return err
	}

	return nil
}

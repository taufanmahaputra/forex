package service

import (
	"github.com/taufanmahaputra/forex/pkg/repository"
	"log"
)

type RateService struct {
	rateRepository repository.RateRepositoryItf
}

func InitRateService(rateRepository repository.RateRepository) RateService {
	return RateService{
		rateRepository: rateRepository,
	}
}

func (rs RateService) CreateExchangeRate(rate *repository.ExchangeRate) error {
	err := rs.rateRepository.InsertExchangeRate(rate)

	if err != nil {
		log.Printf("[Service - CreateExchangeRate] : %s", err)
		return err
	}
	return nil
}
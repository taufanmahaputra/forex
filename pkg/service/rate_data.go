package service

import (
	"github.com/taufanmahaputra/forex/pkg/repository"
	"log"
)

type RateDataService struct {
	rateDataRepository repository.RateDataRepositoryItf
}

func InitRateDataService(rateDataRepository repository.RateDataRepositoryItf) RateDataService {
	return RateDataService{
		rateDataRepository: rateDataRepository,
	}
}

func (rs RateDataService) CreateDailyExchangeRateData(rate *repository.ExchangeRate, data *repository.ExchangeRateData) error {
	err := rs.rateDataRepository.InsertDailyExchangeRateData(rate, data)

	if err != nil {
		log.Printf("[RateDataService - CreateDailyExchangeRate] : %s", err)
		return err
	}

	return nil
}

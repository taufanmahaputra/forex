package service

import (
	"github.com/taufanmahaputra/forex/pkg/repository"
	"log"
	"math"
)

type Map map[string]interface{}

type RateDataService struct {
	rateRepository     repository.RateRepositoryItf
	rateDataRepository repository.RateDataRepositoryItf
}

func InitRateDataService(rateRepository repository.RateRepositoryItf, rateDataRepository repository.RateDataRepositoryItf) RateDataService {
	return RateDataService{
		rateRepository:     rateRepository,
		rateDataRepository: rateDataRepository,
	}
}

func (rs RateDataService) CreateDailyExchangeRateData(rate *repository.ExchangeRate, data *repository.ExchangeRateData) error {
	rateId, err := rs.rateRepository.GetExchangeRateIdByCurrencyPair(rate)
	if err != nil {
		log.Printf("[RateDataService - GetTrendBySevenExchangeRateData] : %s", err)
		return err
	}

	err = rs.rateDataRepository.InsertDailyExchangeRateData(rateId, data)
	if err != nil {
		log.Printf("[RateDataService - CreateDailyExchangeRate] : %s", err)
		return err
	}

	return nil
}

func (rs RateDataService) GetTrendBySevenExchangeRateData(rate *repository.ExchangeRate) (map[string]interface{}, error) {
	rateId, err := rs.rateRepository.GetExchangeRateIdByCurrencyPair(rate)
	if err != nil {
		log.Printf("[RateDataService - GetTrendBySevenExchangeRateData] : %s", err)
		return nil, err
	}

	rateDataList, err := rs.rateDataRepository.GetSevenSpecificExchangeRateData(rateId)
	if err != nil {
		log.Printf("[RateDataService - GetTrendBySevenExchangeRateData] : %s", err)
		return nil, err
	}

	minRate := math.Inf(1)
	maxRate := math.Inf(-1)

	var resultData []Map
	var sumRate float64
	for _, rateData := range rateDataList {
		rate := rateData.Rate

		data := Map{
			"rate": rate,
			"date": rateData.ValidTime.Format("2006-01-02"),
		}

		sumRate += rate
		if rate > maxRate {
			maxRate = rate
		}

		if rate < minRate {
			minRate = rate
		}

		resultData = append(resultData, data)
	}

	result := Map{
		"average":  sumRate / float64(len(rateDataList)),
		"variance": maxRate - minRate,
		"data":     resultData,
	}

	return result, nil
}

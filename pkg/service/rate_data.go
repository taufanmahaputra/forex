package service

import (
	"errors"
	"fmt"
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
	rateId := rs.rateRepository.GetExchangeRateIdByCurrencyPair(rate)
	if rateId == 0 {
		log.Printf("[RateDataService - GetTrendBySevenExchangeRateData] : Exchange rate doesnt exist")
		return errors.New(fmt.Sprintf("Exchange rate doesnt exist"))
	}

	rate.Id = rateId

	err := rs.rateDataRepository.InsertDailyExchangeRateData(rate, data)
	if err != nil {
		log.Printf("[RateDataService - CreateDailyExchangeRate] : %s", err)
		return err
	}

	return nil
}

func (rs RateDataService) GetTrendBySevenExchangeRateData(rate *repository.ExchangeRate) (map[string]interface{}, error) {
	rateId := rs.rateRepository.GetExchangeRateIdByCurrencyPair(rate)
	if rateId == 0 {
		log.Printf("[RateDataService - GetTrendBySevenExchangeRateData] : Exchange rate doesnt exist")
		return nil, errors.New(fmt.Sprintf("Exchange rate doesnt exist"))
	}

	rate.Id = rateId

	rateDataList := rs.rateDataRepository.GetSevenSpecificExchangeRateData(rate)
	if rateDataList == nil {
		log.Printf("[RateDataService - GetTrendBySevenExchangeRateData] : No Data for specific exchange rate")
		return nil, errors.New(fmt.Sprintf("No Data for specific exchange rate"))
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

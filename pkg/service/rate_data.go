package service

import (
	"errors"
	"fmt"
	"github.com/taufanmahaputra/forex/pkg/repository"
	"log"
	"math"
	"time"
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
		log.Printf("[RateDataService - GetTrendBySevenExchangeRateData] : Exchange rate doesnt exist")
		return errors.New(fmt.Sprintf("Exchange rate doesnt exist"))
	}

	rate.Id = rateId

	err = rs.rateDataRepository.InsertDailyExchangeRateData(rate, data)
	if err != nil {
		log.Printf("[RateDataService - CreateDailyExchangeRate] : %s", err)
		return err
	}

	return nil
}

func (rs RateDataService) GetTrendBySevenExchangeRateData(rate *repository.ExchangeRate) (map[string]interface{}, error) {
	rateId, err := rs.rateRepository.GetExchangeRateIdByCurrencyPair(rate)
	if err != nil {
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

func (rs RateDataService) GetListTrackedExchangeRateDataByDate(date time.Time) ([]map[string]interface{}, error) {
	var trackedRateDataList []map[string]interface{}

	rateList := rs.rateRepository.GetExchangeRateList()
	if rateList == nil {
		log.Printf("[RateDataService - GetListTrackedExchangeRateDataByDate] : Cannot get exchange rate list")
		return nil, errors.New(fmt.Sprintf("Cannot get exchange rate list"))
	}

	for _, rate := range rateList {
		var averageRateData float64
		var err error

		trackedData := Map{
			"id":            rate.Id,
			"currency_from": rate.CurrencyFrom,
			"currency_to":   rate.CurrencyTo,
		}

		currRateData := new(repository.ExchangeRateData)
		currRateData.ExchangeRateId = rate.Id
		currRateData.ValidTime = date

		currRateData = rs.rateDataRepository.GetExchangeRateDataByExchangeRateIdAndDate(currRateData)

		if currRateData != nil {
			averageRateData, err = rs.rateDataRepository.GetSevenDaysAverageExchangeRateDataByExchangeRateIdAndDate(currRateData)
		}

		if currRateData == nil || err != nil || averageRateData == 0 {
			trackedData["rate"] = "insufficient data"
			trackedData["average"] = "insufficient data"
		} else {
			trackedData["rate"] = currRateData.Rate
			trackedData["average"] = averageRateData
		}

		trackedRateDataList = append(trackedRateDataList, trackedData)
	}

	return trackedRateDataList, nil
}

package service

import (
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
	rateID, err := rs.rateRepository.GetExchangeRateIDByCurrencyPair(rate)
	if err != nil {
		log.Printf("[RateDataService - GetTrendBySevenExchangeRateData] : Exchange rate doesnt exist")
		return fmt.Errorf("exchange rate doesnt exist")
	}

	rate.ID = rateID

	err = rs.rateDataRepository.InsertDailyExchangeRateData(rate, data)
	if err != nil {
		log.Printf("[RateDataService - CreateDailyExchangeRate] : %s", err)
		return err
	}

	return nil
}

func (rs RateDataService) GetTrendBySevenExchangeRateData(rate *repository.ExchangeRate) (map[string]interface{}, error) {
	rateID, err := rs.rateRepository.GetExchangeRateIDByCurrencyPair(rate)
	if err != nil {
		log.Printf("[RateDataService - GetTrendBySevenExchangeRateData] : Exchange rate doesnt exist")
		return nil, fmt.Errorf("exchange rate doesnt exist")
	}

	rate.ID = rateID

	rateDataList := rs.rateDataRepository.GetSevenSpecificExchangeRateData(rate)
	if rateDataList == nil {
		log.Printf("[RateDataService - GetTrendBySevenExchangeRateData] : No Data for specific exchange rate")
		return nil, fmt.Errorf("no data for specific exchange rate")
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
		return nil, fmt.Errorf("cannot get exchange rate list")
	}

	for _, rate := range rateList {
		var averageRateData float64
		var err error

		trackedData := Map{
			"id":            rate.ID,
			"currency_from": rate.CurrencyFrom,
			"currency_to":   rate.CurrencyTo,
		}

		currRateData := new(repository.ExchangeRateData)
		currRateData.ExchangeRateID = rate.ID
		currRateData.ValidTime = date

		currRateData = rs.rateDataRepository.GetExchangeRateDataByExchangeRateIDAndDate(currRateData)

		if currRateData != nil {
			averageRateData, err = rs.rateDataRepository.GetSevenDaysAverageExchangeRateDataByExchangeRateIDAndDate(currRateData)
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

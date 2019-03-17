package controller

import (
	"github.com/taufanmahaputra/forex/pkg/repository"
	"github.com/taufanmahaputra/forex/pkg/service"
	"log"
	"time"
)

type ExchangeRate struct {
	CurrencyFrom string `json:"currency_from" query:"currency_from"`
	CurrencyTo   string `json:"currency_to" query:"currency_to"`
}

type ExchangeRateData struct {
	Date         string  `json:"date"`
	CurrencyFrom string  `json:"currency_from"`
	CurrencyTo   string  `json:"currency_to"`
	Rate         float64 `json:"rate"`
}

type RateController struct {
	rateService     service.RateService
	rateDataService service.RateDataService
}

func InitRateController(rateService service.RateService, rateDataService service.RateDataService) *RateController {
	return &RateController{
		rateService:     rateService,
		rateDataService: rateDataService,
	}
}

func (rc *RateController) GetExchangeRateList() []repository.ExchangeRate {
	return rc.rateService.GetExchangeRateList()
}

func (rc *RateController) PutNewExchangeRate(rate ExchangeRate) error {
	newExchangeRate := repository.ExchangeRate{
		CurrencyFrom: rate.CurrencyFrom,
		CurrencyTo:   rate.CurrencyTo,
	}

	err := rc.rateService.CreateExchangeRate(&newExchangeRate)
	if err != nil {
		log.Printf("[RateController - PutNewExchangeRate] : %s", err)
		return err
	}

	return nil
}

func (rc *RateController) RemoveExchangeRateByID(id int64) error {
	exchangeRate := repository.ExchangeRate{
		ID: id,
	}

	err := rc.rateService.DeleteExchangeRate(&exchangeRate)
	if err != nil {
		log.Printf("[RateController - RemoveExchangeRateByID] : %s", err)
		return err
	}

	return nil
}

func (rc *RateController) PutNewDailyExchangeRateData(data ExchangeRateData) error {
	exchangeRate := repository.ExchangeRate{
		CurrencyFrom: data.CurrencyFrom,
		CurrencyTo:   data.CurrencyTo,
	}

	date, err := time.Parse("2006-01-02", data.Date)
	if err != nil {
		return err
	}

	exchangeRateData := repository.ExchangeRateData{
		ValidTime: date,
		Rate:      data.Rate,
	}

	err = rc.rateDataService.CreateDailyExchangeRateData(&exchangeRate, &exchangeRateData)
	if err != nil {
		log.Printf("[RateController - PutNewDailyExchangeRateData] : %s", err)
		return err
	}

	return nil
}

func (rc *RateController) FindTrendBySevenExchangeRateData(rate ExchangeRate) (map[string]interface{}, error) {
	exchangeRate := repository.ExchangeRate{
		CurrencyFrom: rate.CurrencyFrom,
		CurrencyTo:   rate.CurrencyTo,
	}

	trend, err := rc.rateDataService.GetTrendBySevenExchangeRateData(&exchangeRate)
	if err != nil {
		log.Printf("[RateController - FindTrendBySevenExchangeRateData] : %s", err)
		return nil, err
	}

	return trend, nil
}

func (rc *RateController) GetListTrackedExchangeRateData(date string) ([]map[string]interface{}, error) {
	dateParsed, _ := time.Parse("2006-01-02", date)

	return rc.rateDataService.GetListTrackedExchangeRateDataByDate(dateParsed)
}

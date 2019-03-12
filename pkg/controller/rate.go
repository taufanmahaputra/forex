package controller

import (
	"github.com/taufanmahaputra/forex/pkg/repository"
	"github.com/taufanmahaputra/forex/pkg/service"
	"log"
	"strconv"
	"time"
)

type ExchangeRate struct {
	CurrencyFrom string `json:"currency_from"`
	CurrencyTo   string `json:"currency_to"`
}

type ExchangeRateData struct {
	Date         string `json:"date"`
	CurrencyFrom string `json:"currency_from"`
	CurrencyTo   string `json:"currency_to"`
	Rate         string `json:"rate"`
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

func (rc *RateController) RemoveExchangeRateById(id int64) error {
	exchangeRate := repository.ExchangeRate{
		Id: id,
	}

	err := rc.rateService.DeleteExchangeRate(&exchangeRate)
	if err != nil {
		log.Printf("[RateController - RemoveExchangeRateById] : %s", err)
		return err
	}

	return nil
}

func (rc *RateController) PutNewDailyExchangeRateData(data ExchangeRateData) error {
	exchangeRate := repository.ExchangeRate{
		CurrencyFrom: data.CurrencyFrom,
		CurrencyTo:   data.CurrencyTo,
	}

	rate, err := strconv.ParseFloat(data.Rate, 64)
	if err != nil {
		return err
	}

	date, err := time.Parse("2006-01-02", data.Date)
	if err != nil {
		return err
	}

	exchangeRateData := repository.ExchangeRateData{
		ValidTime: date,
		Rate:      rate,
	}

	err = rc.rateDataService.CreateDailyExchangeRateData(&exchangeRate, &exchangeRateData)
	if err != nil {
		log.Printf("[RateController - PutNewDailyExchangeRateData] : %s", err)
		return err
	}

	return nil
}

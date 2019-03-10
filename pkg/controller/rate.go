package controller

import (
	"github.com/taufanmahaputra/forex/pkg/repository"
	"github.com/taufanmahaputra/forex/pkg/service"
	"log"
)

type ExchangeRate struct {
	CurrencyFrom string `json:"currency_from"`
	CurrencyTo string `json:"currency_to"`
}

type RateController struct {
	rateService service.RateService
}

func InitRateController(rateService service.RateService) *RateController  {
	return &RateController{
		rateService: rateService,
	}
}

func (rc *RateController) PutNewExchangeRate(rate ExchangeRate) error {
	newExchangeRate := repository.ExchangeRate{
		CurrencyFrom: rate.CurrencyFrom,
		CurrencyTo: rate.CurrencyTo,
	}

	err := rc.rateService.CreateExchangeRate(&newExchangeRate)
	if err != nil {
		log.Printf("[Controller - PutNewExchangeRate] : %s", err)
		return err
	}
	return nil
}
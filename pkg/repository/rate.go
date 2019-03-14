package repository

import (
	"github.com/jinzhu/gorm"
	"log"
)

type ExchangeRate struct {
	Id           int64
	CurrencyFrom string
	CurrencyTo   string
}

type RateRepositoryItf interface {
	GetExchangeRateList() ([]ExchangeRate)
	InsertExchangeRate(*ExchangeRate) error
	DeleteExchangeRateById(*ExchangeRate) error
	GetExchangeRateIdByCurrencyPair(*ExchangeRate) int64
}

type RateRepository struct {
	DB *gorm.DB
}

func (r RateRepository) GetExchangeRateList() []ExchangeRate {
	var rateList []ExchangeRate

	result := r.DB.Find(&rateList)
	if result.Error != nil {
		log.Printf("[RateRepository - GetExchangeRateList] : %s", result.Error)
		return nil
	}

	return rateList
}

func (r RateRepository) InsertExchangeRate(rate *ExchangeRate) error {
	result := r.DB.Create(rate)

	if result.Error != nil {
		log.Printf("[RateRepository - InsertExchangeRate] : %s", result.Error)
		return result.Error
	}

	return nil
}

func (r RateRepository) DeleteExchangeRateById(rate *ExchangeRate) error {
	result := r.DB.Delete(rate)

	if result.Error != nil {
		log.Printf("[RateRepository - DeleteExchangeRateById] : %s", result.Error)
		return result.Error
	}

	return nil
}

func (r RateRepository) GetExchangeRateIdByCurrencyPair(rate *ExchangeRate) (id int64) {
	err := r.DB.Raw("SELECT id "+
		"FROM exchange_rates "+
		"WHERE currency_from = ? AND currency_to = ?", rate.CurrencyFrom, rate.CurrencyTo).Row().Scan(&id)

	if err != nil {
		log.Printf("[RateRepository - GetExchangeRateIdByCurrencyPair] : %s", err)
		return 0
	}

	return id
}

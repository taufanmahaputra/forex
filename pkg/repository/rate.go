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
	InsertExchangeRate(*ExchangeRate) error
	DeleteExchangeRateById(*ExchangeRate) error
}

type RateRepository struct {
	DB *gorm.DB
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

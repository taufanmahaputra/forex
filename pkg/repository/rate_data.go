package repository

import (
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

type ExchangeRateData struct {
	Id             int64
	ExchangeRateId int64
	Rate           float64
	ValidTime      time.Time
}

type RateDataRepositoryItf interface {
	InsertDailyExchangeRateData(*ExchangeRate, *ExchangeRateData, ) error
}

type RateDataRepository struct {
	DB *gorm.DB
}

func (rd RateDataRepository) InsertDailyExchangeRateData(rate *ExchangeRate, data *ExchangeRateData) error {
	result := rd.DB.Exec("INSERT INTO exchange_rate_datas (exchange_rate_id, rate, valid_time) "+
		"VALUES ("+
		"	(SELECT id FROM exchange_rates WHERE currency_from = ? AND currency_to = ?), ?, ?"+
		")"+
		"ON CONFLICT (exchange_rate_id, valid_time)"+
		"DO UPDATE"+
		"	SET rate = EXCLUDED.rate;", rate.CurrencyFrom, rate.CurrencyTo, data.Rate, data.ValidTime)

	if result.Error != nil {
		log.Printf("[RateDataRepository - InsertExchangeRateData] : %s", result.Error)
		return result.Error
	}

	return nil
}

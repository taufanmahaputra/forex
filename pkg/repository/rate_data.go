package repository

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

type ExchangeRateData struct {
	ID             int64
	ExchangeRateID int64
	Rate           float64
	ValidTime      time.Time
}

type RateDataRepositoryItf interface {
	InsertDailyExchangeRateData(*ExchangeRate, *ExchangeRateData) error
	GetExchangeRateDataByExchangeRateIDAndDate(*ExchangeRateData) *ExchangeRateData
	GetSevenSpecificExchangeRateData(*ExchangeRate) []ExchangeRateData
	GetSevenDaysAverageExchangeRateDataByExchangeRateIDAndDate(*ExchangeRateData) (float64, error)
}

type RateDataRepository struct {
	DB *gorm.DB
}

func (rd RateDataRepository) InsertDailyExchangeRateData(rate *ExchangeRate, data *ExchangeRateData) error {
	result := rd.DB.Exec("INSERT INTO exchange_rate_datas (exchange_rate_id, rate, valid_time) "+
		"VALUES (?, ?, ?) "+
		"ON CONFLICT (exchange_rate_id, valid_time) "+
		"DO UPDATE"+
		"	SET rate = EXCLUDED.rate", rate.ID, data.Rate, data.ValidTime)

	if result.Error != nil {
		log.Printf("[RateDataRepository - InsertDailyExchangeRateData] : %s", result.Error)
		return result.Error
	}

	return nil
}

func (rd RateDataRepository) GetExchangeRateDataByExchangeRateIDAndDate(data *ExchangeRateData) *ExchangeRateData {
	result := rd.DB.Table("exchange_rate_datas").
		Where("exchange_rate_id = ? AND valid_time = ?", data.ExchangeRateID, data.ValidTime).Find(&data)

	if result.Error != nil {
		log.Printf("[RateDataRepository - GetExchangeRateDataByExchangeRateIDAndDate] : %s", result.Error)
		return nil
	}

	return data
}

func (rd RateDataRepository) GetSevenSpecificExchangeRateData(rate *ExchangeRate) []ExchangeRateData {
	var rateDataList []ExchangeRateData

	result := rd.DB.Raw("SELECT * "+
		"FROM exchange_rate_datas WHERE exchange_rate_id = ?"+
		"ORDER BY valid_time DESC LIMIT 7", rate.ID).Scan(&rateDataList)

	if result.Error != nil {
		log.Printf("[RateDataRepository - GetSevenSpecificExchangeRateData] : %s", result.Error)
		return nil
	}

	return rateDataList
}

func (rd RateDataRepository) GetSevenDaysAverageExchangeRateDataByExchangeRateIDAndDate(data *ExchangeRateData) (float64, error) {
	var averageRate float64

	date := data.ValidTime.Format("2006-01-02")
	queryStmt := fmt.Sprintf("SELECT (CASE WHEN COUNT(*) < 7 THEN 0 ELSE AVG(rate) END) AS average "+
		"FROM exchange_rate_datas "+
		"WHERE "+
		" exchange_rate_id = %d AND "+
		" valid_time BETWEEN (DATE '%s' - interval '6 days') AND (DATE '%s')", data.ExchangeRateID, date, date)

	result := rd.DB.Raw(queryStmt).Row()
	err := result.Scan(&averageRate)
	if err != nil {
		log.Printf("[RateDataRepository - GetSevenDaysAverageExchangeRateDataByExchangeRateIDAndDate] : %s", err)
		return 0, err
	}

	return averageRate, nil

}

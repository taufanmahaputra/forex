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
	InsertDailyExchangeRateData(int64, *ExchangeRateData, ) error
	GetSevenSpecificExchangeRateData(int64) ([]ExchangeRateData, error)
	GetSevenDaysAverageExchangeRateDataByExchangeRateIdAndDate(int64, string) (map[string]interface{}, error)
}

type RateDataRepository struct {
	DB *gorm.DB
}

func (rd RateDataRepository) InsertDailyExchangeRateData(rateId int64, data *ExchangeRateData) error {
	result := rd.DB.Exec("INSERT INTO exchange_rate_datas (exchange_rate_id, rate, valid_time) "+
		"VALUES (?, ?, ?) "+
		"ON CONFLICT (exchange_rate_id, valid_time) "+
		"DO UPDATE"+
		"	SET rate = EXCLUDED.rate", rateId, data.Rate, data.ValidTime)

	if result.Error != nil {
		log.Printf("[RateDataRepository - InsertExchangeRateData] : %s", result.Error)
		return result.Error
	}

	return nil
}

func (rd RateDataRepository) GetSevenSpecificExchangeRateData(rateId int64) ([]ExchangeRateData, error) {
	var rateDataList []ExchangeRateData

	result := rd.DB.Raw("SELECT * "+
		"FROM exchange_rate_datas WHERE exchange_rate_id = ?"+
		"ORDER BY valid_time DESC LIMIT 7", rateId).Scan(&rateDataList)

	if result.Error != nil {
		log.Printf("[RateDataRepository - GetSevenSpecificExchangeRateData] : %s", result.Error)
		return nil, result.Error
	}

	return rateDataList, nil
}

func (rd RateDataRepository) GetSevenDaysAverageExchangeRateDataByExchangeRateIdAndDate(rateId int64, date string) (map[string]interface{}, error) {
	var days int
	var averageRate float64

	result := rd.DB.Raw("SELECT COUNT(*) as days, AVG(rate) "+
		"FROM exchange_rate_datas "+
		"WHERE "+
		" exchange_rate_id = ? AND "+
		" valid_time BETWEEN (DATE ? - interval '6 days') "+
		"				 AND (DATE ?)", rateId, date, date).Row()

	err := result.Scan(&days, &averageRate)
	if err != nil {
		log.Printf("[RateDataRepository - GetSevenDaysAverageExchangeRateDataByExchangeRateIdAndDate] : %s", err)
		return nil, err
	}

	res := make(map[string]interface{})
	res["days"] = days
	res["average"] = averageRate

	return res, nil

}

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
	InsertDailyExchangeRateData(*ExchangeRate, *ExchangeRateData) error
	GetExchangeRateDataByExchangeRateIdAndDate(*ExchangeRateData) *ExchangeRateData
	GetSevenSpecificExchangeRateData(*ExchangeRate) []ExchangeRateData
	GetSevenDaysAverageExchangeRateDataByExchangeRateIdAndDate(*ExchangeRateData) map[string]interface{}
}

type RateDataRepository struct {
	DB *gorm.DB
}

func (rd RateDataRepository) InsertDailyExchangeRateData(rate *ExchangeRate, data *ExchangeRateData) error {
	result := rd.DB.Exec("INSERT INTO exchange_rate_datas (exchange_rate_id, rate, valid_time) "+
		"VALUES (?, ?, ?) "+
		"ON CONFLICT (exchange_rate_id, valid_time) "+
		"DO UPDATE"+
		"	SET rate = EXCLUDED.rate", rate.Id, data.Rate, data.ValidTime)

	if result.Error != nil {
		log.Printf("[RateDataRepository - InsertExchangeRateData] : %s", result.Error)
		return result.Error
	}

	return nil
}

func (rd RateDataRepository) GetExchangeRateDataByExchangeRateIdAndDate(data *ExchangeRateData) *ExchangeRateData {
	result := rd.DB.
		Where("exchange_rate_id = ? AND valid_time = ?", data.ExchangeRateId, data.ValidTime).Find(&data)
	if result.Error != nil {
		log.Printf("[RateRepository - GetExchangeRateList] : %s", result.Error)
		return nil
	}

	return data
}

func (rd RateDataRepository) GetSevenSpecificExchangeRateData(rate *ExchangeRate) []ExchangeRateData {
	var rateDataList []ExchangeRateData

	result := rd.DB.Raw("SELECT * "+
		"FROM exchange_rate_datas WHERE exchange_rate_id = ?"+
		"ORDER BY valid_time DESC LIMIT 7", rate.Id).Scan(&rateDataList)

	if result.Error != nil {
		log.Printf("[RateDataRepository - GetSevenSpecificExchangeRateData] : %s", result.Error)
		return nil
	}

	return rateDataList
}

func (rd RateDataRepository) GetSevenDaysAverageExchangeRateDataByExchangeRateIdAndDate(data *ExchangeRateData) map[string]interface{} {
	var days int
	var averageRate float64

	result := rd.DB.Raw("SELECT COUNT(*) as days, AVG(rate) "+
		"FROM exchange_rate_datas "+
		"WHERE "+
		" exchange_rate_id = ? AND "+
		" valid_time BETWEEN (DATE ? - interval '6 days') "+
		"				 AND (DATE ?)", data.ExchangeRateId, data.ValidTime, data.ValidTime).Row()

	err := result.Scan(&days, &averageRate)
	if err != nil {
		log.Printf("[RateDataRepository - GetSevenDaysAverageExchangeRateDataByExchangeRateIdAndDate] : %s", err)
		return nil
	}

	res := make(map[string]interface{})
	res["days"] = days
	res["average"] = averageRate

	return res

}

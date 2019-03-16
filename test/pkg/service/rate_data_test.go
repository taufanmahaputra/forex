package service

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/taufanmahaputra/forex/pkg/repository"
	"github.com/taufanmahaputra/forex/pkg/service"
	mock "github.com/taufanmahaputra/forex/test/pkg/repository"
	"testing"
)

type RateDataServiceTestSuite struct {
	suite.Suite
	assert             *assert.Assertions
	rateRepository     *mock.MockRateRepositoryItf
	rateDataRepository *mock.MockRateDataRepositoryItf
	rateDataService    service.RateDataService
}

func (suite *RateDataServiceTestSuite) SetupTest() {
	ctl := gomock.NewController(suite.T())
	defer ctl.Finish()

	suite.assert = assert.New(suite.T())

	suite.rateRepository = mock.NewMockRateRepositoryItf(ctl)
	suite.rateDataRepository = mock.NewMockRateDataRepositoryItf(ctl)
	suite.rateDataService = service.InitRateDataService(suite.rateRepository, suite.rateDataRepository)
}

func (suite *RateDataServiceTestSuite) TestCreateDailyExchangeRateDataShouldReturnSuccess() {
	var rate repository.ExchangeRate
	var rateData repository.ExchangeRateData

	suite.rateRepository.
		EXPECT().
		GetExchangeRateIDByCurrencyPair(&rate).
		Return(int64(1), nil)

	suite.rateDataRepository.
		EXPECT().
		InsertDailyExchangeRateData(&rate, &rateData).
		Return(nil)

	result := suite.rateDataService.CreateDailyExchangeRateData(&rate, &rateData)

	suite.assert.Nil(result)
}

func (suite *RateDataServiceTestSuite) TestCreateDailyExchangeRateDataWithRateIDNotFoundShouldReturnError() {
	var rate repository.ExchangeRate
	var rateData repository.ExchangeRateData

	suite.rateRepository.
		EXPECT().
		GetExchangeRateIDByCurrencyPair(&rate).
		Return(int64(0), errors.New(""))

	result := suite.rateDataService.CreateDailyExchangeRateData(&rate, &rateData)

	suite.assert.EqualError(result, "exchange rate doesnt exist")
}

func (suite *RateDataServiceTestSuite) TestCreateDailyExchangeRateDataWithFailedInsertToDBShouldReturnError() {
	var rate repository.ExchangeRate
	var rateData repository.ExchangeRateData

	suite.rateRepository.
		EXPECT().
		GetExchangeRateIDByCurrencyPair(&rate).
		Return(int64(1), nil)

	suite.rateDataRepository.
		EXPECT().
		InsertDailyExchangeRateData(&rate, &rateData).
		Return(errors.New("error"))

	result := suite.rateDataService.CreateDailyExchangeRateData(&rate, &rateData)

	suite.assert.EqualError(result, "error")
}

func (suite *RateDataServiceTestSuite) TestGetTrendBySevenExchangeRateDataShouldReturnTrend() {
	var resultExpectedType map[string]interface{}

	var rate repository.ExchangeRate
	var rateDataList []repository.ExchangeRateData

	suite.rateRepository.
		EXPECT().
		GetExchangeRateIDByCurrencyPair(&rate).
		Return(int64(1), nil)

	rateDataList = []repository.ExchangeRateData{}
	suite.rateDataRepository.
		EXPECT().
		GetSevenSpecificExchangeRateData(&rate).
		Return(rateDataList)

	result, err := suite.rateDataService.GetTrendBySevenExchangeRateData(&rate)

	suite.assert.NoError(err)
	suite.assert.NotNil(result)
	suite.assert.IsType(resultExpectedType, result)
	suite.assert.Contains(result, "average")
	suite.assert.Contains(result, "variance")
	suite.assert.Contains(result, "data")
}

func (suite *RateDataServiceTestSuite) TestGetTrendBySevenExchangeRateDataWithRateIDNotFoundShouldReturnTrend() {
	var rate repository.ExchangeRate

	suite.rateRepository.
		EXPECT().
		GetExchangeRateIDByCurrencyPair(&rate).
		Return(int64(0), errors.New(""))

	result, err := suite.rateDataService.GetTrendBySevenExchangeRateData(&rate)

	suite.assert.Nil(result)
	suite.assert.EqualError(err, "exchange rate doesnt exist")
}

func (suite *RateDataServiceTestSuite) TestGetTrendBySevenExchangeRateDataWithNoRateDataOrDBErrorShouldReturnTrend() {
	var rate repository.ExchangeRate

	suite.rateRepository.
		EXPECT().
		GetExchangeRateIDByCurrencyPair(&rate).
		Return(int64(1), nil)

	suite.rateDataRepository.
		EXPECT().
		GetSevenSpecificExchangeRateData(&rate).
		Return(nil)

	result, err := suite.rateDataService.GetTrendBySevenExchangeRateData(&rate)

	suite.assert.Nil(result)
	suite.assert.EqualError(err, "no data for specific exchange rate")
}

func TestRateDataServiceTestSuite(t *testing.T) {
	suite.Run(t, new(RateDataServiceTestSuite))
}

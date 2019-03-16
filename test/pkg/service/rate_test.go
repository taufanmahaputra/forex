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

type RateServiceTestSuite struct {
	suite.Suite
	assert         *assert.Assertions
	rateRepository *mock.MockRateRepositoryItf
	rateService    service.RateService
}

func (suite *RateServiceTestSuite) SetupTest() {
	ctl := gomock.NewController(suite.T())
	defer ctl.Finish()

	suite.assert = assert.New(suite.T())

	suite.rateRepository = mock.NewMockRateRepositoryItf(ctl)
	suite.rateService = service.InitRateService(suite.rateRepository)
}

func (suite *RateServiceTestSuite) TestGetExchangeListShouldReturnExchangeList() {
	suite.rateRepository.
		EXPECT().
		GetExchangeRateList().
		Return([]repository.ExchangeRate{})

	result := suite.rateService.GetExchangeRateList()

	suite.assert.NotNil(result)
	suite.assert.Equal(0, len(result))
}

func (suite *RateServiceTestSuite) TestGetExchangeListWithDatabaseErrorShouldReturnNil() {
	suite.rateRepository.
		EXPECT().
		GetExchangeRateList().
		Return(nil)

	result := suite.rateService.GetExchangeRateList()

	suite.assert.Nil(result)
}

func (suite *RateServiceTestSuite) TestCreateExchangeRateShouldReturnSuccess() {
	rate := repository.ExchangeRate{}

	suite.rateRepository.
		EXPECT().
		InsertExchangeRate(&rate).
		Return(nil)

	result := suite.rateService.CreateExchangeRate(&rate)

	suite.assert.Nil(result)
}

func (suite *RateServiceTestSuite) TestCreateExchangeRateWithDatabaseErrorShouldReturnError() {
	rate := repository.ExchangeRate{}

	suite.rateRepository.
		EXPECT().
		InsertExchangeRate(&rate).
		Return(errors.New("error"))

	result := suite.rateService.CreateExchangeRate(&rate)

	suite.assert.EqualError(result, "error")
}

func (suite *RateServiceTestSuite) TestDeleteExchangeRateShouldReturnSuccess() {
	rate := repository.ExchangeRate{}

	suite.rateRepository.
		EXPECT().
		DeleteExchangeRateByID(&rate).
		Return(nil)

	result := suite.rateService.DeleteExchangeRate(&rate)

	suite.assert.Nil(result)
}

func (suite *RateServiceTestSuite) TestDeleteExchangeRateWithDatabaseErrorShouldReturnError() {
	rate := repository.ExchangeRate{}

	suite.rateRepository.
		EXPECT().
		DeleteExchangeRateByID(&rate).
		Return(errors.New("error"))

	result := suite.rateService.DeleteExchangeRate(&rate)

	suite.assert.EqualError(result, "error")
}

func TestRateServiceTestSuite(t *testing.T) {
	suite.Run(t, new(RateServiceTestSuite))
}

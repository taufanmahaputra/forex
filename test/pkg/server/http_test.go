package server

import (
	"github.com/Selvatico/go-mocket"
	mocket "github.com/Selvatico/go-mocket"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/suite"
	"github.com/taufanmahaputra/forex/pkg/server"
	"gopkg.in/gavv/httpexpect.v1"
	"net/http"
	"testing"
	"time"
)

type HTTPHandlerTestSuite struct {
	suite.Suite
	client *httpexpect.Expect
}

func (suite *HTTPHandlerTestSuite) SetupTest() {
	e := echo.New()

	httpServer := server.NewHTTPServer()
	httpServer.RegisterHandler(e)

	mocket.Catcher.Register()
	mocket.Catcher.Logging = true
	db, _ := gorm.Open(mocket.DriverName, "test")

	server.Init(db)

	suite.client = httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(e),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(suite.T()),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(suite.T(), true),
		},
	})
}

func (suite *HTTPHandlerTestSuite) TestIndex() {
	body := suite.client.
		GET("/").
		Expect().
		Status(http.StatusOK).JSON().Object()

	body.Keys().ContainsOnly("title", "message")
	body.Value("title").String().Equal("Foreign exchange rate API")
	body.Value("message").String().Equal("OK")
}

func (suite *HTTPHandlerTestSuite) TestHandleGetRateListShouldReturnOK() {
	reply := []map[string]interface{}{{"id": 1, "currency_from": "IDR", "currency_to": "SGD"}}

	gomocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "exchange_rates"`).WithReply(reply)

	body := suite.client.
		GET("/api/v1/rate").
		Expect().
		Status(http.StatusOK).JSON().Array()

	body.First().Object().Keys().ContainsOnly("id", "currency_from", "currency_to")
	body.First().Object().Value("id").Number().Equal(1)
	body.First().Object().Value("currency_from").String().Equal("IDR")
	body.First().Object().Value("currency_to").String().Equal("SGD")
}

func (suite *HTTPHandlerTestSuite) TestHandleGetRateListWithRateListIsNilShouldReturnError() {
	gomocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "exchange_rates"`).WithQueryException()

	body := suite.client.
		GET("/api/v1/rate").
		Expect().
		Status(http.StatusInternalServerError).JSON().Object()

	body.Keys().ContainsOnly("title", "message")
	body.Value("title").String().Equal("")
	body.Value("message").String().Equal("Internal server error")
}

func (suite *HTTPHandlerTestSuite) TestHandleNewRateShouldReturnCreated() {
	rate := map[string]interface{}{
		"currency_from": "SGD",
		"currency_to":   "USD",
	}

	body := suite.client.
		POST("/api/v1/rate").WithJSON(rate).
		Expect().
		Status(http.StatusCreated).JSON().Object()

	body.Keys().ContainsOnly("currency_from", "currency_to")
	body.Value("currency_from").String().Equal("SGD")
	body.Value("currency_to").String().Equal("USD")
}

func (suite *HTTPHandlerTestSuite) TestHandleNewRateWithWrongPayloadShouldReturnBadRequest() {
	rate := map[string]interface{}{
		"currency_from": 1,
	}

	body := suite.client.
		POST("/api/v1/rate").WithJSON(rate).
		Expect().
		Status(http.StatusBadRequest).JSON().Object()

	body.Keys().ContainsOnly("title", "message")
	body.Value("title").String().Equal("")
	body.Value("message").String().Equal("Cannot process your request")
}

func (suite *HTTPHandlerTestSuite) TestHandleNewRateWithDBErrorShouldReturnError() {
	rate := map[string]interface{}{
		"currency_from": "SGD",
		"currency_to":   "USD",
	}

	gomocket.Catcher.Reset().NewMock().
		WithQuery(`INSERT INTO "exchange_rates" ("currency_from","currency_to")`).
		WithExecException()

	gomocket.Catcher.NewMock().
		WithQuery(`INSERT  INTO "exchange_rates" ("currency_from","currency_to")`).
		WithExecException()

	body := suite.client.
		POST("/api/v1/rate").WithJSON(rate).
		Expect().
		Status(http.StatusInternalServerError).JSON().Object()

	body.Keys().ContainsOnly("title", "message")
	body.Value("title").String().Equal("")
	body.Value("message").String().Equal("Internal server error")
}

func (suite *HTTPHandlerTestSuite) TestHandleDeleteRateByIDShouldReturnSuccess() {
	suite.client.
		DELETE("/api/v1/rate/{id}", 1).
		Expect().
		Status(http.StatusNoContent)

}

func (suite *HTTPHandlerTestSuite) TestHandleDeleteRateByIDWithIDIsNotFoundShouldReturnError() {
	gomocket.Catcher.Reset().NewMock().
		WithQuery(`DELETE FROM "exchange_rates"  WHERE`).
		WithExecException()

	body := suite.client.
		DELETE("/api/v1/rate/{id}", 1).
		Expect().
		Status(http.StatusInternalServerError).JSON().Object()

	body.Keys().ContainsOnly("title", "message")
	body.Value("title").String().Equal("")
	body.Value("message").String().Equal("Internal server error")
}

func (suite *HTTPHandlerTestSuite) TestHandleNewDailyRateDataShouldReturnCreated() {
	rateData := map[string]interface{}{
		"date":          "2018-03-11",
		"currency_from": "IDR",
		"currency_to":   "USD",
		"rate":          0.6,
	}

	replyWithId := []map[string]interface{}{{"id": 1}}
	gomocket.Catcher.Reset().NewMock().
		WithQuery(`SELECT id FROM exchange_rates WHERE `).
		WithReply(replyWithId)

	body := suite.client.
		POST("/api/v1/rate/input_daily").WithJSON(rateData).
		Expect().
		Status(http.StatusCreated).JSON().Object()

	body.Keys().ContainsOnly("date", "currency_from", "currency_to", "rate")
	body.Value("date").String().Equal("2018-03-11")
	body.Value("currency_from").String().Equal("IDR")
	body.Value("currency_to").String().Equal("USD")
	body.Value("rate").Number().Equal(0.6)
}

func (suite *HTTPHandlerTestSuite) TestHandleNewDailyRateDataWithWrongPayloadShouldReturnBadRequest() {
	rateData := map[string]interface{}{
		"date":          1,
		"currency_from": "",
		"currency_to":   "",
		"rate":          "",
	}

	body := suite.client.
		POST("/api/v1/rate/input_daily").WithJSON(rateData).
		Expect().
		Status(http.StatusBadRequest).JSON().Object()

	body.Keys().ContainsOnly("title", "message")
	body.Value("title").String().Equal("")
	body.Value("message").String().Equal("Cannot process your request")
}

func (suite *HTTPHandlerTestSuite) TestHandleNewDailyRateDataWithRateIDNotExistShouldReturnError() {
	rateData := map[string]interface{}{
		"date":          "2018-03-11",
		"currency_from": "IDR",
		"currency_to":   "USD",
		"rate":          0.6,
	}

	gomocket.Catcher.Reset()

	body := suite.client.
		POST("/api/v1/rate/input_daily").WithJSON(rateData).
		Expect().
		Status(http.StatusInternalServerError).JSON().Object()

	body.Keys().ContainsOnly("title", "message")
	body.Value("title").String().Equal("")
	body.Value("message").String().Equal("exchange rate doesnt exist")
}

func (suite *HTTPHandlerTestSuite) TestHandleNewDailyRateDataWithDBErrorShouldReturnError() {
	rateData := map[string]interface{}{
		"date":          "2018-03-11",
		"currency_from": "IDR",
		"currency_to":   "USD",
		"rate":          0.6,
	}

	replyWithId := []map[string]interface{}{{"id": 1}}
	gomocket.Catcher.Reset().NewMock().
		WithQuery(`SELECT id FROM exchange_rates WHERE `).
		WithReply(replyWithId)

	gomocket.Catcher.NewMock().
		WithQuery(`INSERT INTO exchange_rate_datas (exchange_rate_id, rate, valid_time)`).
		WithExecException()

	body := suite.client.
		POST("/api/v1/rate/input_daily").WithJSON(rateData).
		Expect().
		Status(http.StatusInternalServerError).JSON().Object()

	body.Keys().ContainsOnly("title", "message")
	body.Value("title").String().Equal("")
	body.Value("message").String().Equal("driver: bad connection")
}

func (suite *HTTPHandlerTestSuite) TestHandleGetTrendBySevenExchangeRateDataShouldReturnOK() {
	replyWithId := []map[string]interface{}{{"id": 1}}
	gomocket.Catcher.Reset().NewMock().
		WithQuery(`SELECT id FROM exchange_rates WHERE `).
		WithReply(replyWithId)

	replyRateData := []map[string]interface{}{
		{
			"id":               1,
			"exchange_rate_id": 1,
			"rate":             99,
			"valid_time":       time.Now(),
		},
	}

	gomocket.Catcher.NewMock().
		WithQuery(`SELECT * FROM exchange_rate_datas WHERE `).
		WithReply(replyRateData)

	body := suite.client.GET("/api/v1/rate/trend").
		Expect().
		Status(http.StatusOK).JSON().Object()

	body.Keys().ContainsOnly("average", "data", "variance")
	body.Value("average").Number().Equal(99)
	body.Value("variance").Number().Equal(0)

	data := body.Value("data").Array().First().Object()
	data.Keys().ContainsOnly("rate", "date")
	data.Value("rate").Number().Equal(99)
	data.Value("date").String().Equal(time.Now().Format("2006-01-02"))
}

func (suite *HTTPHandlerTestSuite) TestHandleGetTrendBySevenExchangeRateDataShouldReturnBadRequest() {
	gomocket.Catcher.Reset()

	body := suite.client.GET("/api/v1/rate/trend").
		WithText("test").
		Expect().
		Status(http.StatusBadRequest).JSON().Object()

	body.Keys().ContainsOnly("title", "message")
	body.Value("title").String().Equal("")
	body.Value("message").String().Equal("Cannot process your request")
}

func (suite *HTTPHandlerTestSuite) TestHandleGetTrendBySevenExchangeRateDataWithRateIDIsNilShouldReturnError() {
	gomocket.Catcher.Reset()

	body := suite.client.GET("/api/v1/rate/trend").
		Expect().
		Status(http.StatusInternalServerError).JSON().Object()

	body.Keys().ContainsOnly("title", "message")
	body.Value("title").String().Equal("")
	body.Value("message").String().Equal("exchange rate doesnt exist")
}

func (suite *HTTPHandlerTestSuite) TestHandleGetTrendBySevenExchangeRateDataWithRateDataIsNilShouldReturnError() {
	replyWithId := []map[string]interface{}{{"id": 1}}
	gomocket.Catcher.Reset().NewMock().
		WithQuery(`SELECT id FROM exchange_rates WHERE`).
		WithReply(replyWithId)

	gomocket.Catcher.NewMock().
		WithQuery(`SELECT * FROM exchange_rate_datas WHERE `).
		WithQueryException()

	body := suite.client.GET("/api/v1/rate/trend").
		Expect().
		Status(http.StatusInternalServerError).JSON().Object()

	body.Keys().ContainsOnly("title", "message")
	body.Value("title").String().Equal("")
	body.Value("message").String().Equal("no data for specific exchange rate")
}

func (suite *HTTPHandlerTestSuite) TestHandleGetListTrackedExchangeRateDataShouldReturnOK() {
	replyGetExchageRateList := []map[string]interface{}{{"id": 1, "currency_from": "IDR", "currency_to": "SGD"}}

	gomocket.Catcher.Reset().NewMock().
		WithQuery(`SELECT * FROM "exchange_rates"`).
		WithReply(replyGetExchageRateList)

	body := suite.client.
		GET("/api/v1/rate/track").WithQuery("date", "2018-01-01").
		Expect().
		Status(http.StatusOK).JSON().Array()

	body.First().Object().Keys().ContainsOnly("id", "currency_from", "currency_to", "rate", "average")
}

func (suite *HTTPHandlerTestSuite) TestHandleGetListTrackedExchangeRateDataWithRateListIsNilShouldReturnError() {
	gomocket.Catcher.Reset().NewMock().
		WithQuery(`SELECT * FROM "exchange_rates"`).
		WithQueryException()

	body := suite.client.
		GET("/api/v1/rate/track").
		Expect().
		Status(http.StatusInternalServerError).JSON().Object()

	body.Keys().ContainsOnly("title", "message")
	body.Value("title").String().Equal("")
	body.Value("message").String().Equal("cannot get exchange rate list")
}

func TestHTTPHanlderTestSuite(t *testing.T) {
	suite.Run(t, new(HTTPHandlerTestSuite))
}

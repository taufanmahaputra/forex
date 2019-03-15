package server

import (
	"github.com/taufanmahaputra/forex/pkg/controller"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type Response struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type HTTPService struct {
}

func NewHTTPServer() HTTPService {
	return HTTPService{}
}

func (s HTTPService) RegisterHandler(e *echo.Echo) {
	e.GET("/", index)

	api := e.Group("/api")
	apiV1 := api.Group("/v1")

	//TODO: validate request payload
	apiV1.GET("/rate", handleGetRateList)
	apiV1.POST("/rate", handleNewRate)
	apiV1.DELETE("/rate/:id", handleDeleteRateByID)

	apiV1.POST("/rate/input_daily", handleNewDailyRateData)
	apiV1.GET("/rate/trend", handleGetTrendBySevenExchangeRateData)
	apiV1.GET("/rate/track", handleGetListTrackedExchangeRateData)
}

func index(ctx echo.Context) error {
	return response(ctx, http.StatusOK, Response{
		"Foreign exchange rate API",
		"OK",
	})
}

func handleGetRateList(ctx echo.Context) error {
	exchangeRateList := rateController.GetExchangeRateList()
	if exchangeRateList == nil {
		return response(ctx, http.StatusInternalServerError, Response{Message: "Internal server error"})
	}

	return response(ctx, http.StatusOK, exchangeRateList)
}

func handleNewRate(ctx echo.Context) error {
	rate := new(controller.ExchangeRate)
	if err := ctx.Bind(rate); err != nil {
		return response(ctx, http.StatusBadRequest, Response{Message: "Cannot process your request"})
	}

	err := rateController.PutNewExchangeRate(*rate)
	if err != nil {
		return response(ctx, http.StatusInternalServerError, Response{Message: "Internal server error"})
	}

	return response(ctx, http.StatusCreated, rate)
}

func handleDeleteRateByID(ctx echo.Context) error {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)

	err := rateController.RemoveExchangeRateByID(id)
	if err != nil {
		return response(ctx, http.StatusInternalServerError, Response{Message: "Internal server error"})
	}

	return response(ctx, http.StatusNoContent, nil)
}

func handleNewDailyRateData(ctx echo.Context) error {
	rateData := new(controller.ExchangeRateData)
	if err := ctx.Bind(rateData); err != nil {
		return response(ctx, http.StatusBadRequest, Response{Message: "Cannot process your request"})
	}

	err := rateController.PutNewDailyExchangeRateData(*rateData)
	if err != nil {
		return response(ctx, http.StatusInternalServerError, Response{Message: err.Error()})
	}

	return response(ctx, http.StatusCreated, rateData)
}

func handleGetTrendBySevenExchangeRateData(ctx echo.Context) error {
	rate := new(controller.ExchangeRate)
	if err := ctx.Bind(rate); err != nil {
		return response(ctx, http.StatusBadRequest, Response{Message: "Cannot process your request"})
	}

	trend, err := rateController.FindTrendBySevenExchangeRateData(*rate)
	if err != nil {
		return response(ctx, http.StatusInternalServerError, Response{Message: err.Error()})
	}

	return response(ctx, http.StatusOK, trend)
}

func handleGetListTrackedExchangeRateData(ctx echo.Context) error {
	date := ctx.QueryParam("date")

	trackedList, err := rateController.GetListTrackedExchangeRateData(date)
	if err != nil {
		return response(ctx, http.StatusInternalServerError, Response{Message: err.Error()})
	}

	return response(ctx, http.StatusOK, trackedList)
}

func response(ctx echo.Context, statusCode int, response interface{}) error {
	return ctx.JSON(statusCode, response)
}

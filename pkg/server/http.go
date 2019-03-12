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

type HttpService struct {
}

func NewHttpServer() HttpService {
	return HttpService{}
}

func (s HttpService) RegisterHandler(e *echo.Echo) {
	e.GET("/", index)

	api := e.Group("/api")
	apiV1 := api.Group("/v1")

	//TODO: validate request payload
	apiV1.POST("/rate", handleNewRate)
	apiV1.DELETE("/rate/:id", handleDeleteRateById)

	apiV1.POST("/rate/input_daily", handleNewDailyRateData)
}

func index(ctx echo.Context) error {
	return response(ctx, http.StatusOK, Response{
		"Foreign exchange rate API",
		"OK",
	})
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

func handleDeleteRateById(ctx echo.Context) error {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)

	err := rateController.RemoveExchangeRateById(id)
	if err != nil {
		return response(ctx, http.StatusInternalServerError, Response{Message: "Internal server error"})
	}

	return response(ctx, http.StatusOK, Response{Message: "OK"})
}

func handleNewDailyRateData(ctx echo.Context) error {
	rateData := new(controller.ExchangeRateData)
	if err := ctx.Bind(rateData); err != nil {
		return response(ctx, http.StatusBadRequest, Response{Message: "Cannot process your request"})
	}

	err := rateController.PutNewDailyExchangeRateData(*rateData)
	if err != nil {
		return response(ctx, http.StatusInternalServerError, Response{Message: "Internal server error"})
	}

	return response(ctx, http.StatusCreated, rateData)
}

func response(ctx echo.Context, statusCode int, response interface{}) error {
	return ctx.JSON(statusCode, response)
}

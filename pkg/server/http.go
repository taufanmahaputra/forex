package server

import (
	"github.com/taufanmahaputra/forex/pkg/controller"
	"net/http"

	"github.com/labstack/echo"
)

type HttpService struct {
}

func NewHttpServer() HttpService {
	return HttpService{}
}

func (s HttpService) RegisterHandler(e *echo.Echo) {
	e.GET("/", index)
	e.POST("api/v1/rate", handleNewRate)
}

type Response struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

func index(ctx echo.Context) error {
	return handleResponse(ctx, http.StatusOK, Response{
		"Foreign exchange rate API",
		"OK",
	})
}

func handleNewRate(ctx echo.Context) error {
	rate := new(controller.ExchangeRate)
	if err := ctx.Bind(rate); err != nil {
		return handleResponse(ctx, http.StatusBadRequest, Response{Message: "Invalid payload"})
	}

	err := rateController.PutNewExchangeRate(*rate)
	if err != nil {
		return handleResponse(ctx, http.StatusInternalServerError, Response{Message: "Internal server error"})
	}

	return handleResponse(ctx, http.StatusCreated, rate)
}

func handleResponse(ctx echo.Context, statusCode int, response interface{}) error {
	return ctx.JSON(statusCode, response)
}

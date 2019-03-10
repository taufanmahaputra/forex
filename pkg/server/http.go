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
	Title  string `json:"title"`
	Message string `json:"message"`
}

func index(ctx echo.Context) error {
	return handleReponse(ctx, http.StatusOK, Response{
		"Foreign exchange rate API",
		"OK",
	})
}

func handleNewRate(ctx echo.Context) error {
	rate := new(controller.ExchangeRate)
	if err := ctx.Bind(rate); err != nil {
		return handleReponse(ctx, http.StatusBadRequest, Response{Message: "Invalid payload"})
	}

	err := rateController.PutNewExchangeRate(*rate)
	if err != nil {
		return handleReponse(ctx, http.StatusInternalServerError, Response{Message: "Internal server error"})
	}

	return handleReponse(ctx, http.StatusCreated, rate)
}

func handleReponse(ctx echo.Context, statusCode int, response interface{}) error {
	return ctx.JSON(statusCode, response)
}
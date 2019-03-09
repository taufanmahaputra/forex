package server


import (
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
}

type Response struct {
	Title  string `json:"title"`
	Status string `json:"status"`
}

func index(c echo.Context) error {
	return c.JSON(http.StatusOK, Response{
		"Foreign exchange rate API",
		"OK",
	})
}

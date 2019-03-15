package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/taufanmahaputra/forex/pkg/lib/config"
	"github.com/taufanmahaputra/forex/pkg/server"
	"log"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	httpServer := server.NewHTTPServer()
	httpServer.RegisterHandler(e)

	err := server.Init()
	if err != nil {
		log.Fatal(err)
	}

	cfg := config.GetConfig()

	e.Logger.Fatal(e.Start(cfg.App.Port))
}

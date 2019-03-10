package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/taufanmahaputra/forex/pkg/server"
	"log"
)

func main()  {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	httpServer := server.NewHttpServer()
	httpServer.RegisterHandler(e)

	err := server.Init()
	if err != nil {
		log.Fatal(err)
	}

	e.Logger.Fatal(e.Start(":9999"))
}

package main

import (
	"github.com/labstack/echo"
	"github.com/taufanmahaputra/forex/pkg/server"
)

func main()  {
	e := echo.New()

	httpServer := server.NewHttpServer()
	httpServer.RegisterHandler(e)

	e.Logger.Fatal(e.Start(":9999"))
}

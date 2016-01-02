package main

import (
	"time"

	"github.com/Dataman-Cloud/seckilling/gate/src/handler"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
	"github.com/tylerb/graceful"
)

func componentInit() {
	initConfig()
}

func main() {
	// inti config and component
	componentInit()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())
	e.Use(handler.Auth)
	e.Use(handler.CrossDomain)

	// Routes
	e.Get("/hello", handler.Hello)
	e.Post("/v1/tickets", handler.Tickets)
	e.Get("/v1/over", handler.Over)
	e.Post("/v1/push", handler.Push)

	// go kafka.StartKafkaProducer()
	// Start server
	graceful.ListenAndServe(e.Server(viper.GetString("port")), 1*time.Second)

}

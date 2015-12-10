package main

import (
	"github.com/Dataman-Cloud/seckilling/queue/src/handler"
	"github.com/Dataman-Cloud/seckilling/queue/src/kafka"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

func main() {
	initConfig()
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())
	e.Use(handler.Auth)
	go kafka.StartKafkaProducer()
	// Routes
	e.Get("/hello", handler.Hello)

	e.Get("/v1/events/:id", handler.Countdown)

	// Start server
	e.Run(viper.GetString("port"))
}

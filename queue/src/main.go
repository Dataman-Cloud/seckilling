package main

import (
	"time"

	"github.com/Dataman-Cloud/seckilling/queue/src/handler"
	"github.com/Dataman-Cloud/seckilling/queue/src/kafka"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
	"github.com/tylerb/graceful"
)

func main() {
	initConfig()
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())
	e.Use(handler.Auth)

	// server favicon
	e.Favicon("public/favicon.ico")

	// server indec file
	e.Index("public/static/index.html")

	// Serve static files
	e.Static("/", "public/static")

	// Routes
	e.Get("/hello", handler.Hello)
	e.Get("/v1/events/:id", handler.Countdown)
	e.Post("/v1/tickets", handler.Tickets)

	go kafka.StartKafkaProducer()
	// Start server
	graceful.ListenAndServe(e.Server(viper.GetString("port")), 5*time.Second)

}

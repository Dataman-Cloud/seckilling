package main

import (
	"net/http"

	"github.com/Dataman-Cloud/seckilling/queue/src/kafka"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

// Handler
func hello(c *echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!\n")
}

func main() {
	initConfig()
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())
	go kafka.StartKafkaProducer()
	// Routes
	e.Get("/test", hello)

	// Start server
	e.Run(viper.GetString("port"))
}

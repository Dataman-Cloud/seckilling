package main

import (
	"log"
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
	// e.Use(handler.CrossDomain)

	// Routes
	e.Get("/hello", handler.Hello)
	e.Get("/api/v1/seckill", handler.Tickets)
	e.Get("/api/v1/reset", handler.Reset)

	// go cache.StartUpdateEventStatus()

	// go kafka.StartKafkaProducer()
	// Start server
	port := viper.GetString("port")
	log.Println("server port: ", port)
	graceful.ListenAndServe(e.Server(port), 1*time.Second)

}

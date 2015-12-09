package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Dataman-Cloud/seckilling/queue/src/kafka"
	"github.com/Dataman-Cloud/seckilling/queue/src/model"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

// Handler
func hello(c *echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!\n")
}

func countdown(c *echo.Context) error {
	seckillTime := viper.GetTime("seckillTime")
	log.Println(seckillTime)
	curTime := time.Now().UTC()
	data := model.CountdownData{
		CurTime:  curTime.Unix(),
		UnlockOn: seckillTime.Unix(),
		Locked:   curTime.Before(seckillTime),
	}

	return c.JSON(http.StatusOK, model.CommonResponse{
		Code:  0,
		Data:  data,
		Error: "",
	})
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
	e.Get("/hello", hello)

	e.Get("/v1/events/:id", countdown)

	// Start server
	e.Run(viper.GetString("port"))
}

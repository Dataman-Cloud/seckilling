package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Dataman-Cloud/seckilling/queue/src/model"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

// Handler
func Hello(c *echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!\n")
}

func Countdown(c *echo.Context) error {
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

func Auth(c *echo.Context) error {
	req := c.Request()
	if req == nil {
		return fmt.Errorf("context request is null")
	}

	cookie, err := req.Cookie(model.SkCookie)
	if err != nil {
		cookie = &http.Cookie{Name: model.SkCookie, Value: model.NewUUID(), MaxAge: 300}
		http.SetCookie(c.Response(), cookie)
	} else {
		log.Println(cookie.Value)
	}
	return nil
}

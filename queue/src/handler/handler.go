package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Dataman-Cloud/seckilling/queue/src/kafka"
	"github.com/Dataman-Cloud/seckilling/queue/src/model"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

func Auth(c *echo.Context) error {
	req := c.Request()
	if req == nil {
		return fmt.Errorf("context request is null")
	}

	cookie, err := req.Cookie(model.SkCookie)
	if err != nil {
		cookie = &http.Cookie{Name: model.SkCookie, Value: model.NewUUID(), MaxAge: 300}
		req.AddCookie(cookie)
		http.SetCookie(c.Response(), cookie)
	} else {
		log.Println(cookie.Value)
	}
	return nil
}

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

func checkCookie(c *echo.Context) string {
	cookies := c.Request().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == model.SkCookie {
			return cookie.Value
		}
	}
	return ""
}

func Tickets(c *echo.Context) error {
	cookie := checkCookie(c)
	if cookie == "" {
		return c.JSON(model.CookieCheckFailed, model.TicketData{UID: cookie, Timestamp: time.Now().UTC().Unix()})
	}

	ticket := model.TicketData{UID: cookie, Timestamp: time.Now().UTC().Unix()}
	bytes, err := json.Marshal(ticket)
	if err != nil {
		log.Printf("Marshal ticket has error: %s", err.Error())
		return c.JSON(model.PushQueueError, ticket)
	}
	kafka.ProducerMessage <- string(bytes)
	return c.JSON(http.StatusOK, ticket)

}

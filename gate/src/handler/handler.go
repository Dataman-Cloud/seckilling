package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"github.com/Dataman-Cloud/seckilling/queue/src/cache"
	"github.com/Dataman-Cloud/seckilling/queue/src/kafka"
	"github.com/Dataman-Cloud/seckilling/queue/src/model"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

var seckillTime *time.Time

func SeckillTime() *time.Time {
	if seckillTime == nil {
		time := viper.GetTime("seckillTime")
		seckillTime = &time
	}

	return seckillTime
}

func ResetSeckillTime() {
	seckillTime = nil
}

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

func CrossDomain(c *echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Request().Header.Set("Access-Control-Allow-Credentials", "true")
	c.Request().Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Request().Header.Set("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, X-Requested-By, If-Modified-Since, X-File-Name, Cache-Control, X-XSRFToken, Authorization")
	c.Request().Header.Set("Content-Type", "application/json")
	if c.Request().Method == "OPTIONS" {
		c.String(204, "")
	}
	return nil
}

// Handler
func Hello(c *echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!\n")
}

func Countdown(c *echo.Context) error {
	seckillTime := SeckillTime()
	log.Println(*seckillTime)
	curTime := time.Now().UTC()
	data := model.CountdownData{
		CurTime:  curTime.Unix(),
		UnlockOn: seckillTime.Unix(),
		Locked:   curTime.Before(*seckillTime),
	}

	return c.JSON(http.StatusOK, model.CommonResponse{
		Code:  0,
		Data:  data,
		Error: "",
	})
}

func Reset(c *echo.Context) error {
	offsetParam := c.Param("offset")
	offset := int64(10)
	if offsetSec, err := strconv.ParseFloat(offsetParam, 64); err == nil {
		offset = int64(offsetSec)
	}
	curTime := time.Now().UTC()
	newTime := curTime.Add(time.Second * time.Duration(offset))
	seckillTime = &newTime

	cmd := exec.Command("touch", viper.GetString("proxyTrigger"))
	err := cmd.Run()
	if err != nil {
		log.Fatalln("can't touch proxyTrigger")
	}
	log.Println("reset done")

	return c.Redirect(http.StatusFound, "/")
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
	curTime := time.Now().UTC()
	if curTime.Before(*seckillTime) {
		return c.Redirect(http.StatusFound, "/")
	}

	cookie := checkCookie(c)
	if cookie == "" {
		return c.JSON(model.CookieCheckFailed, model.TicketData{UID: cookie, Timestamp: time.Now().UTC().Unix()})
	}

	ticket := model.TicketData{UID: cookie, Timestamp: time.Now().UTC().Unix()}
	model.UIDMap[cookie] = false
	bytes, err := json.Marshal(ticket)
	if err != nil {
		log.Printf("Marshal ticket has error: %s", err.Error())
		return c.JSON(model.PushQueueError, ticket)
	}
	kafka.ProducerMessage <- string(bytes)
	err = cache.WriteHashToRedis(cookie, "Status", "0", -1)
	if err != nil {
		log.Printf("write ticket to redis has error %s", err.Error())
		return c.JSON(model.PushQueueError, ticket)
	}

	return c.JSON(http.StatusOK, model.CommonResponse{
		Code:  0,
		Data:  ticket,
		Error: "",
	})

}

func Over(c *echo.Context) error {
	return c.JSON(http.StatusOK, model.CommonResponse{
		Code:  99,
		Data:  "Game Over",
		Error: "",
	})
}

func Push(c *echo.Context) error {
	if (rand.Intn(10)) > 5 {
		time.Sleep(time.Second * 1)
	} else {
		time.Sleep(time.Millisecond * 50)
	}
	req := c.Request()
	if req == nil {
		return fmt.Errorf("context request is null")
	}
	cookie, err := req.Cookie(model.SkCookie)
	if err != nil {
		return c.JSON(http.StatusOK, model.CommonResponse{
			Code:  99,
			Data:  "tickets error",
			Error: "",
		})
	}

	model.UIDMap[cookie.Value] = true
	return c.JSON(http.StatusOK, model.CommonResponse{
		Code:  0,
		Data:  "Game Over",
		Error: "",
	})
}

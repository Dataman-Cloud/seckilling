package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Dataman-Cloud/seckilling/gate/src/cache"
	"github.com/Dataman-Cloud/seckilling/gate/src/model"
	"github.com/labstack/echo"
)

var CountKey = "Stock:1"

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
		return c.JSON(model.CookieCheckFailed, model.OrderInfo{Timestamp: time.Now().UTC().Unix()})
	}

	// Get user info by cookie(UUID)
	// if cookie is not exit in redis return error
	// if status is null or status not 1 return StatusNotOne/StatusNull
	// if phone is null return UserPhoneNumNull
	// if event is null or event is not match current event return EventNull/EventNotMatch
	user, code := cache.GetOrderInfo(cookie)
	if code != 0 {
		return c.JSON(code, model.OrderInfo{Timestamp: time.Now().UTC().Unix()})
	}

	if user == nil {
		return c.JSON(model.UnknownError, model.OrderInfo{Timestamp: time.Now().UTC().Unix()})
	}

	// check phone number and event id make sure one phone number only have once chance in one activity
	repeat, err := cache.CheckPhoneNum(user.Phone)
	if err != nil {
		log.Println("check user phone number hs error: ", err)
	}

	if !repeat {
		return c.JSON(model.PhoneRepaet, user)
	}

	code = ProduceOrder(user)
	if code == 0 {
		go SaveOrder(user)
		return c.JSON(http.StatusOK, user)
	}

	return c.JSON(code, user)
}

func ProduceOrder(user *model.OrderInfo) int {
	// TODO use MultiBulk
	stock, err := cache.Decr(CountKey)
	if err != nil {
		log.Println("get stock has error err", err)
		return model.RedisError
	} else if stock <= 0 {
		log.Println("short of stock!! stock is ", stock)
		return model.ShortageStock
	}
	serialNum, err := cache.GetSerialNum(user.EventId, stock)
	if err != nil {
		log.Println("Get serial number has error: ", err)
		return model.RedisError
	}

	user.Index = stock
	user.SerialNum = serialNum
	return 0
}

func SaveOrder(user *model.OrderInfo) {
	key := fmt.Sprintf(model.OrderKey, user.Eid, user.SerialNum)
	cache.WriteStructToRedis(user, key)
}

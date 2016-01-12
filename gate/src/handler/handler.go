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

// func checkCookie(c *echo.Context) string {
// 	cookies := c.Request().Cookies()
// 	log.Println("cookies: ", cookies)
// 	for _, cookie := range cookies {
// 		log.Println(cookie.Name)
// 		if cookie.Name == model.SkCookie {
// 			return cookie.Value
// 		}
// 	}
// 	ck, err := c.Request().Cookie(model.SkCookie)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	log.Println(ck)

// 	log.Println(c.Get(model.SkCookie))

// 	return ""
// }

func Tickets(c *echo.Context) error {
	// cookie := checkCookie(c)
	// cookie := c.Param(model.SkCookie)
	cookie := c.Query(model.SkCookie)

	if cookie == "" {
		log.Println("Error!! Cookie is null")
		return c.JSON(model.CookieCheckFailed, model.OrderInfo{Timestamp: time.Now().UTC().Unix()})
	}

	phoneNum := c.Query("phone")
	if phoneNum == "" {
		log.Println("Error!! Phone is null")
		return c.JSON(model.UserPhoneNumNull, model.OrderInfo{Timestamp: time.Now().UTC().Unix(), UID: cookie})
	}

	eid := c.Query("id")
	if eid == "" {
		log.Println("Error!! Id is null")
		return c.JSON(model.UserPhoneNumNull, model.OrderInfo{Timestamp: time.Now().UTC().Unix(), UID: cookie})
	}

	model.CurrentEventId = eid

	// Get user info by cookie(UUID)
	// if cookie is not exit in redis return error
	// if status is null or status not 1 return StatusNotOne/StatusNull
	// if phone is null return UserPhoneNumNull
	// if event is null or event is not match current event return EventNull/EventNotMatch
	ckey := fmt.Sprintf(model.CookHashKey, cookie, eid)
	err := cache.CheckStatus(ckey)
	if err != nil {
		log.Println("Error!! Status of %s is invalid", ckey)
		return c.JSON(model.InvalidStatus, model.OrderInfo{Timestamp: time.Now().UTC().Unix(), UID: cookie, EventId: eid})
	}

	// check phone number and event id make sure one phone number only have once chance in one activity
	repeat, err := cache.CheckPhoneNum(phoneNum)
	if err != nil {
		log.Println("check user phone number hs error: ", err)
	}

	order := &model.OrderInfo{
		Timestamp: time.Now().UTC().Unix(),
		UID:       cookie,
		EventId:   eid,
		Phone:     phoneNum,
	}

	if !repeat {
		return c.JSON(model.PhoneRepaet, order)
	}

	code := ProduceOrder(order)
	if code == 0 {
		go SaveOrder(order)
		return c.JSON(http.StatusOK, order)
	}

	return c.JSON(code, order)
}

func ProduceOrder(user *model.OrderInfo) int {
	// TODO use MultiBulk
	serialNum, stock, err := cache.GetSerialNum()
	if err != nil {
		log.Println("Get serial number has error: ", err)
		return model.RedisError
	}

	user.Index = stock
	user.SerialNum = serialNum
	return 0
}

func SaveOrder(user *model.OrderInfo) {
	key := fmt.Sprintf(model.OrderKey, user.EventId, user.SerialNum)
	cache.WriteStructToRedis(user, key)
}

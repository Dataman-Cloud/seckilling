package cache

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Dataman-Cloud/seckilling/gate/src/model"
	redis "github.com/garyburd/redigo/redis"
	"github.com/spf13/viper"
)

var (
	pool         *redis.Pool
	CtEventId    string
	ValidityTime int64
)

func Open() redis.Conn {
	if pool != nil {
		return pool.Get()
	}

	mutex := &sync.Mutex{}
	mutex.Lock()
	InitCache()
	defer mutex.Unlock()

	return pool.Get()
}

func initConn() (redis.Conn, error) {
	cacheHost := viper.GetString("cache.host")
	cachePort := viper.GetInt("cache.port")
	addr := fmt.Sprintf("%s:%d", cacheHost, cachePort)
	c, err := redis.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return c, err
}

func InitCache() {
	poolSize := viper.GetInt("cache.poolSize")
	pool = redis.NewPool(initConn, poolSize)
	conn := Open()
	defer conn.Close()
	pong, err := conn.Do("ping")
	if err != nil {
		log.Panicln("can't connect cache server has error", err)
	}
	log.Println("reach cache server ", pong)
}

func DestroyCache() {
	log.Println("destroying Cache")
	if pool != nil {
		pool.Close()
		log.Println("cache was closed")
	}
}

func Decr(key string) (int64, error) {
	conn := Open()
	defer conn.Close()

	return redis.Int64(conn.Do("DECR", key))
}

func WriteHashToRedis(key, field, value string, timeout int) error {
	conn := Open()
	defer conn.Close()
	var err error
	log.Printf("redis HSET: %s, field: %s, value: %s", key, field, value)
	if _, err = conn.Do("HSET", key, field, value); err != nil {
		return err
	}

	if timeout != -1 {
		_, err = conn.Do("EXPIRE", key, timeout)
		return err
	}
	return nil
}

func GetUserInfo(cookie string) (*model.UserInfo, int) {
	conn := Open()
	defer conn.Close()

	status, err := redis.String(conn.Do("HGET", cookie, "status"))
	if err == redis.ErrNil {
		return nil, model.StatusNull
	} else if err != nil {
		return nil, model.GetStatusFailed
	}

	if status != "1" {
		return nil, model.StatusNotOne
	}

	phoneNum, err := redis.String(conn.Do("HGET", cookie, "phone"))
	if err == redis.ErrNil {
		return nil, model.UserPhoneNumNull
	} else if err != nil {
		return nil, model.GetPhoneNumFailed
	}

	eventId, err := redis.String(conn.Do("HGET", cookie, "event"))
	if err == redis.ErrNil {
		return nil, model.EventNull
	} else if err != nil {
		return nil, model.GetEventFailed
	}

	ctEvent, err := GetCurrentEventId()
	if err != nil {
		return nil, model.GetCtEventFailed
	}

	if ctEvent != eventId {
		return nil, model.EventNotMatch
	}

	return &model.UserInfo{
		UID:       cookie,
		Phone:     phoneNum,
		EventId:   eventId,
		Timestamp: time.Now().UTC().Unix(),
	}, 0

}

func GetCurrentEventId() (string, error) {

	if time.Now().UTC().Unix() < ValidityTime && CtEventId != "" {
		return CtEventId, nil
	}

	conn := Open()
	defer conn.Close()

	eventId, err := redis.String(conn.Do("HGET", "CurrentEvent", "ID"))
	if err != nil {
		log.Println("get current event id has error: ", err)
		return "", err
	}

	CtEventId = eventId
	start, err := redis.Int64(conn.Do("HGET", "CurrentEvent", "start"))
	if err != nil {
		log.Println("get current event start time has error: ", err)
		return "", err
	}

	duration, err := redis.Int64(conn.Do("HGET", "CurrentEvent", "duration"))
	if err != nil {
		log.Println("get current event duration time has error: ", err)
		return "", err
	}

	ValidityTime = start + duration
	return CtEventId, err
}

func CheckPhoneNum(phone string) (bool, error) {
	conn := Open()
	defer conn.Close()

	eventId, err := redis.String(conn.Do("HGET", phone, "event"))
	if err == redis.ErrNil {
		return true, nil
	} else if err != nil {
		return false, err
	}

	// TODO check once when new event begin
	ctEventId, err := GetCurrentEventId()
	if err != nil {
		return false, err
	}

	return ctEventId == eventId, nil
}

func GetPhoneNum(cookie string) (string, error) {
	conn := Open()
	defer conn.Close()

	phoneNum, err := redis.String(conn.Do("HGET", cookie, "phone"))
	if err == redis.ErrNil {
		return "", nil
	}

	return phoneNum, err
}

func GetSerialNum(key, field string) (string, error) {
	conn := Open()
	defer conn.Close()
	return redis.String(conn.Do("HGET", key, field))
}

func UpdateCurEventId(curEid string) error {
	conn := Open()
	defer conn.Close()

	_, err := conn.Do("SET", model.CurrentEventKey, curEid)
	return err
}

package cache

import (
	"fmt"
	"log"
	"sync"

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

func CheckStatus(key string) error {
	conn := Open()
	defer conn.Close()

	status, err := redis.String(conn.Do("HGET", key, "status"))
	if err == redis.ErrNil {
		return fmt.Errorf("empty status of %s", key)
	} else if err != nil {
		return fmt.Errorf("get status of %s has error: ", key, err.Error())
	}

	if status != "1" {
		return fmt.Errorf("invalid status of %s", key)
	}

	return nil
}

func GetCurrentEventId() (string, error) {
	// TODO check again
	return model.CurrentEventId, nil
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

func GetSerialNum() (string, int64, error) {
	index, err := GetSeriaIndex()
	if err != nil {
		log.Println("GetSeriaIndex has error: ", err)
		return "", -1, err
	}
	log.Println("indec: ", index)
	conn := Open()
	defer conn.Close()
	eid, _ := GetCurrentEventId()
	eidKey := fmt.Sprintf(model.EventIdKey, eid)
	log.Println(eidKey)
	indexKey := fmt.Sprintf(model.WorkOffIndexKey, eid)
	log.Println(indexKey)
	conn.Send("MULTI")
	conn.Send("ZRANGE", eidKey, index, index)
	conn.Send("INCR", indexKey)
	r, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		log.Println("get seria num has error: ", err)
		return "", index, err
	}

	if r[0] == redis.ErrNil {
		return "", index, model.ShortageStockError
	}

	var seriaNum string
	if slice, ok := r[0].([]interface{}); ok {
		if bytes, ok := slice[0].([]byte); ok {
			seriaNum = string(bytes)
			return seriaNum, index, nil
		}
	}

	return seriaNum, index, fmt.Errorf("unknown result")
}

func UpdateCurEventId(curEid string) error {
	conn := Open()
	defer conn.Close()

	_, err := conn.Do("SET", model.CurrentEventKey, curEid)
	return err
}

func GetSeriaIndex() (int64, error) {
	conn := Open()
	defer conn.Close()
	id, _ := GetCurrentEventId()
	indexKey := fmt.Sprintf(model.WorkOffIndexKey, id)
	index, err := redis.Int64(conn.Do("GET", indexKey))
	if err == redis.ErrNil {
		return 0, nil
	} else if err != nil {
		return -1, err
	}

	return index, nil
}

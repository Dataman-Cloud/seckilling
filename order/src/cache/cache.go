package cache

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	redis "github.com/garyburd/redigo/redis"
	"github.com/spf13/viper"
)

var pool *redis.Pool

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

func Decr(key string) (int64, error) {
	conn := Open()
	defer conn.Close()
	result, err := conn.Do("DECR", key)
	if err != nil {
		return -1, err
	}

	sku, err := strconv.ParseInt(fmt.Sprint(result), 10, 64)
	if err != nil {
		return -1, err
	}

	return sku, nil
}

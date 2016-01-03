package cache

import (
	"reflect"

	redis "github.com/garyburd/redigo/redis"
)

func ReadStructFromRedis(v interface{}, key string) {
	conn := Open()
	defer conn.Close()
	val := reflect.ValueOf(v).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		if !valueField.CanSet() {
			continue
		}
		switch typeField.Type.Kind() {
		case reflect.String:
			str, err := redis.String(conn.Do("HGET", key, tag.Get("json")))
			if err != nil {
				continue
			}

			valueField.SetString(str)
		case reflect.Int:
			integer, err := redis.Int64(conn.Do("HGET", key, tag.Get("json")))
			if err != nil {
				continue
			}

			valueField.SetInt(integer)
		}
	}
}

// func WriteStructToRedis(v interface{}, key string) error {
// 	conn := Open()
// 	defer conn.Close()
// 	val := reflect.ValueOf(v).Elem()
// 	for i := 0; i < val.NumField(); i++ {
// 		valueField := val.Field(i)
// 		typeField := val.Type().Field(i)
// 		tag := typeField.Tag
// 		switch typeField.Type.Kind() {
// 		case reflect.String:
// 			str, err := redis.String(conn.Do("HSET", key, tag.Get("json"), valueField.Interface()))
// 			if err != nil {
// 				return err
// 			}
// 		case reflect.Int:
// 			integer, err := redis.Int(conn.Do("HSET", key, tag.Get("json"), valueField.Interface()))
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
// }

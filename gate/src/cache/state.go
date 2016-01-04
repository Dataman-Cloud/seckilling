package cache

import (
	"fmt"
	"log"
	"time"

	"github.com/Dataman-Cloud/seckilling/gate/src/model"
	redis "github.com/garyburd/redigo/redis"
)

const (
	eventUpdateInterval = time.Minute * 1
)

func LoadEventList() ([]string, error) {
	conn := Open()
	defer conn.Close()

	eventNum, err := redis.Int(conn.Do("LLEN", model.EventListKey))
	if err != nil {
		log.Println("Get event list length has error: ", err)
		return nil, err
	}

	var eventIdList = make([]string, eventNum)

	for i := 0; i < eventNum; i++ {
		eventId, err := redis.String(conn.Do("LINDEX", model.EventListKey, i))
		if err != nil {
			log.Printf("Get event list %s index %d has error %f", model.EventListKey, i, err.Error())
			continue
		}

		eventIdList = append(eventIdList, eventId)
	}

	return eventIdList, nil

}

func LoadEventData() ([]*model.EventInfo, error) {
	eventIds, err := LoadEventList()
	if err != nil {
		return nil, err
	}

	var eventInfoList = make([]*model.EventInfo, len(eventIds))
	for _, eventId := range eventIds {
		eventInfoKey := fmt.Sprintf(model.EventInfoKey, eventId)
		evenInfo := &model.EventInfo{}
		ReadStructFromRedis(evenInfo, eventInfoKey)
		eventInfoList = append(eventInfoList, evenInfo)
	}

	return eventInfoList, nil
}

func UpdateEvent() error {
	eventInfos, err := LoadEventData()
	if err != nil {
		log.Println("update event ststus failed error: ", err)
		return err
	}

	// judgment event status
	// 1 is under way
	// 2 is has not started
	// 3 is over
	for _, event := range eventInfos {
		timestamp := time.Now().UTC().Unix()
		if event == nil {
			continue
		}
		if event.EffectOn <= timestamp && timestamp <= event.EffectOn+event.Duration {
			err = WriteHashToRedis(event.Id, "status", "1", -1)
			if err != nil {
				log.Println("write event status to redis has error: ", err)
				continue
			}
			log.Println("update under way event succes id ", event.Id)
			model.CurrentEventId = event.Id
			err = UpdateCurEventId(model.CurrentEventId)
			if err != nil {
				log.Println("update current id failed id ", model.CurrentEventId)
			}
			log.Println("update current id success id ", model.CurrentEventId)
		} else if event.EffectOn > timestamp {
			err = WriteHashToRedis(event.Id, "status", "2", -1)
			if err != nil {
				log.Println("write event status to redis has error: ", err)
				continue
			}
		} else if timestamp > event.EffectOn+event.Duration {
			err = WriteHashToRedis(event.Id, "status", "3", -1)
			if err != nil {
				log.Println("write event status to redis has error: ", err)
				continue
			}
		}
	}
	return nil
}

func StartUpdateEventStatus() {
	updateTicker := time.NewTicker(eventUpdateInterval)
	defer updateTicker.Stop()

	for {
		select {
		case <-updateTicker.C:
			UpdateEvent()
		}
	}
}

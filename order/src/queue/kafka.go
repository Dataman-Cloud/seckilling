package queue

import (
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/Dataman-Cloud/seckilling/order/src/cache"
	"github.com/Dataman-Cloud/seckilling/order/src/db"
	"github.com/Dataman-Cloud/seckilling/order/src/model"
	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
	csgroup "github.com/wvanbergen/kafka/consumergroup"
	kazoo "github.com/wvanbergen/kazoo-go"
)

var (
	zkURI   string
	group   string
	topic   string
	zkNodes []string
	config  *csgroup.Config
)

//InitClient initialize client
func InitClient() {
	zkURI = viper.GetString("kafka.zkURI")
	log.Printf("zk nodes: %s \n", zkURI)

	group = viper.GetString("kafka.group")
	log.Println("group", group)

	topic = viper.GetString("kafka.topic")
	log.Printf("kafka topic name %s \n", topic)

	//sarama.Logger = log.New(os.Stdout, "[Sarama] ", log.LstdFlags)

	config = csgroup.NewConfig()
	config.Offsets.Initial = sarama.OffsetNewest
	config.Offsets.ProcessingTimeout = 10 * time.Second

	zkNodes, config.Zookeeper.Chroot = kazoo.ParseConnectionString(zkURI)
}

//Start start consumers
func Start() {
	workers := viper.GetInt("workers")
	log.Printf("workers: %d\n", workers)
	startConsumer()
}

func startConsumer() {
	log.Println("starting consumer...")
	csGroup, err := csgroup.JoinConsumerGroup(group, []string{topic}, zkNodes, config)
	if err != nil {
		log.Panicln(err)
	}

	defer func() {
		if err = csGroup.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	go handleErrors(csGroup)
	handleMessages(csGroup)
}

func handleErrors(csGroup *csgroup.ConsumerGroup) {
	for err := range csGroup.Errors() {
		log.Println(err)
	}
}

func handleMessages(csGroup *csgroup.ConsumerGroup) {
	consumed := 0
	for message := range csGroup.Messages() {
		log.Println("received message#", consumed)
		logMessage(message)
		checkUerTicket(message)

		consumed++
		csGroup.CommitUpto(message)
	}
}

func logMessage(message *sarama.ConsumerMessage) {
	log.Println("================================")
	log.Println("key:", message.Key)
	log.Println("value:", string(message.Value))
	log.Println("topic:", message.Topic)
	log.Println("partition:", message.Partition)
	log.Println("offset:", message.Offset)
	log.Println("--------------------------------")
}

func checkUerTicket(message *sarama.ConsumerMessage) {
	eid, err := db.GetEId()
	if err != err {
		log.Println("Get EId has error do nothing")
		return
	}

	seq, err := updateStockInfo(eid)
	if err != nil {
		log.Println("Get merchandise inventory has error do nothing")
		return
	}

	order := model.Order{
		EId:    eid,
		UId:    string(message.Key),
		Seq:    seq,
		Status: 1,
		Ext:    "",
		Create: time.Now().UTC(),
	}

	//TODO if has error or panic must rollback
	err = writeOrderToDb(order)
	if err != nil {
		log.Println("Insert order to db has error: ", err)
		return
	}

	err = writeOrderToCache(order)
	if err != nil {
		log.Println("Insert order to cache has error: ", err)
		return
	}

}

func writeOrderToCache(order model.Order) error {
	val := reflect.ValueOf(order).Elem()
	conn := cache.Open()
	defer conn.Close()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		conn.Send("HSET", order.UId, typeField.Name, valueField.Interface())
	}

	_, err := conn.Do("EXEC")
	return err
}

func writeOrderToDb(order model.Order) error {
	return db.InsertOrder(order)
}

func updateStockInfo(eid int64) (int64, error) {
	eidStr := strconv.FormatInt(eid, 10)
	countId := "dataman-" + eidStr
	seq, err := cache.Decr(countId)
	if err != nil {
		log.Println("get merchandise inventory has error: ", err)
		return -1, err
	}
	return seq, nil
}

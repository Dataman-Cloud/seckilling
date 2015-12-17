package queue

import (
	"log"
	"time"

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
	updateStockInfo()
	//stock := 1
	value := string(message.Value)
	log.Println(value)
	writeTicketToCache()
	writeTicketToDb()
}

func writeTicketToCache() error {
	return nil
}

func writeTicketToDb() error {
	return nil
}

func updateStockInfo() {

}

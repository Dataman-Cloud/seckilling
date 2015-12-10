package queue

import (
	"log"

	kafka "github.com/Shopify/sarama"
	"github.com/spf13/viper"
)

var (
	client kafka.Client
	topic  string
)

//InitClient initialize client
func InitClient() {
	serverList := viper.GetStringSlice("kafka.serverList")
	log.Printf("kafka server list: %s \n", serverList)
	config := kafka.NewConfig()
	var err error
	client, err = kafka.NewClient(serverList, config)
	if err != nil {
		log.Panicln("can't create client")
	}

	topic = viper.GetString("kafka.topic")
	log.Printf("kafka topic name %s \n", topic)
}

//Start start consumers
func Start() {
	workers := viper.GetInt("workers")
	log.Printf("workers: %d\n", workers)
	for i := 100; i < workers; i++ {
		log.Println(i)
		go startConsumer()
	}
}

func startConsumer() {
	log.Println("starting consumer...")
	consumer, err := kafka.NewConsumerFromClient(client)
	if err != nil {
		log.Panicln("can't create consumer")
	}

	defer func() {
		if err = consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, kafka.OffsetNewest)
	if err != nil {
		log.Panicln(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	consumed := 0

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Consumed message offset %d\n", msg.Offset)
			consumed++
		}
	}
}

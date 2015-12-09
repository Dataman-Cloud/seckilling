package kafka

import (
	"log"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
)

var (
	ProducerMessage chan string
)

func init() {
	ProducerMessage = make(chan string, 5)
}

func StartKafkaProducer() {
	kafkaServerList := viper.GetStringSlice("kafka.serverList")
	log.Printf("kafka server list: %s \n", kafkaServerList)

	topic := viper.GetString("kafka.topic")
	log.Printf("kafka topic name %s \n", topic)

	producer, err := sarama.NewAsyncProducer(kafkaServerList, nil)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = producer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

ProducerLoop:
	for {
		select {
		case message := <-ProducerMessage:
			producer.Input() <- &sarama.ProducerMessage{Topic: topic, Key: nil, Value: sarama.StringEncoder(message)}
		case err = <-producer.Errors():
			log.Println("Failed to produce message", err)
		case <-signals:
			break ProducerLoop
		}
	}

}

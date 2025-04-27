package main

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers": "fc2-gokafka-kafka-1:9092",
		"client.id":         "goapp-consumer",
		"group.id":       "goapp-group",
		"auto.offset.reset": "earliest",
	}

	c, err := kafka.NewConsumer(configMap)
	if err != nil {
		log.Println("Failed to create consumer:", err.Error())
		panic(err)
	}

	topic := []string{"novo-teste"}

	c.SubscribeTopics(topic, nil)

	for {
		msg, err := c.ReadMessage(-1)

		if err == nil {
			log.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		}
	}

}

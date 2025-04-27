package main

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
	deliveryChannel := make(chan kafka.Event)
	producer, err := NewKafkaProducer()

	if err != nil {
		log.Println("Failed to create producer:", err.Error())
		return
	}

	for i := 0; i < 10; i++ {
		Publish("transferiu", "novo-teste", producer, []byte("transferencia"), deliveryChannel)
	}

	go DeliveryReport(deliveryChannel) //async

	producer.Flush(1000)
	// Wait for all messages to be delivered
}

// Cria um novo producer Kafka
func NewKafkaProducer() (*kafka.Producer, error) {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers":   "fc2-gokafka-kafka-1:9092",
		"delivery.timeout.ms": "0",
		"acks":                "all",
		"enable.idempotence":  "true", // por padrão é false
	}

	producer, err := kafka.NewProducer(configMap)

	if err != nil {
		log.Println("Failed to create producer:", err.Error())
		return nil, err
	}

	return producer, nil
}

// Publica uma mensagem no Kafka
func Publish(
	message string,
	topic string,
	producer *kafka.Producer,
	key []byte,
	deliveryChannel chan kafka.Event) error {

	msg := &kafka.Message{
		Value: []byte(message),
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key: key,
	}

	err := producer.Produce(msg, deliveryChannel)

	if err != nil {
		log.Println("Failed to produce message:", err.Error())
		return err
	}
	return nil
}

func DeliveryReport(deliveryChannel chan kafka.Event) {
	for e := range deliveryChannel {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				log.Printf("Failed to deliver message: %v\n", ev.TopicPartition.Error)
			} else {
				// poderiamos anotar no banco de dados que a mensagem foi entregue
				// ou fazer qualquer outra coisa
				log.Printf("Message delivered to %v [%d] at offset %v\n",
					*ev.TopicPartition.Topic,
					ev.TopicPartition.Partition,
					ev.TopicPartition.Offset)
			}
		default:
			log.Printf("Ignored event: %s\n", ev)

		}
	}
}

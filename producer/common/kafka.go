package common

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
)

var producer *kafka.Producer

func GetProducer() (*kafka.Producer, error) {
	if producer != nil {
		return producer, nil
	}

	e := GetEnv()
	logrus.Info(&kafka.ConfigMap{
		"bootstrap.servers":  e.KafkaBootstrapServers,
		"message.timeout.ms": e.KafkaMessageTimeout,
		"retries":            e.KafkaRetries,
		"retry.backoff.ms":   e.KafkaRetryBackoff,
	})
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":  e.KafkaBootstrapServers,
		"message.timeout.ms": e.KafkaMessageTimeout,
		"retries":            e.KafkaRetries,
		"retry.backoff.ms":   e.KafkaRetryBackoff,
	})
	if err != nil {
		return nil, err
	}

	producer = p

	return producer, err
}

// send a message
func ProduceEvent(message []byte) error {
	p, _ := GetProducer()

	kafkaMessage := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &GetEnv().KafkaTopic,
			Partition: kafka.PartitionAny,
		},
		Value: message,
	}
	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)
	err := p.Produce(kafkaMessage, deliveryChan)
	if err != nil {
		return err
	}

	// pick delivery report
	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}

	return nil
}

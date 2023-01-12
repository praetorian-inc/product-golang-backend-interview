package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	dto "github.com/praetorian-inc/product-golang-backend-interview/internal"
)

// Producer wraps a kafka producer.
type Producer struct{ *kafka.Producer }

// Ensure Producer satisfies the required interface.
var _ dto.Producer = (*Producer)(nil)

// SetupKafkaProducer returns a producer.
func SetupKafkaProducer() (*Producer, error) {
	// Default port is 9092
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	return &Producer{producer}, err
}

// ProduceMsg calls Produce with the topic and message.
func (p *Producer) ProduceMsg(topic string, message []byte) error {
	delivery_chan := make(chan kafka.Event, 10000)
	return p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message},
		delivery_chan,
	)
}

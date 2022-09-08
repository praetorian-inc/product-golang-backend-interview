package handlers

import "github.com/confluentinc/confluent-kafka-go/kafka"

func (s *Server) Produce(topic string, message []byte) error {
	delivery_chan := make(chan kafka.Event, 10000)
	return s.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message},
		delivery_chan,
	)
}

package dto

type KafkaMessage struct {
	Type    string
	Payload map[string]interface{}
}

package dto

type (
	// IngestDto is the ingest data structure.
	IngestDto struct {
		Domain string
		Id     uint32
	}

	// RootDomainDto is the root domain data structure.
	RootDomainDto struct {
		Id     uint32
		Root   string
		Domain string
		Status string
		Owner  string
	}

	// SubdomainDto is the subdomain data structure.
	SubdomainDto struct {
		Id     uint32
		Root   string
		Source string
	}

	// KafkaMessage defines a simplified kafka message.
	KafkaMessage struct {
		Type    string
		Payload interface{}
	}

	// EnumerateFn defines a function that is used to enumerate a subdomain.
	EnumerateFn func(string) []string

	// MessageHandlerFn defines a function that takes acts on a message.
	MessageHandlerFn func(message KafkaMessage) error

	// Producer defines the interface for producing kafka messages.
	Producer interface {
		ProduceMsg(topic string, message []byte) error
	}

	// Poller defines the interface for consuming kafka messages.
	Poller interface {
		PollKafka(mh MessageHandlerFn)
	}
)

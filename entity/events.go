package entity

type Events struct {
	Enabled           bool           `json:"enabled"`
	Transport         EventTransport `json:"transport"`
	AllEntities       bool           `json:"all-entities"`
	EntityIdentifiers []string       `json:"entity-identifiers"`
}

type EventTransport string

const (
	InvalidEventTransport EventTransport = "invalid"
	KafkaEventTransport   EventTransport = "kafka"
)

type TransportConfig struct {
}

type KafkaTransportConfig struct {
	Version string   `json:"version"`
	Brokers []string `json:"brokers"`
	Topics  []string `json:"topics"`
}

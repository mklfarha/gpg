package entity

func (p Project) KafkaEnabled() bool {
	if !p.Events.Enabled {
		return false
	}

	if p.Events.Transport == KafkaEventTransport && p.Events.TransportConfig.Kafka != nil {
		return true
	}

	return false
}

func (p Project) KafkaConfig() *KafkaTransportConfig {
	if p.KafkaEnabled() {
		return p.Events.TransportConfig.Kafka
	}
	return nil
}

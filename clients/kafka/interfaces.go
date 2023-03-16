package kafka

type PublisherConfig interface {
	GetBrokers() []string
	GetTopic() string
}

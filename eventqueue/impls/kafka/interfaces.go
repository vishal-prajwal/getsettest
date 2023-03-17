package kafka

type PublisherConfig interface {
	GetBrokers() []string
	GetTopic() string
}

type ConsumerConfig interface {
	GetBrokers() []string
	GetTopic() string
	GetGroupId() string
}

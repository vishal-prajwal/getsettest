package kafka

type PublisherConfig interface {
	GetBrokers() []string
	GetTopic() string
	GetPublisherCount() int
	GetAsyncQueueSize() int
}

type ConsumerConfig interface {
	GetBrokers() []string
	GetTopic() string
	GetGroupId() string
}

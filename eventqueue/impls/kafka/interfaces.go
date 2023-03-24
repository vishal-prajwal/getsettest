package kafka

type PublisherConfig interface {
	GetBrokers() []string
	GetTopic() string
	GetPublisherCount() int
	GetAsyncQueueSize() int
	GetMaxRetries() int
}

type ConsumerConfig interface {
	GetBrokers() []string
	GetTopic() string
	GetGroupId() string
	GetAsyncQueueSize() int
	GetBatchSize() int
	GetMaxWaitSeconds() int
}

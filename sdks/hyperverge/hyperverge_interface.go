package hyperverge

type Hyperverge interface {
	ReadDocument(documentType string, hypervergeRequest HypervergeRequest) (*HypervergeResponse, error)
}

type HypervergeConfig interface {
	GetHypervergeAppID() string
	GetHypervergeAppKey() string
	GetHypervergeEndpoint() string
	GetHypervergeTransactionId() string
}

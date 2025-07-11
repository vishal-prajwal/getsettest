package apilogger

import "github.com/pkg/errors"

const (
	FAILED_TO_CREATE_KAFKA_PUBLISHER = "failed to create Kafka publisher for API usage logger"
)

var NIL_CONFIG = errors.New("config cannot be nil")

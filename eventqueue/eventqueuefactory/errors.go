package eventqueuefactory

import "fmt"

var (
	ErrInvalidPublisherName = fmt.Errorf("invalid event pulisher name")
	ErrInvalidConsumerName  = fmt.Errorf("invalid event consumer name")
)

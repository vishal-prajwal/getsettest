package eventqueue

import (
	"context"
)

type Publisher interface {
	Publish(ctx context.Context, key any, msg any) error
	PublishAsync(ctx context.Context, key any, msg any)
	GetAsyncPublishResponseChan() *chan error
	Close()
}

type Consumer interface {
}

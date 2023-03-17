package eventqueue

import (
	"context"
)

type Publisher interface {
	Publish(ctx context.Context, key any, msg any) error
	Close()
}

type Consumer interface {
}

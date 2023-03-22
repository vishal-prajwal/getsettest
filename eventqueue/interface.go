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

type Message struct {
	Key   []byte
	Value []byte
}
type Consumer interface {
	// it will return a message and also commits
	ReadMessage(ctx context.Context) (*Message, error)

	// it will fetch and return a batch without commiting it , and commits the previously fetched batch
	ReadBatch(ctx context.Context) ([]Message, error)

	// it will close the reader
	Close() error
}

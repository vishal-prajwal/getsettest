package eventqueue

import (
	"context"
)

type Publisher interface {
	Publish(ctx context.Context, key any, msg any) error
	PublishMany(ctx context.Context, messages []Message) error
	PublishAsync(ctx context.Context, key any, msg any)
	GetAsyncPublishResponseChan() chan error
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

	// it will return a message without commiting it
	ReadMessageWithUnCommit(ctx context.Context) (*Message, error)

	// it will commit the message
	Commit(ctx context.Context) error

	// it will close the reader
	Close() error
}

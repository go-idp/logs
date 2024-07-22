package pubsub

import (
	"context"
	"io"

	"github.com/go-zoox/core-utils/safe"
)

type PubSub interface {
	Create(ctx context.Context, topic string) error
	Destroy(ctx context.Context, topic string) error
	//
	Publish(ctx context.Context, topic string, message string) error
	Subscribe(ctx context.Context, topic string, handler func(message string)) error
}

type Message struct {
	ID        int64  `json:"id"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

type File struct {
	Path string `json:"path"`
	//
	Writer io.WriteCloser
	//
	Clean func() error
	//
	GetReader func() (io.ReadCloser, error)
}

var topicStore = safe.NewMap[string, bool]()

var countsStore = safe.NewMap[string, int64]()

var messagesStore = safe.NewMap[string, *safe.Queue[*Message]]()

var subscribersStore = safe.NewMap[string, *safe.Map[string, func(message *Message)]]()

var filesStore = safe.NewMap[string, *File]()

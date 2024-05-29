package pubsub

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/go-zoox/core-utils/safe"
	"github.com/go-zoox/fs"
)

func Open(ctx context.Context, topic string) error {
	// check if topic exists
	if ok := topicStore.Has(topic); ok {
		return fmt.Errorf("topic %s already exists", topic)
	}

	// create topic store
	if err := topicStore.Set(topic, true); err != nil {
		return fmt.Errorf("failed to create topic store: %s", err)
	}

	// create messages store
	messages := safe.NewQueue[*Message](func(qc *safe.QueueConfig) {
		qc.Capacity = DefaultMessageCapacityForEachTopic
	})
	messages.Enqueue(&Message{
		ID:        1,
		Content:   "waiting for logs ...",
		Timestamp: time.Now().UnixMilli(),
	})
	if err := messagesStore.Set(topic, messages); err != nil {
		return fmt.Errorf("failed to create messages store: %s", err)
	}

	// create counts store
	if err := countsStore.Set(topic, 1); err != nil {
		return fmt.Errorf("failed to create counts store: %s", err)
	}

	// create subscribes store
	subscribe := safe.NewMap[string, func(message *Message)](func(mc *safe.MapConfig) {
		mc.Capacity = DefaultSubscruberCapacityForEachTopic
	})
	if err := subscribersStore.Set(topic, subscribe); err != nil {
		return fmt.Errorf("failed to create subscribes store: %s", err)
	}

	// create files store
	filedir := "/tmp/go-idp/logs"
	if err := fs.Mkdirp(filedir); err != nil {
		return fmt.Errorf("failed to create file dir: %s", err)
	}
	filepath := fmt.Sprintf("%s/%s.log", filedir, topic)
	writer, err := fs.CreateFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %s", err)
	}
	file := &File{
		Path:   filepath,
		Writer: writer,
		Clean: func() error {
			return fs.RemoveFile(filepath)
		},
		GetReader: func() (io.ReadCloser, error) {
			return fs.Open(filepath)
		},
	}
	if err := filesStore.Set(topic, file); err != nil {
		return fmt.Errorf("failed to create files store: %s", err)
	}

	return nil
}

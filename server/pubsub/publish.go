package pubsub

import (
	"context"
	"fmt"
	"time"
)

func Publish(ctx context.Context, topic string, message string) error {
	if ok := topicStore.Has(topic); !ok {
		return fmt.Errorf("topic %s not found", topic)
	}

	messages := messagesStore.Get(topic)
	if messages == nil {
		return fmt.Errorf("messages store for topic %s not found", topic)
	}
	counts := countsStore.Get(topic)
	if counts == 0 {
		return fmt.Errorf("counts store for topic %s not found", topic)
	}
	subscribes := subscribersStore.Get(topic)
	if subscribes == nil {
		return fmt.Errorf("subscribes store for topic %s not found", topic)
	}
	file := filesStore.Get(topic)

	counts++
	msg := &Message{
		ID:        counts,
		Content:   message,
		Timestamp: time.Now().UnixMilli(),
	}

	// store message
	messages.Enqueue(msg)
	// store count
	countsStore.Set(topic, counts)

	// notify subscribers
	subscribes.ForEach(func(uid string, handler func(message *Message)) (stop bool) {
		handler(msg)
		return
	})

	if message != "[DONE]" {
		// write to file
		// _, err := file.Writer.Write([]byte(fmt.Sprintf("%d %s %d\n", msg.ID, msg.Content, msg.Timestamp)))
		_, err := file.Writer.Write([]byte(msg.Content))
		if err != nil {
			return fmt.Errorf("failed to write to file: %s", err)
		}
	}

	return messagesStore.Set(topic, messages)
}

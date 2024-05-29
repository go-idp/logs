package pubsub

import (
	"context"
	"fmt"
)

func Close(ctx context.Context, topic string) error {
	if ok := topicStore.Has(topic); !ok {
		return fmt.Errorf("topic %s not found", topic)
	}

	file := filesStore.Get(topic)
	file.Writer.Close()
	// // clean file
	// go func() {
	// 	time.Sleep(2 * time.Hour)
	// 	file.Clean()
	// }()

	if err := Publish(ctx, topic, "[DONE]"); err != nil {
		return fmt.Errorf("failed to publish done message: %s", err)
	}

	if err := topicStore.Del(topic); err != nil {
		return fmt.Errorf("failed to delete topic store: %s", err)
	}

	if err := messagesStore.Del(topic); err != nil {
		return fmt.Errorf("failed to delete messages store: %s", err)
	}

	if err := countsStore.Del(topic); err != nil {
		return fmt.Errorf("failed to delete counts store: %s", err)
	}

	if err := subscribersStore.Del(topic); err != nil {
		return fmt.Errorf("failed to delete subscribes store: %s", err)
	}

	if err := filesStore.Del(topic); err != nil {
		return fmt.Errorf("failed to delete files store: %s", err)
	}

	return nil
}

package pubsub

import (
	"context"
	"fmt"

	"github.com/go-zoox/uuid"
)

func Subscribe(ctx context.Context, topic string, handler func(message *Message)) error {
	if ok := messagesStore.Has(topic); !ok {
		return fmt.Errorf("topic %s not found", topic)
	}

	messages := messagesStore.Get(topic)
	if messages == nil {
		return fmt.Errorf("messages store for topic %s not found", topic)
	}
	subscribes := subscribersStore.Get(topic)
	if subscribes == nil {
		return fmt.Errorf("subscribes store for topic %s not found", topic)
	}

	// store subscribe
	uid := uuid.V4()
	subscribes.Set(uid, handler)
	go func() {
		<-ctx.Done()
		subscribes.Del(uid)
	}()

	// notify messages
	messages.ForEach(func(message *Message, index int) (stop bool) {
		handler(message)
		return
	})

	return nil
}

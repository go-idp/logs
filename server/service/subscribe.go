package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-idp/logs/server/pubsub"
)

func (s *service) Subscribe(ctx context.Context, id string, handler func(err error, message string)) error {
	if s.cfg == nil {
		return fmt.Errorf("service is not setup")
	}

	if id == "" {
		return fmt.Errorf("id is required")
	}

	done, cancel := context.WithCancel(context.Background())

	err := pubsub.Subscribe(ctx, id, func(message *pubsub.Message) {
		msg, err := json.Marshal(map[string]any{
			"id":  message.ID,
			"log": message.Content,
			"ts":  message.Timestamp,
		})
		if err != nil {
			handler(fmt.Errorf("failed to marshal message: %s", err), "")
			return
		}

		handler(nil, string(msg))

		// finished
		if message.Content == "[DONE]" {
			cancel()
		}
	})
	if err != nil {
		// errMessage := map[string]any{
		// 	"code":    400,
		// 	"message": err.Error(),
		// }
		// errMessageBytes, errx := json.Marshal(errMessage)
		// if errx != nil {
		// 	return fmt.Errorf("failed to marshal error message: %s", errx)
		// }

		handler(err, "")
		// handler(nil, "[DONE]")
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-done.Done():
			return nil
		}
	}
}

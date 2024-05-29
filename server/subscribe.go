package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-idp/logs/server/pubsub"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/zoox"
)

func Subscribe() func(ctx *zoox.Context) {
	return func(ctx *zoox.Context) {
		id := ctx.Param().Get("id").String()
		if id == "" {
			ctx.Fail(nil, http.StatusBadRequest, "id is required")
		}

		done, cancel := context.WithCancel(context.Background())

		handler := func(message *pubsub.Message) {
			msg, err := json.Marshal(map[string]any{
				"id":  message.ID,
				"log": message.Content,
				"ts":  message.Timestamp,
			})
			if err != nil {
				ctx.Fail(err, http.StatusInternalServerError, fmt.Sprintf("failed to marshal message: %s", err))
				return
			}

			ctx.SSE().Event("message", string(msg))

			// finished
			if message.Content == "[DONE]" {
				cancel()
			}
		}

		if err := pubsub.Subscribe(ctx.Context(), id, handler); err != nil {
			ctx.Fail(err, http.StatusInternalServerError, fmt.Sprintf("failed to subscribe topic: %s", err))
			return
		}

		for {
			select {
			case <-ctx.Context().Done():
				return
			case <-done.Done():
				ctx.Success(zoox.H{
					"code":    200,
					"message": "success",
				})
				return
			}
		}
	}
}

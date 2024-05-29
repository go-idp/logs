package server

import (
	"net/http"

	"github.com/go-idp/logs/server/pubsub"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/zoox"
)

type PublishDTO struct {
	Messages []string `body:"messages"`
}

func Publish() func(ctx *zoox.Context) {
	return func(ctx *zoox.Context) {
		id := ctx.Param().Get("id").String()
		if id == "" {
			ctx.Fail(nil, http.StatusBadRequest, "id is required")
		}

		data := &PublishDTO{}
		if err := ctx.BindBody(data); err != nil {
			ctx.Fail(err, http.StatusBadRequest, fmt.Sprintf("failed to bind body: %s", err))
			return
		}
		if len(data.Messages) == 0 {
			ctx.Fail(nil, http.StatusBadRequest, "messages is required")
			return
		}

		for _, message := range data.Messages {
			if err := pubsub.Publish(ctx.Context(), id, message); err != nil {
				ctx.Fail(err, http.StatusInternalServerError, fmt.Sprintf("failed to publish topic: %s", err))
				return
			}
		}

		ctx.Success(nil)
	}
}

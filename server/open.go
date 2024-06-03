package server

import (
	"net/http"

	"github.com/go-idp/logs/server/pubsub"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/datetime"
	"github.com/go-zoox/zoox"
)

func Open() func(ctx *zoox.Context) {
	return func(ctx *zoox.Context) {
		id := ctx.Param().Get("id").String()
		if id == "" {
			ctx.Fail(nil, http.StatusBadRequest, "id is required")
		}

		if err := pubsub.Open(ctx.Context(), id); err != nil {
			ctx.Fail(err, http.StatusInternalServerError, fmt.Sprintf("failed to create topic: %s", err))
			return
		}

		welcomMessage := fmt.Sprintf("[%s][ID: %s] ...", datetime.Now().Format("YYYY-MM-DD HH:mm:ss"), id)
		if err := pubsub.Publish(ctx.Context(), id, welcomMessage); err != nil {
			ctx.Fail(err, http.StatusInternalServerError, fmt.Sprintf("failed to publish topic: %s", err))
			return
		}

		ctx.Success(nil)
	}
}

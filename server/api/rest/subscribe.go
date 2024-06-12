package rest

import (
	"net/http"

	"github.com/go-idp/logs/server/service"
	"github.com/go-zoox/zoox"
)

func Subscribe() func(ctx *zoox.Context) {
	return func(ctx *zoox.Context) {
		id := ctx.Param().Get("id").String()
		if id == "" {
			ctx.Fail(nil, http.StatusBadRequest, "id is required")
		}

		// ctx.SSE().Retry(10)

		err := service.Get().Subscribe(ctx.Context(), id, func(err error, message string) {
			if err != nil {
				ctx.Fail(err, http.StatusBadRequest, err.Error())
				return
			}

			ctx.SSE().Event("message", message)
		})
		if err != nil {
			ctx.Fail(err, http.StatusBadRequest, err.Error())
		}
	}
}

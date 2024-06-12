package rest

import (
	"net/http"

	"github.com/go-idp/logs/server/service"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/zoox"
)

func Finish() func(ctx *zoox.Context) {
	return func(ctx *zoox.Context) {
		id := ctx.Param().Get("id").String()
		if id == "" {
			ctx.Fail(nil, http.StatusBadRequest, "id is required")
		}

		if err := service.Get().Finish(ctx.Context(), id); err != nil {
			ctx.Fail(err, http.StatusInternalServerError, fmt.Sprintf("failed to finish: %s", err))
			return
		}

		ctx.Success(nil)
	}
}

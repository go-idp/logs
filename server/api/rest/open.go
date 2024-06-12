package rest

import (
	"net/http"

	"github.com/go-idp/logs/server/service"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/zoox"
)

func Open() func(ctx *zoox.Context) {
	return func(ctx *zoox.Context) {
		id := ctx.Param().Get("id").String()

		if err := service.Get().Open(ctx.Context(), id); err != nil {
			ctx.Fail(err, http.StatusInternalServerError, fmt.Sprintf("failed to open: %s", err))
			return
		}

		ctx.Success(nil)
	}
}

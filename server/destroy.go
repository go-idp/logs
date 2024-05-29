package server

import (
	"net/http"

	"github.com/go-idp/logs/config"
	"github.com/go-idp/logs/server/pubsub"
	"github.com/go-idp/logs/server/storage/fs"
	"github.com/go-idp/logs/server/storage/oss"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/zoox"
)

func Destroy() func(ctx *zoox.Context) {
	cfg := config.Get()
	return func(ctx *zoox.Context) {
		id := ctx.Param().Get("id").String()
		if id == "" {
			ctx.Fail(nil, http.StatusBadRequest, "id is required")
		}

		if err := pubsub.Destroy(ctx.Context(), id); err != nil {
			ctx.Fail(err, http.StatusInternalServerError, fmt.Sprintf("failed to destroy topic: %s", err))
			return
		}

		//
		file, err := pubsub.GetFile(id)
		if err != nil {
			ctx.Logger.Errorf("failed to get file path: %s (err: %s)", ctx.Path, err)
			return
		}
		reader, err := file.GetReader()
		if err != nil {
			ctx.Logger.Errorf("failed to get file reader: %s (err: %s)", ctx.Path, err)
			return
		}
		defer reader.Close()

		switch cfg.Storage.Driver {
		case "oss":
			err := oss.Get().Put(ctx.Path, reader)
			if err != nil {
				ctx.Logger.Errorf("failed to get file path: %s (err: %s)", ctx.Path, err)
				ctx.Error(http.StatusNotFound, "Not Found")
				return
			}
		default:
			err = fs.Get().Put(ctx.Path, reader)
			if err != nil {
				ctx.Logger.Errorf("failed to get file path: %s (err: %s)", ctx.Path, err)
				ctx.Error(http.StatusNotFound, "Not Found")
				return
			}
		}

		ctx.Success(nil)
	}
}

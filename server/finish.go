package server

import (
	"net/http"

	"github.com/go-idp/logs/server/config"
	"github.com/go-idp/logs/server/pubsub"
	"github.com/go-idp/logs/server/storage/fs"
	"github.com/go-idp/logs/server/storage/oss"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/zoox"
)

func Finish() func(ctx *zoox.Context) {
	cfg := config.Get()
	return func(ctx *zoox.Context) {
		id := ctx.Param().Get("id").String()
		if id == "" {
			ctx.Fail(nil, http.StatusBadRequest, "id is required")
		}

		// get file
		file, err := pubsub.GetFile(id)
		if err != nil {
			ctx.Fail(err, http.StatusInternalServerError, fmt.Sprintf("failed to get file path (topic: %s, err: %s)", id, err))
			return
		}
		defer file.Clean()

		if err := pubsub.Close(ctx.Context(), id); err != nil {
			ctx.Fail(err, http.StatusInternalServerError, fmt.Sprintf("failed to destroy topic(%s): %s", id, err))
			return
		}

		reader, err := file.GetReader()
		if err != nil {
			ctx.Fail(err, http.StatusInternalServerError, fmt.Sprintf("failed to get file reader (topic: %s, err: %s)", id, err))
			return
		}
		defer reader.Close()

		switch cfg.Storage.Driver {
		case "oss":
			err := oss.Get().Put(id, reader)
			if err != nil {
				ctx.Fail(err, http.StatusInternalServerError, fmt.Sprintf("failed to get file path from oss (topic: %s, err: %s)", id, err))
				return
			}
		default:
			err = fs.Get().Put(id, reader)
			if err != nil {
				ctx.Fail(err, http.StatusInternalServerError, fmt.Sprintf("failed to get file path from fs (topic: %s, err: %s)", id, err))
				return
			}
		}

		ctx.Success(nil)
	}
}

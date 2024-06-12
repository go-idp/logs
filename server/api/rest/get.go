package rest

import (
	"io"
	"net/http"
	"time"

	"github.com/go-idp/logs/server/config"
	"github.com/go-idp/logs/server/storage/fs"
	"github.com/go-idp/logs/server/storage/oss"
	"github.com/go-zoox/headers"
	"github.com/go-zoox/zoox"
)

func Get() func(ctx *zoox.Context) {
	cfg := config.Get()

	return func(ctx *zoox.Context) {
		if ctx.Path == "" {
			ctx.Error(http.StatusNotFound, "Not Found")
			return
		}

		var reader io.ReadCloser
		var err error

		switch cfg.Storage.Driver {
		case "oss":
			reader, err = oss.Get().Get(ctx.Path)
			if err != nil {
				ctx.Logger.Errorf("failed to get file path from fs: %s (err: %s)", ctx.Path, err)
				ctx.Error(http.StatusNotFound, "Not Found")
				return
			}
		default:
			reader, err = fs.Get().Get(ctx.Path)
			if err != nil {
				ctx.Logger.Errorf("failed to get file path from oss: %s (err: %s)", ctx.Path, err)
				ctx.Error(http.StatusNotFound, "Not Found")
				return
			}
		}
		defer reader.Close()

		ctx.SetCacheControlWithMaxAge(365 * 24 * time.Hour)
		ctx.Set(headers.Vary, "origin")

		if _, err := io.Copy(ctx.Writer, reader); err != nil {
			ctx.Logger.Errorf("failed to send file reader: %s (err: %s)", ctx.Path, err)
		}
	}
}

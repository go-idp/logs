package server

import (
	"github.com/go-idp/logs/server/config"
	"github.com/go-zoox/zoox"
)

func Auth() zoox.Middleware {
	cfg := config.Get()
	isAuthEnabled := cfg.Auth.Username != "" && cfg.Auth.Password != ""
	return func(ctx *zoox.Context) {
		if isAuthEnabled {
			user, pass, ok := ctx.Request.BasicAuth()
			if !ok {
				ctx.Set("WWW-Authenticate", `Basic realm="go-idp"`)
				ctx.Status(401)
				return
			}

			if !(user == cfg.Auth.Username && pass == cfg.Auth.Password) {
				ctx.Status(401)
				return
			}
		}

		ctx.Next()
	}
}

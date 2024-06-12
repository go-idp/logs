package rest

import (
	"github.com/go-zoox/zoox"
)

func Stream() func(ctx *zoox.Context) {
	return Subscribe()
}

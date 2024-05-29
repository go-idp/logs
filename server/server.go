package server

import (
	"github.com/go-idp/logs/config"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/zoox/defaults"
)

func New(cfg *config.Config) error {
	app := defaults.Default()

	//
	app.Post("/:id/create", Create())
	app.Post("/:id/destroy", Destroy())
	//
	app.Post("/:id/publish", Publish())
	app.Post("/:id/subscribe", Subscribe())
	//
	app.Get("/:id/stream", Stream())

	//
	// app.Get("/:id", Get())

	return app.Run(fmt.Sprintf(":%d", cfg.Port))
}

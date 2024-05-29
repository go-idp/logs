package server

import (
	"github.com/go-idp/logs"
	"github.com/go-idp/logs/server/config"
	"github.com/go-idp/logs/server/storage/fs"
	"github.com/go-idp/logs/server/storage/oss"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/zoox"
	"github.com/go-zoox/zoox/defaults"
)

type Server interface {
	Run() error
}

type server struct {
	cfg *config.Config
}

func New() (Server, error) {
	cfg := config.Get()
	fs.Get().Setup(func(c *fs.Config) {
		c.RootDIR = cfg.Storage.RootDIR
	})
	oss.Get().SetUp(func(c *oss.Config) {
		c.RootDIR = cfg.Storage.RootDIR
		c.AccessKeyID = cfg.Storage.OSSAccessKeyID
		c.AccessKeySecret = cfg.Storage.OSSAccessKeySecret
		c.Bucket = cfg.Storage.OSSBucket
		c.Endpoint = cfg.Storage.OSSEndpoint
	})

	s := &server{
		cfg: cfg,
	}

	return s, nil
}

func (s *server) Run() error {
	app := defaults.Default()

	app.Use(Auth())

	//
	app.Post("/:id/open", Open())
	app.Post("/:id/finish", Finish())
	//
	app.Post("/:id/publish", Publish())
	app.Post("/:id/subscribe", Subscribe())
	//
	app.Get("/:id/stream", Stream())

	//
	app.Get("/:id", Get())

	app.Get("/", func(ctx *zoox.Context) {
		ctx.JSON(200, zoox.H{
			"name":    "logs service for idp",
			"version": logs.Version,
		})
	})

	return app.Run(fmt.Sprintf(":%d", s.cfg.Port))
}

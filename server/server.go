package server

import (
	"github.com/go-idp/logs"
	"github.com/go-idp/logs/server/api/rest"
	"github.com/go-idp/logs/server/api/ws"
	"github.com/go-idp/logs/server/config"
	"github.com/go-idp/logs/server/pubsub"
	"github.com/go-idp/logs/server/service"
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

	{ // rest
		app.Post("/:id/open", rest.Open())
		app.Post("/:id/finish", rest.Finish())
		//
		app.Post("/:id/publish", rest.Publish())
		app.Post("/:id/subscribe", rest.Subscribe())
		//
		app.Get("/:id/stream", rest.Stream())
	}

	{ // websocket
		s, err := app.WebSocket("/")
		if err != nil {
			return fmt.Errorf("failed to create websocket server: %s", err)
		}

		//
		s.Event("open", ws.Open())
		s.Event("finish", ws.Finish())
		s.Event("publish", ws.Publish())
		s.Event("subscribe", ws.Subscribe())
	}

	//
	app.Get("/:id", rest.Get())

	app.Get("/", func(ctx *zoox.Context) {
		ctx.JSON(200, zoox.H{
			"name":    "logs service for idp",
			"version": logs.Version,
			"logs":    pubsub.GetTopics(),
		})
	})

	service.Get().Setup(s.cfg)

	return app.Run(fmt.Sprintf(":%d", s.cfg.Port))
}

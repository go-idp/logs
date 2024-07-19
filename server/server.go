package server

import (
	"github.com/go-idp/logs"
	"github.com/go-idp/logs/server/api/rest"
	"github.com/go-idp/logs/server/api/ws"
	"github.com/go-idp/logs/server/config"
	"github.com/go-idp/logs/server/service"
	"github.com/go-idp/logs/server/storage/fs"
	"github.com/go-idp/logs/server/storage/oss"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/datetime"
	es "github.com/go-zoox/websocket/extension/event/server"
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

	// objects: list + retrieve
	app.Group("/objects", func(group *zoox.RouterGroup) {
		group.Get("/", func(ctx *zoox.Context) {
			tasks, err := service.Get().Data().List()
			if err != nil {
				ctx.Error(500, err.Error())
				return
			}

			ctx.Success(zoox.H{
				"data":  tasks,
				"total": len(tasks),
			})
		})

		group.Get("/:id", func(ctx *zoox.Context) {
			id := ctx.Param().Get("id").String()
			if id == "" {
				ctx.Error(400, "id is required")
				return
			}

			task, err := service.Get().Data().Retrieve(id)
			if err != nil {
				ctx.Error(404, err.Error())
				return
			}

			ctx.Success(task)
		})
	})

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

		et := es.New(s)
		et.On("open", ws.Open())
		et.On("finish", ws.Finish())
		et.On("publish", ws.Publish())
		et.On("subscribe", ws.Subscribe())
	}

	//
	app.Get("/:id", rest.Get())

	app.Get("/", func(ctx *zoox.Context) {
		ctx.JSON(200, zoox.H{
			"name":       "logs service for idp",
			"version":    logs.Version,
			"status":     service.Get().Status(),
			"running_at": datetime.Now().Format("YYYY-MM-DD HH:mm:ss"),
		})
	})

	service.Get().Setup(s.cfg)

	return app.Run(fmt.Sprintf(":%d", s.cfg.Port))
}

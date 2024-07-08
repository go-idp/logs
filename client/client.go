package client

import (
	"context"
	"fmt"
	gurl "net/url"
	"time"

	"github.com/go-zoox/core-utils/safe"
	"github.com/go-zoox/websocket"

	ec "github.com/go-zoox/websocket/extension/event/client"
)

type Client interface {
	Connect() error
	Close() error
	//
	Open(ctx context.Context, id string) error
	Finish(ctx context.Context, id string) error
	//
	Publish(ctx context.Context, id string, message string) error
	Subscribe(ctx context.Context, id string, fn func(message string)) error
}

type client struct {
	cfg *Config

	ws    websocket.Client
	event ec.Client

	publishStore *safe.Map[string, *publishTopic]
}

type Option func(cfg *Config)

func New(opts ...Option) (c Client, err error) {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}

	// fix the server URL
	cfg.Server, err = GetServerURL(cfg.Engine, cfg.Server)
	if err != nil {
		return nil, err
	}

	var event ec.Client
	var ws websocket.Client
	if cfg.Engine == "websocket" {
		ws, err = websocket.NewClient(func(opt *websocket.ClientOption) {
			// opt.Context = ctx
			opt.Addr = cfg.Server
			opt.ConnectTimeout = 10 * time.Second
		})
		if err != nil {
			return nil, err
		}

		event = ec.New(ws)
	}

	return &client{
		cfg:   cfg,
		ws:    ws,
		event: event,
		//
		publishStore: safe.NewMap[string, *publishTopic](),
	}, nil
}

func GetServerURL(engine string, url string) (string, error) {
	u, err := gurl.Parse(url)
	if err != nil {
		return "", err
	}

	switch engine {
	case "http":
		return url, nil
	case "websocket":
		// if https
		if u.Scheme == "https" {
			u.Scheme = "wss"
		} else {
			u.Scheme = "ws"
		}

		return u.String(), nil
	case "tcp":
		return url, nil
	case "grpc":
		return url, nil
	default:
		return "", fmt.Errorf("unsupported engine: %s, only support http, websocket, tcp, and grpc", engine)
	}
}

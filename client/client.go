package client

import (
	"context"
	"fmt"
	gurl "net/url"
)

type Client interface {
	Open(ctx context.Context, id string) error
	Finish(ctx context.Context, id string) error
	//
	Publish(ctx context.Context, id string, message string) error
	Subscribe(ctx context.Context, id string, fn func(message string)) error
}

type client struct {
	cfg *Config
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

	return &client{
		cfg: cfg,
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

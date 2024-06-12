package client

import "context"

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

func New(opts ...Option) (Client, error) {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}

	// fix the server URL
	cfg.Server = GetServerURL(cfg.Engine, cfg.Server)

	return &client{
		cfg: cfg,
	}, nil
}

func GetServerURL(engine string, url string) string {
	switch engine {
	case "http":
		return url
	case "websocket":
		// if https
		if url[:5] == "https" {
			return "wss" + url[5:]
		} else if url[:4] == "http" {
			// http
			return "ws" + url[4:]
		} else {
			return url
		}
	case "tcp":
		return "tcp" + url[4:]
	case "grpc":
		return "grpc" + url[4:]
	default:
		return url
	}
}

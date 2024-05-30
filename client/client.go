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

	return &client{
		cfg: cfg,
	}, nil
}

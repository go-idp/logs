package client

type Client interface {
	Open(id string) error
	Finish(id string) error
	//
	Publish(id string, message string) error
	Subscribe(id string, fn func(message string)) error
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

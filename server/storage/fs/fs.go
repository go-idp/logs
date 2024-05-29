package fs

import (
	"io"

	"github.com/go-idp/logs/server/storage"
	"github.com/go-zoox/fs"
	"github.com/go-zoox/once"
)

type FS struct {
	cfg *Config
}

type Config struct {
	RootDIR string
}

type Option func(cfg *Config)

func New() *FS {
	return &FS{}
}

func (o *FS) Get(path string) (io.ReadCloser, error) {
	if o.cfg == nil {
		return nil, storage.ErrNotSetup
	}

	fullpath := fs.JoinPath(o.cfg.RootDIR, path)
	return fs.Open(fullpath)
}

func (o *FS) Put(path string, stream io.Reader) error {
	if o.cfg == nil {
		return storage.ErrNotSetup
	}

	fullpath := fs.JoinPath(o.cfg.RootDIR, path)
	f, err := fs.CreateFile(fullpath)
	if err != nil {
		return err
	}

	if _, err := io.Copy(f, stream); err != nil {
		return err
	}

	return nil
}

func (o *FS) Setup(opts ...Option) error {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}

	o.cfg = cfg

	return nil
}

func Get() *FS {
	return once.Get("storage.fs", New)
}

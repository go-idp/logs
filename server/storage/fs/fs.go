package fs

import (
	"io"

	"github.com/go-idp/logs/server/storage"
	"github.com/go-zoox/fs"
	"github.com/go-zoox/once"
)

type FS struct {
	cfg *FSConfig
}

type FSConfig struct {
	RootDIR string
}

type FSOption func(cfg *FSConfig)

func New() storage.Storage {
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

func (o *FS) Setup(opts ...FSOption) error {
	cfg := &FSConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	o.cfg = cfg

	return nil
}

func Get() storage.Storage {
	return once.Get("storage.fs", New)
}

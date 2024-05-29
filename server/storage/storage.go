package storage

import (
	"errors"
	"io"
)

type Storage interface {
	Get(path string) (io.ReadCloser, error)
	Put(path string, stream io.Reader) error
}

var ErrNotSetup = errors.New("not setup")

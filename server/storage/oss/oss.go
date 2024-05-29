package oss

import (
	"io"

	alioss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-idp/logs/server/storage"
	"github.com/go-zoox/core-utils/regexp"
	"github.com/go-zoox/fs"
	"github.com/go-zoox/once"
)

type OSS struct {
	cfg *OSSConfig

	bucket *alioss.Bucket
}

type OSSConfig struct {
	AccessKeyID     string
	AccessKeySecret string
	Bucket          string
	Endpoint        string
	RootDIR         string
}

type OSSOption func(cfg *OSSConfig)

func New() storage.Storage {
	return &OSS{}
}

func (o *OSS) Get(path string) (io.ReadCloser, error) {
	if o.cfg == nil {
		return nil, storage.ErrNotSetup
	}

	fullpath := fs.JoinPath(o.cfg.RootDIR, path)
	return o.bucket.GetObject(fullpath)
}

func (o *OSS) Put(path string, stream io.Reader) error {
	if o.cfg == nil {
		return storage.ErrNotSetup
	}

	fullpath := fs.JoinPath(o.cfg.RootDIR, path)
	return o.bucket.PutObject(fullpath, stream)
}

func (o *OSS) SetUp(opts ...OSSOption) error {
	cfg := &OSSConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	// fix root dir path
	// remove prefix slash
	if matched := regexp.Match("^/", cfg.RootDIR); matched {
		cfg.RootDIR = cfg.RootDIR[1:]
	}

	o.cfg = cfg

	client, err := alioss.New(
		cfg.Endpoint,
		cfg.AccessKeyID,
		cfg.AccessKeySecret,
	)
	if err != nil {
		return err
	}

	bucket, err := client.Bucket(cfg.Bucket)
	if err != nil {
		return err
	}

	o.bucket = bucket

	return nil
}

func Get() storage.Storage {
	return once.Get("storage.oss", New)
}

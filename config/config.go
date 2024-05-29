package config

import (
	"github.com/go-zoox/once"
)

type Config struct {
	Port    int
	Storage Storage `json:"storage"`
}

type Storage struct {
	Driver string `json:"driver"` // oss, fs, default: fs

	RootDIR string `json:"root_dir"`

	OSSAccessKeyID     string
	OSSAccessKeySecret string
	OSSBucket          string
	OSSEndpoint        string
}

func New() *Config {
	return &Config{}
}

func Get() *Config {
	return once.Get("config", New)
}

package config

import (
	"github.com/go-zoox/once"
)

type Config struct {
	Port    int     `json:"port"`
	Storage Storage `json:"storage"`
	Auth    Auth    `json:"auth"`
}

type Storage struct {
	Driver string `json:"driver"` // oss, fs, default: fs

	RootDIR string `json:"root_dir"`

	OSSAccessKeyID     string `json:"oss_access_key_id"`
	OSSAccessKeySecret string `json:"oss_access_key_secret"`
	OSSBucket          string `json:"oss_bucket"`
	OSSEndpoint        string `json:"oss_endpoint"`
}

type Auth struct {
	Username string
	Password string
}

func New() *Config {
	return &Config{
		Port: 8080,
		Storage: Storage{
			Driver:  "fs",
			RootDIR: "/data",
		},
	}
}

func Get() *Config {
	return once.Get("config", New)
}

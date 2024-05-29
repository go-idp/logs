package server

import (
	"github.com/go-idp/logs/server"
	"github.com/go-idp/logs/server/config"
	"github.com/go-zoox/cli"
	"github.com/go-zoox/fs"
)

func Register(app *cli.MultipleProgram) {
	app.Register("server", &cli.Command{
		Name:  "server",
		Usage: "the logs server",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Usage:   "server port",
				Aliases: []string{"p"},
				EnvVars: []string{"PORT"},
				Value:   8080,
			},
			&cli.StringFlag{
				Name:    "storage-driver",
				Usage:   "Driver for Storage, available: fs, oss",
				EnvVars: []string{"STORAGE_DRIVER"},
				Value:   "fs",
			},
			&cli.StringFlag{
				Name:    "storage-root-dir",
				Usage:   "Root Directory for Storage",
				EnvVars: []string{"STORAGE_ROOT_DIR"},
			},
			&cli.StringFlag{
				Name:    "storage-oss-access-key-id",
				Usage:   "OSS Acess Key ID for Storage",
				EnvVars: []string{"STORAGE_OSS_ACCESS_KEY_ID"},
				// Required: true,
			},
			&cli.StringFlag{
				Name:    "storage-oss-access-key-secret",
				Usage:   "OSS Acess Key Secret for Storage",
				EnvVars: []string{"STORAGE_OSS_ACCESS_KEY_SECRET"},
				// Required: true,
			},
			&cli.StringFlag{
				Name:    "storage-oss-bucket",
				Usage:   "OSS Bucket for Storage",
				EnvVars: []string{"STORAGE_OSS_BUCKET"},
				// Required: true,
			},
			&cli.StringFlag{
				Name:    "storage-oss-endpoint",
				Usage:   "OSS Endpoint for Storage",
				EnvVars: []string{"STORAGE_OSS_ENDPOINT"},
				// Required: true,
			},
			&cli.StringFlag{
				Name:    "username",
				Usage:   "Username for Basic Auth",
				EnvVars: []string{"USERNAME"},
			},
			&cli.StringFlag{
				Name:    "password",
				Usage:   "Password for Basic Auth",
				EnvVars: []string{"PASSWORD"},
			},
		},
		Action: func(ctx *cli.Context) (err error) {
			cfg := config.Get()
			// Port
			cfg.Port = ctx.Int("port")
			// Storage
			cfg.Storage.Driver = ctx.String("storage-driver")
			cfg.Storage.RootDIR = ctx.String("storage-root-dir")
			cfg.Storage.OSSAccessKeyID = ctx.String("storage-oss-access-key-id")
			cfg.Storage.OSSAccessKeySecret = ctx.String("storage-oss-access-key-secret")
			cfg.Storage.OSSBucket = ctx.String("storage-oss-bucket")
			cfg.Storage.OSSEndpoint = ctx.String("storage-oss-endpoint")
			// Auth
			cfg.Auth.Username = ctx.String("username")
			cfg.Auth.Password = ctx.String("password")

			if cfg.Storage.Driver == "fs" {
				if cfg.Storage.RootDIR == "" {
					cfg.Storage.RootDIR = fs.JoinCurrentDir("data")
				}
			} else {
				if cfg.Storage.RootDIR == "" {
					cfg.Storage.RootDIR = "/data"
				}
			}

			s, err := server.New()
			if err != nil {
				return err
			}

			return s.Run()
		},
	})
}

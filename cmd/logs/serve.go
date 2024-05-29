package main

import (
	"github.com/go-idp/logs/config"
	"github.com/go-idp/logs/server"
	"github.com/go-zoox/cli"
)

func registerServe(app *cli.MultipleProgram) {
	app.Register("serve", &cli.Command{
		Name:  "serve",
		Usage: "Start the logs service",
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
		},
		Action: func(ctx *cli.Context) (err error) {
			return server.New(&config.Config{
				Port: ctx.Int("port"),
				Storage: config.Storage{
					Driver: ctx.String("storage-driver"),
					//
					RootDIR: ctx.String("storage-root-dir"),
					//
					OSSAccessKeyID:     ctx.String("storage-oss-access-key-id"),
					OSSAccessKeySecret: ctx.String("storage-oss-access-key-secret"),
					OSSBucket:          ctx.String("storage-oss-bucket"),
					OSSEndpoint:        ctx.String("storage-oss-endpoint"),
				},
			})
		},
	})
}

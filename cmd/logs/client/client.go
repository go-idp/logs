package client

import (
	"github.com/go-zoox/cli"
)

func Register(app *cli.MultipleProgram) {
	app.Register("client", &cli.Command{
		Name:  "client",
		Usage: "the logs client",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "server",
				Usage:   "logs server",
				Aliases: []string{"s"},
				EnvVars: []string{"SERVER"},
				Value:   "http://127.0.0.1:8080",
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
			&cli.StringFlag{
				Name:    "engine",
				Usage:   "engine to use, avaliable: http, websocket, tcp, grpc",
				EnvVars: []string{"ENGINE"},
				Value:   "http",
			},
		},
		Subcommands: []*cli.Command{
			Open(),
			Finish(),
			Publish(),
			Subscribe(),
			//
			Pipe(),
		},
	})
}

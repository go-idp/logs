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

func GetServerURL(engine string, url string) string {
	switch engine {
	case "http":
		return url
	case "websocket":
		// if https
		if url[:5] == "https" {
			return "wss" + url[5:]
		} else if url[:4] == "http" {
			// http
			return "ws" + url[4:]
		} else {
			return url
		}
	case "tcp":
		return "tcp" + url[4:]
	case "grpc":
		return "grpc" + url[4:]
	default:
		return url
	}
}

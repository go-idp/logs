package main

import (
	"github.com/go-idp/logs"
	"github.com/go-idp/logs/cmd/logs/client"
	"github.com/go-idp/logs/cmd/logs/server"
	"github.com/go-zoox/cli"
)

func main() {
	app := cli.NewMultipleProgram(&cli.MultipleProgramConfig{
		Name:    "logs",
		Usage:   "Logs Service for IDP",
		Version: logs.Version,
	})

	server.Register(app)
	client.Register(app)

	app.Run()
}

package main

import (
	"github.com/go-idp/logs"
	"github.com/go-zoox/cli"
)

func main() {
	app := cli.NewMultipleProgram(&cli.MultipleProgramConfig{
		Name:    "logs",
		Usage:   "Logs Service for IDP",
		Version: logs.Version,
	})

	registerServe(app)

	app.Run()
}

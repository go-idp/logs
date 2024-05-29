package client

import (
	"fmt"
	"os"

	"github.com/go-idp/logs/client"
	"github.com/go-zoox/cli"
)

func Subscribe() *cli.Command {
	return &cli.Command{
		Name:      "subscribe",
		Usage:     "subscribe logs",
		ArgsUsage: "<id>",
		Action: func(ctx *cli.Context) (err error) {
			c, err := client.New(func(cfg *client.Config) {
				cfg.Server = ctx.String("server")
				cfg.Username = ctx.String("username")
				cfg.Password = ctx.String("password")
			})
			if err != nil {
				return err
			}

			id := ctx.Args().Get(0)
			if id == "" {
				// return cli.ShowCommandHelp(ctx, "subscribe")
				return fmt.Errorf("id is required")
			}

			return c.Subscribe(id, func(message string) {
				os.Stdout.WriteString(message)
			})
		},
	}
}

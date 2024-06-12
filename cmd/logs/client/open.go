package client

import (
	"context"
	"fmt"

	"github.com/go-idp/logs/client"
	"github.com/go-zoox/cli"
)

func Open() *cli.Command {
	return &cli.Command{
		Name:      "open",
		Usage:     "open logs",
		ArgsUsage: "<id>",
		Action: func(ctx *cli.Context) (err error) {
			c, err := client.New(func(cfg *client.Config) {
				cfg.Server = ctx.String("server")
				cfg.Username = ctx.String("username")
				cfg.Password = ctx.String("password")
				cfg.Engine = ctx.String("engine")
			})
			if err != nil {
				return err
			}

			id := ctx.Args().Get(0)
			if id == "" {
				// return cli.ShowCommandHelp(ctx, "open")
				return fmt.Errorf("id is required")
			}

			return c.Open(context.Background(), id)
		},
	}
}

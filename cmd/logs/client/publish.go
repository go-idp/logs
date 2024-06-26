package client

import (
	"context"
	"fmt"

	"github.com/go-idp/logs/client"
	"github.com/go-zoox/cli"
)

func Publish() *cli.Command {
	return &cli.Command{
		Name:      "publish",
		Usage:     "publish logs",
		ArgsUsage: "<id> <message>",
		Action: func(ctx *cli.Context) (err error) {
			id := ctx.Args().Get(0)
			message := ctx.Args().Get(1)
			if id == "" || message == "" {
				// return cli.ShowCommandHelp(ctx, "publish")
				return fmt.Errorf("id and message are required")
			}

			c, err := client.New(func(cfg *client.Config) {
				cfg.Server = ctx.String("server")
				cfg.Username = ctx.String("username")
				cfg.Password = ctx.String("password")
				cfg.Engine = ctx.String("engine")
			})
			if err != nil {
				return err
			}

			if err := c.Connect(); err != nil {
				return err
			}
			defer c.Close()

			return c.Publish(context.Background(), id, message)
		},
	}
}

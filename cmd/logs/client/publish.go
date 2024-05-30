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
			c, err := client.New(func(cfg *client.Config) {
				cfg.Server = ctx.String("server")
				cfg.Username = ctx.String("username")
				cfg.Password = ctx.String("password")
			})
			if err != nil {
				return err
			}

			id := ctx.Args().Get(0)
			message := ctx.Args().Get(1)
			if id == "" || message == "" {
				// return cli.ShowCommandHelp(ctx, "publish")
				return fmt.Errorf("id and message are required")
			}

			return c.Publish(context.Background(), id, message)
		},
	}
}

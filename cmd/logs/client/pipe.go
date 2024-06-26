package client

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/go-idp/logs/client"
	"github.com/go-zoox/cli"
)

func Pipe() *cli.Command {
	return &cli.Command{
		Name:      "pipe",
		Usage:     "pipe logs",
		ArgsUsage: "<id>",
		Action: func(ctx *cli.Context) (err error) {
			id := ctx.Args().Get(0)
			if id == "" {
				// return cli.ShowCommandHelp(ctx, "publish")
				return fmt.Errorf("id is required")
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

			for {
				buf := make([]byte, 1024)
				n, err := os.Stdin.Read(buf)
				if err == io.EOF {
					break
				} else if err != nil {
					return err
				}

				if n != 0 {
					message := string(buf[:n])
					if message == "" {
						continue
					}

					if err := c.Publish(context.Background(), id, message); err != nil {
						return err
					}
				}
			}

			return nil
		},
	}
}

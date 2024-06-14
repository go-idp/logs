package client

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/go-idp/logs/client"
	"github.com/go-zoox/cli"
)

func Subscribe() *cli.Command {
	return &cli.Command{
		Name:      "subscribe",
		Usage:     "subscribe logs",
		ArgsUsage: "<id...>",
		Action: func(ctx *cli.Context) (err error) {
			ids := ctx.Args().Slice()
			if len(ids) == 0 {
				// return cli.ShowCommandHelp(ctx, "subscribe")
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

			wg := &sync.WaitGroup{}
			for _, id := range ids {
				wg.Add(1)
				go func(id string) {
					err := c.Subscribe(context.Background(), id, func(message string) {
						os.Stdout.Write([]byte(message))
					})
					if err != nil {
						fmt.Println(err)
					}

					wg.Done()
				}(id)
			}

			wg.Wait()

			return nil
		},
	}
}

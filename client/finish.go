package client

import (
	"context"
	"fmt"

	"github.com/go-zoox/fetch"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/websocket"
	"github.com/go-zoox/websocket/event/cs"
)

func (c *client) Finish(ctx context.Context, id string) error {
	if c.cfg == nil {
		return fmt.Errorf("client is not setup")
	}

	if id == "" {
		return fmt.Errorf("id is required")
	}

	switch c.cfg.Engine {
	case "websocket":
		return c.finishWithWebsocket(ctx, id)
	case "http":
		return c.finishWithHTTP(ctx, id)
	default:
		return fmt.Errorf("unsupported engine: %s, only support websocket and http", c.cfg.Engine)
	}
}

func (c *client) finishWithHTTP(ctx context.Context, id string) error {
	response, err := fetch.Post(fmt.Sprintf("%s/:id/finish", c.cfg.Server), &fetch.Config{
		Context: ctx,
		Params: fetch.Params{
			"id": id,
		},
	})
	if err != nil {
		return err
	}

	// if response.StatusCode() != 200 {
	// 	return fmt.Errorf("failed to finish: %s", response.String())
	// }

	if response.Get("code").Int() != 200 {
		return fmt.Errorf("failed to finish: %s", response.Get("message").String())
	}

	return nil
}

func (c *client) finishWithWebsocket(ctx context.Context, id string) error {
	ws, err := websocket.NewClient(func(opt *websocket.ClientOption) {
		opt.Context = ctx
		opt.Addr = c.cfg.Server
	})
	if err != nil {
		return err
	}

	if err := ws.Connect(); err != nil {
		return err
	}
	defer ws.Close()

	err = ws.Event("finish", cs.EventPayload{
		"id": id,
	}, func(err error, payload cs.EventPayload) {
		if err != nil {
			logger.Infof("failed to finish: %s", err)
			return
		}
	})
	if err != nil {
		return err
	}

	return nil
}

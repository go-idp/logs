package client

import (
	"context"
	"fmt"

	"github.com/go-zoox/fetch"
	"github.com/go-zoox/websocket"
	"github.com/go-zoox/websocket/event/cs"
)

func (c *client) Publish(ctx context.Context, id string, message string) error {
	if c.cfg == nil {
		return fmt.Errorf("client is not setup")
	}

	if id == "" {
		return fmt.Errorf("id is required")
	}

	switch c.cfg.Engine {
	case "websocket":
		return c.publishWithWebsocket(ctx, id, message)
	case "http":
		return c.publishWithHTTP(ctx, id, message)
	default:
		return fmt.Errorf("unsupported engine: %s, only support websocket and http", c.cfg.Engine)
	}
}

func (c *client) publishWithHTTP(ctx context.Context, id string, message string) error {
	response, err := fetch.Post(fmt.Sprintf("%s/:id/publish", c.cfg.Server), &fetch.Config{
		Context: ctx,
		Params: fetch.Params{
			"id": id,
		},
		Headers: fetch.Headers{
			"Content-Type": "application/json",
		},
		Body: map[string]interface{}{
			"message": message,
		},
	})
	if err != nil {
		return err
	}

	// if response.StatusCode() != 200 {
	// 	return fmt.Errorf("failed to publish: %s", response.String())
	// }

	if response.Get("code").Int() != 200 {
		return fmt.Errorf("failed to publish: %s", response.Get("message").String())
	}

	return nil
}

func (c *client) publishWithWebsocket(ctx context.Context, id string, message string) error {
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

	err = ws.Event("publish", cs.EventPayload{
		"id":      id,
		"message": message,
	}, func(err error, payload cs.EventPayload) {
		if err != nil {
			// logger.Infof("failed to publish: %s", err)
			return
		}
	})
	if err != nil {
		return err
	}

	return nil
}

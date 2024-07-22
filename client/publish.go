package client

import (
	"context"
	"fmt"
	"time"

	"github.com/go-zoox/core-utils/strings"
	"github.com/go-zoox/fetch"
	"github.com/go-zoox/logger"
	cs "github.com/go-zoox/websocket/extension/event/entity"
)

type publishTopic struct {
	Data   strings.Builder
	Ticker *time.Ticker
}

func (c *client) Publish(ctx context.Context, id string, message string) error {
	if c.cfg == nil {
		return fmt.Errorf("client is not setup")
	}

	if id == "" {
		return fmt.Errorf("id is required")
	}

	topic := c.publishStore.Get(id)
	if topic == nil {
		return fmt.Errorf("cannot publish before open")
	}

	if _, err := topic.Data.Write([]byte(message)); err != nil {
		return err
	}

	return nil
}

func (c *client) flushPublish(ctx context.Context, id string) error {
	c.flushPublishLocker.Lock()
	defer c.flushPublishLocker.Unlock()

	topic := c.publishStore.Get(id)
	if topic != nil {
		if topic.Data.Len() == 0 {
			return nil
		}
		defer topic.Data.Reset()
		message := topic.Data.String()

		switch c.cfg.Engine {
		case "websocket":
			return c.publishWithWebsocket(ctx, id, message)
		case "http":
			return c.publishWithHTTP(ctx, id, message)
		default:
			return fmt.Errorf("unsupported engine: %s, only support websocket and http", c.cfg.Engine)
		}
	}

	return nil
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
	return c.event.Emit("publish", cs.EventPayload{
		"id":      id,
		"message": message,
	}, func(err error, payload cs.EventPayload) {
		if err != nil {
			logger.Infof("failed to publish: %s", err)
			return
		}
	})
}

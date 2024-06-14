package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-zoox/core-utils/strings"
	"github.com/go-zoox/fetch"
	ec "github.com/go-zoox/websocket/extension/event/client"
	cs "github.com/go-zoox/websocket/extension/event/entity"
)

func (c *client) Subscribe(ctx context.Context, id string, fn func(message string)) error {
	if c.cfg == nil {
		return fmt.Errorf("client is not setup")
	}

	if id == "" {
		return fmt.Errorf("id is required")
	}

	switch c.cfg.Engine {
	case "websocket":
		return c.subscribeWithWebsocket(ctx, id, fn)
	case "http":
		return c.subscribeWithHTTP(ctx, id, fn)
	default:
		return fmt.Errorf("unsupported engine: %s, only support websocket and http", c.cfg.Engine)
	}
}

func (c *client) subscribeWithHTTP(ctx context.Context, id string, fn func(message string)) error {
	response, err := fetch.Stream(fmt.Sprintf("%s/:id/stream", c.cfg.Server), &fetch.Config{
		Context: ctx,
		Params: fetch.Params{
			"id": id,
		},
	})
	if err != nil {
		return err
	}

	// // io.Copy(os.Stdout, response.Stream)
	// buf := make([]byte, 1024)
	// for {
	// 	n, err := response.Stream.Read(buf)
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			break
	// 		}
	// 		return err
	// 	}

	// 	fn(string(buf[:n]))
	// }

	scanner := bufio.NewScanner(response.Stream)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}

		if line == "event: message" {
			continue
		}

		if strings.StartsWith(line, "id:") {
			continue
		}

		if strings.HasPrefix(line, "data:") {
			data := strings.TrimPrefix(line, "data:")
			dataObject := map[string]any{}
			if err := json.Unmarshal([]byte(data), &dataObject); err != nil {
				return err
			}

			if v, ok := dataObject["log"]; ok {
				if vv, ok := v.(string); ok {
					fn(vv)
				}
			}
		}
	}

	return nil
}

func (c *client) subscribeWithWebsocket(ctx context.Context, id string, fn func(message string)) error {

	done := make(chan struct{})
	errCh := make(chan error)
	defer close(done)
	defer close(errCh)

	err := c.event.Emit(
		"subscribe",
		cs.EventPayload{
			"id": id,
		},
		func(err error, payload cs.EventPayload) {
			if err != nil {
				// logger.Infof("failed to subscribe: %s", err)
				errCh <- err
				return
			}

			log, ok := payload.Get("log").(string)
			if ok {
				if log == "[DONE]" {
					done <- struct{}{}
				} else {
					fn(log)
				}
			}
		},
		func(cfg *ec.EmitConfig) {
			cfg.IsSubscribe = true
		},
	)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-errCh:
			return err
		case <-done:
			return nil
		}
	}
}

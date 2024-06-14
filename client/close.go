package client

import "fmt"

func (c *client) Close() error {
	if c.cfg == nil {
		return fmt.Errorf("client is not setup")
	}

	if c.cfg.Engine == "websocket" {
		return c.ws.Close()
	}

	return nil
}

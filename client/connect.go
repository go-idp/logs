package client

import "fmt"

func (c *client) Connect() error {
	if c.cfg == nil {
		return fmt.Errorf("client is not setup")
	}

	if c.cfg.Engine == "websocket" {
		return c.ws.Connect()
	}

	return nil
}

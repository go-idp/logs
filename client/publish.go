package client

import (
	"context"
	"fmt"

	"github.com/go-zoox/fetch"
)

func (c *client) Publish(ctx context.Context, id string, message string) error {
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

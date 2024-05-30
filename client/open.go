package client

import (
	"context"
	"fmt"

	"github.com/go-zoox/fetch"
)

func (c *client) Open(ctx context.Context, id string) error {
	response, err := fetch.Post(fmt.Sprintf("%s/:id/open", c.cfg.Server), &fetch.Config{
		Context: ctx,
		Params: fetch.Params{
			"id": id,
		},
	})
	if err != nil {
		return err
	}

	if response.StatusCode() != 200 {
		return fmt.Errorf("failed to open: %s", response.String())
	}

	if response.Get("code").Int() != 200 {
		return fmt.Errorf("failed to open: %s", response.Get("message").String())
	}

	return nil
}

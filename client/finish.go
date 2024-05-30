package client

import (
	"context"
	"fmt"

	"github.com/go-zoox/fetch"
)

func (c *client) Finish(ctx context.Context, id string) error {
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

package client

import (
	"fmt"

	"github.com/go-zoox/fetch"
)

func (c *client) Finish(id string) error {
	response, err := fetch.Post(fmt.Sprintf("%s/:id/finish", c.cfg.Server), &fetch.Config{
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

package client

import (
	"fmt"
	"io"

	"github.com/go-zoox/fetch"
)

func (c *client) Subscribe(id string, fn func(message string)) error {
	response, err := fetch.Stream(fmt.Sprintf("%s/:id/stream", c.cfg.Server), &fetch.Config{
		Params: fetch.Params{
			"id": id,
		},
	})
	if err != nil {
		return err
	}

	// io.Copy(os.Stdout, response.Stream)
	buf := make([]byte, 1024)
	for {
		n, err := response.Stream.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		fn(string(buf[:n]))
	}

	return nil
}

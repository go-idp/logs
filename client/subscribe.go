package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-zoox/core-utils/strings"
	"github.com/go-zoox/fetch"
)

func (c *client) Subscribe(ctx context.Context, id string, fn func(message string)) error {
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

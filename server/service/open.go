package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-idp/logs/server/pubsub"
	"github.com/go-zoox/datetime"
)

func (s *service) Open(ctx context.Context, id string) error {
	if s.cfg == nil {
		return fmt.Errorf("service is not setup")
	}

	if id == "" {
		return fmt.Errorf("id is required")
	}

	if err := pubsub.Open(ctx, id); err != nil {
		return fmt.Errorf("[open] failed to create log(id: %s): %s", id, err)
	}

	s.manager.Create(id)

	welcomMessage := fmt.Sprintf("[%s][ID: %s] ...\n", datetime.Now().Format("YYYY-MM-DD HH:mm:ss"), id)
	if err := pubsub.Publish(ctx, id, welcomMessage); err != nil {
		return fmt.Errorf("[open] failed to publish log(id: %s): %s", id, err)
	}

	// force clean scheduled
	go func() {
		for {
			select {
			// if context is done, return
			case <-ctx.Done():
				return
			case <-time.After(24 * time.Hour):
				if err := s.Finish(ctx, id); err != nil {
					fmt.Printf("[open] failed to finish log(id: %s): %s\n", id, err)
				}
			}
		}
	}()

	return nil
}

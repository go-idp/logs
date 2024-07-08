package service

import (
	"context"
	"fmt"

	"github.com/go-idp/logs/server/pubsub"
)

func (s *service) Publish(ctx context.Context, id string, message string) error {
	if s.cfg == nil {
		return fmt.Errorf("service is not setup")
	}

	if id == "" {
		return fmt.Errorf("id is required")
	}

	s.manager.Update(id, message)

	if err := pubsub.Publish(ctx, id, message); err != nil {
		return fmt.Errorf("[publish] failed to publish log (id: %s): %s", id, err)
	}

	return nil
}

package service

import (
	"context"
	"fmt"

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
		return fmt.Errorf("failed to create topic: %s", err)
	}

	welcomMessage := fmt.Sprintf("[%s][ID: %s] ...\n", datetime.Now().Format("YYYY-MM-DD HH:mm:ss"), id)
	if err := pubsub.Publish(ctx, id, welcomMessage); err != nil {
		return fmt.Errorf("failed to publish topic: %s", err)
	}

	return nil
}

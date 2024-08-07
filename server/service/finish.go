package service

import (
	"context"
	"fmt"

	"github.com/go-idp/logs/server/pubsub"
	"github.com/go-idp/logs/server/storage/fs"
	"github.com/go-idp/logs/server/storage/oss"
)

func (s *service) Finish(ctx context.Context, id string) error {
	if s.cfg == nil {
		return fmt.Errorf("service is not setup")
	}

	if id == "" {
		return fmt.Errorf("id is required")
	}

	s.manager.Delete(id)

	// get file
	file, err := pubsub.GetFile(id)
	if err != nil {
		return fmt.Errorf("[finish] failed to get file path (id: %s, err: %s)", id, err)
	}
	defer file.Clean()

	if err := pubsub.Close(ctx, id); err != nil {
		return fmt.Errorf("[finish] failed to close pubsub(id: %s): %s", id, err)
	}

	reader, err := file.GetReader()
	if err != nil {
		return fmt.Errorf("[finish] failed to get file reader (id: %s, err: %s)", id, err)
	}
	defer reader.Close()

	switch s.cfg.Storage.Driver {
	case "oss":
		err := oss.Get().Put(id, reader)
		if err != nil {
			return fmt.Errorf("[finish] failed to get file path from oss (id: %s, err: %s)", id, err)
		}
	default:
		err = fs.Get().Put(id, reader)
		if err != nil {
			return fmt.Errorf("[finish] failed to get file path from fs (id: %s, err: %s)", id, err)
		}
	}

	return nil
}

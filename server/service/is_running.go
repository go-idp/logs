package service

import "context"

func (s *service) IsRunning(ctx context.Context, id string) bool {
	if s.cfg == nil {
		return false
	}

	if id == "" {
		return false
	}

	return s.manager.IsRunning(id)
}

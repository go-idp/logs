package service

import (
	"context"

	"github.com/go-idp/logs/server/config"
	"github.com/go-idp/logs/server/service/manager"
	"github.com/go-zoox/once"
)

type Service interface {
	Open(ctx context.Context, id string) error
	Finish(ctx context.Context, id string) error
	//
	Publish(ctx context.Context, id string, message string) error
	Subscribe(ctx context.Context, id string, handler func(err error, message string)) error
	//
	Setup(cfg *config.Config) error
	//
	Status() any
}

type service struct {
	cfg *config.Config

	manager manager.Manager
}

func New(cfg *config.Config) Service {
	return &service{
		cfg:     cfg,
		manager: manager.New(),
	}
}

func Get() Service {
	return once.Get("service", func() Service {
		return New(nil)
	})
}

func (s *service) Setup(cfg *config.Config) error {
	s.cfg = cfg
	return nil
}

func (s *service) Status() any {
	return s.manager.Status()
}

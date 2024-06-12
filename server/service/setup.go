package service

import "github.com/go-idp/logs/server/config"

func (s *service) Setup(cfg *config.Config) error {
	s.cfg = cfg
	return nil
}

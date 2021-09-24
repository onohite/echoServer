package service

import (
	"context"
	"echoTest/internal/config"
	"echoTest/internal/service/db"
	"github.com/labstack/gommon/log"
)

type Service struct {
	DB db.DBService
}

func InitService(ctx context.Context, cfg *config.Config) (*Service, error) {
	var service Service
	if err := service.initDb(ctx, cfg); err != nil {
		return nil, err
	}
	log.Info("All services are up")
	return &service, nil
}

func (s *Service) initDb(ctx context.Context, cfg *config.Config) error {
	dbService, err := db.NewPgxCon(ctx, cfg)
	if err != nil {
		return err
	}
	s.DB = dbService
	log.Info("Database connection complete successful")
	return nil
}

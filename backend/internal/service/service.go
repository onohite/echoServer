package service

import (
	"backend/internal/config"
	"backend/internal/service/db"
	"backend/internal/service/queue"
	"context"
	"github.com/labstack/gommon/log"
)

type Service struct {
	DB    db.DBService
	Queue queue.QueueService
}

func InitService(ctx context.Context, cfg *config.Config) (*Service, error) {
	var service Service
	if err := service.initDb(ctx, cfg); err != nil {
		return nil, err
	}
	if err := service.initQueue(ctx, cfg); err != nil {
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

func (s *Service) initQueue(ctx context.Context, cfg *config.Config) error {
	queueService, err := queue.NewRabbitCon(ctx, cfg)
	if err != nil {
		return err
	}
	s.Queue = queueService
	log.Info("Queue connection complete successful")
	return nil
}

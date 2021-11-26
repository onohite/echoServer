package service

import (
	"backend/internal/config"
	"backend/internal/service/cache"
	"backend/internal/service/db"
	"backend/internal/service/graph"
	"backend/internal/service/queue"
	"context"
	"github.com/labstack/gommon/log"
)

type Service struct {
	DB    db.DBService
	Queue queue.QueueService
	Cache cache.CacheService
	Graph graph.GraphService
}

func InitService(ctx context.Context, cfg *config.Config) (*Service, error) {
	var service Service
	if err := service.initDb(ctx, cfg); err != nil {
		return nil, err
	}
	if err := service.initQueue(ctx, cfg); err != nil {
		return nil, err
	}
	if err := service.initCache(ctx, cfg); err != nil {
		return nil, err
	}
	if err := service.initGraph(ctx, cfg); err != nil {
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

func (s *Service) initCache(ctx context.Context, cfg *config.Config) error {
	cacheService, err := cache.InitCache(ctx, cfg)
	if err != nil {
		return err
	}
	s.Cache = cacheService
	log.Info("Cache connection complete successful")
	return nil
}

func (s *Service) initGraph(ctx context.Context, cfg *config.Config) error {
	graphService, err := graph.NewGraphConn(ctx, cfg)
	if err != nil {
		return err
	}
	s.Graph = graphService
	log.Info("Graph connection complete successful")
	return nil
}

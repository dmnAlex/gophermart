package main

import (
	"context"
	"log"

	"github.com/dmnAlex/gophermart/internal/config"
	"github.com/dmnAlex/gophermart/internal/handler"
	"github.com/dmnAlex/gophermart/internal/logger"
	"github.com/dmnAlex/gophermart/internal/repository"
	"github.com/dmnAlex/gophermart/internal/service"
	"github.com/dmnAlex/gophermart/internal/storage/pg"
)

func main() {
	globalCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	if err = logger.Init(cfg.LogLevel); err != nil {
		log.Fatalf("init logger error: %v", err)
	}

	db, err := pg.New(globalCtx, cfg.DatabaseURI, cfg.MigrationsPath)
	if err != nil {
		log.Fatalf("db error: %v", err)
	}

	repo := repository.NewRepository(db)
	defer repo.Close()

	service := service.NewService(repo)
	handler := handler.NewHandler(service, cfg)
	router := newRouter(handler, cfg)

	if err := router.Run(cfg.RunAddress); err != nil {
		log.Fatalf("router run error: %v", err)
	}
}

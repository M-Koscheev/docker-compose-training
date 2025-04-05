package main

import (
	"context"
	"docker-compose-training/config"
	"docker-compose-training/internal/domain"
	"docker-compose-training/internal/repository"
	"docker-compose-training/internal/rest"
	"docker-compose-training/internal/server"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	_ "docker-compose-training/docs"
)

// @title S3 compatible storage API
// @version 1.0
// @description This is a sample swagger for docker compose training project
// @host localhost:8080
// @BasePath /api/v1
func main() {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		panic(err)
	}

	minioClient, err := server.NewMinioClient("files_bucket", cfg.Minio)
	if err != nil {
		slog.Error("failed to setup minio", "error", err)
		panic(err)
	}

	repositoryVar := repository.NewRepository(minioClient, cfg.Minio.BaseBucket, cfg.Minio.BasePath)
	servicesVar := domain.NewService(repositoryVar)
	handlersVar := rest.NewHandler(servicesVar)

	srv := server.Server{}
	exit := make(chan os.Signal, 1)

	go func() {
		if err = srv.Run(handlersVar.InitRoutes().Handler(), cfg.Server); err != nil {
			slog.Error("error while running server", "error", err)
		}
	}()

	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT)
	<-exit

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("error occurred on server shutting down")
	}

	slog.Error("server stopper")
}

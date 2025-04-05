package domain

import (
	"context"
	"docker-compose-training/internal/repository"
)

type Service struct {
	Storage
}

type Storage interface {
	PostMultipleFiles(ctx context.Context, filesNames []string) error
	PostFile(ctx context.Context, name string) error
	GetFilesList(ctx context.Context) ([]string, error)
	GetFileContent(ctx context.Context, name string) (FileContent, error)
	RemoveFile(ctx context.Context, name string) error
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Storage: NewFileStorageService(repo),
	}
}

package repository

import (
	"context"
	"github.com/minio/minio-go"
)

type Repository struct {
	Storage
}

func NewRepository(minioClient *minio.Client, bucketName, basePath string) *Repository {
	return &Repository{
		Storage: NewFileStorageRepository(minioClient, bucketName, basePath),
	}
}

type Storage interface {
	PostFile(ctx context.Context, fileName string) error
	GetFilesList(ctx context.Context) ([]string, error)
	GetFileContent(ctx context.Context, name string) (string, error)
	RemoveFile(ctx context.Context, name string) error
}

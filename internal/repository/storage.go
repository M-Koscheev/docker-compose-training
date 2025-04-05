package repository

import (
	"context"
	"fmt"
	"github.com/minio/minio-go"
)

type StorageRepository struct {
	minioClient *minio.Client
	bucketName  string
	basePath    string
}

func NewFileStorageRepository(minioClient *minio.Client, bucketName, basePath string) *StorageRepository {
	return &StorageRepository{minioClient: minioClient, bucketName: bucketName, basePath: basePath}
}

func (r *StorageRepository) PostFile(ctx context.Context, fileName string) error {
	_, err := r.minioClient.FPutObjectWithContext(ctx, r.bucketName, fileName, r.basePath+fileName, minio.PutObjectOptions{ContentType: ".txt"})
	if err != nil {
		return fmt.Errorf("failed to exec fput minio function: %w", err)
	}

	return nil
}

func (r *StorageRepository) GetFilesList(ctx context.Context) ([]string, error) {
	doneCh := make(chan struct{})
	defer close(doneCh)

	foundFiles := make([]string, 0)
	foundChan := r.minioClient.ListObjectsV2(r.bucketName, r.basePath, true, doneCh)
	for found := range foundChan {
		select {
		case <-ctx.Done():
			close(doneCh)
			return nil, fmt.Errorf("context closed before end of reading: %w", ctx.Err())
		default:
			foundFiles = append(foundFiles, found.Key)
		}
	}

	return foundFiles, nil
}

func (r *StorageRepository) GetFile(ctx context.Context, name string) (*minio.Object, error) {
	object, err := r.minioClient.GetObjectWithContext(ctx, r.bucketName, r.basePath+name, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from minio: %w", err)
	}

	return object, nil

}

func (r *StorageRepository) RemoveFile(ctx context.Context, name string) error {
	return r.minioClient.RemoveObject(r.bucketName, name)
}

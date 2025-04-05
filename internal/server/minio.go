package server

import (
	"docker-compose-training/config"
	"fmt"
	"github.com/minio/minio-go"
	"log/slog"
)

func NewMinioClient(bucketName string, cfg config.Minio) (*minio.Client, error) {
	endpoint := cfg.Host + ":" + cfg.Port
	accessKeyID := cfg.AccessKey
	secretAccessKey := cfg.SecretAccessKey
	useSSL := *cfg.SSLMode

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize minio client: %w", err)
	}

	location := "ru-ural"

	err = minioClient.MakeBucket(bucketName, location)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(bucketName)
		if errBucketExists == nil && exists {
			slog.Info("we already own bucket %s", bucketName)
		} else {
			return nil, fmt.Errorf("failed to create new bucker: %w", err)
		}
	} else {
		slog.Info("successfully created bucket %s", bucketName)
	}

	return minioClient, nil
}

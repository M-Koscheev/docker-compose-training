package domain

import (
	"context"
	"docker-compose-training/internal/repository"
	"fmt"
	"github.com/minio/minio-go"
	"sync"
)

const (
	errMessage = "file with given name does not exist"
)

var (
	NoObjectErr = minio.ErrInvalidObjectName(errMessage)
)

type StorageService struct {
	repo repository.Storage
}

func NewFileStorageService(repo repository.Storage) *StorageService {
	return &StorageService{repo: repo}
}

func (s *StorageService) PostFile(ctx context.Context, name string) error {
	err := s.repo.PostFile(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to post file with name %s to minio: %w", name, err)
	}

	return nil
}

func (s *StorageService) PostMultipleFiles(ctx context.Context, filesNames []string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errChan := make(chan error, 1)

	var wg sync.WaitGroup
	wg.Add(len(filesNames))
	for _, fileName := range filesNames {
		go func(name string) {
			defer wg.Done()

			err := s.repo.PostFile(ctx, name)
			if err == nil {
				return
			}

			select {
			case errChan <- fmt.Errorf("failed to post file %s: %w", name, err):
			default:
			}

			cancel()
		}(fileName)
	}

	wg.Wait()

	var err error
	select {
	case err = <-errChan:
	default:
	}

	if err != nil {
		return fmt.Errorf("failed to post files to minio: %w", err)
	}

	return nil
}

func (s *StorageService) GetFilesList(ctx context.Context) ([]string, error) {
	files, err := s.repo.GetFilesList(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get files list from minio: %w", err)
	}

	return files, nil
}

func (s *StorageService) GetFileContent(ctx context.Context, name string) (FileContent, error) {
	names, err := s.GetFilesList(ctx)
	if err != nil {
		return FileContent{}, fmt.Errorf("failed to get files list: %w", err)
	}

	found := false
	for _, curName := range names {
		if curName == name {
			found = true
			break
		}
	}
	if !found {
		return FileContent{}, NoObjectErr
	}

	content, err := s.repo.GetFileContent(ctx, name)
	if err != nil {
		return FileContent{}, fmt.Errorf("failed to get file %s content: %w", name, err)
	}

	return FileContent{Content: content}, nil
}

func (s *StorageService) RemoveFile(ctx context.Context, name string) error {
	names, err := s.GetFilesList(ctx)
	if err != nil {
		return fmt.Errorf("failed to get files list: %w", err)
	}

	found := false
	for _, curName := range names {
		if curName == name {
			found = true
			break
		}
	}
	if !found {
		return NoObjectErr
	}

	err = s.repo.RemoveFile(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to remove file from minio: %w", err)
	}

	return nil
}

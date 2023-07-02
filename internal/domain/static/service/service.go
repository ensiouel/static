package service

import (
	"context"
	"github.com/ensiouel/static/internal/domain/static"
	"golang.org/x/exp/slog"
	"os"
	filepathpkg "path/filepath"
)

type StaticService struct {
	storage static.Storage
	logger  *slog.Logger
}

func NewStaticService(logger *slog.Logger, storage static.Storage) *StaticService {
	return &StaticService{logger: logger, storage: storage}
}

func (service *StaticService) Upload(ctx context.Context, file *os.File, hashsum string) (string, error) {
	// hash-based file organization
	filedir := filepathpkg.Join(hashsum[0:2], hashsum[2:4])

	filename := hashsum

	filepath := filepathpkg.Join(filedir, filename)

	_, err := os.Stat(filepath)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}

		err = service.storage.MoveTemp(ctx, file.Name(), filepath)
		if err != nil {
			return "", err
		}
	}

	return filename, nil
}

func (service *StaticService) Download(ctx context.Context, filename string) (*os.File, error) {
	// hash-based file organization
	filedir := filepathpkg.Join(filename[0:2], filename[2:4])

	filepath := filepathpkg.Join(filedir, filename)

	file, err := service.storage.Open(ctx, filepath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (service *StaticService) Delete(ctx context.Context, filename string) error {
	// hash-based file organization
	filedir := filepathpkg.Join(filename[0:2], filename[2:4])

	filepath := filepathpkg.Join(filedir, filename)

	err := service.storage.Delete(ctx, filepath)
	if err != nil {
		return err
	}

	return nil
}

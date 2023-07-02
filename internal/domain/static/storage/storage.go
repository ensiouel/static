package storage

import (
	"context"
	"os"
	filepathpkg "path/filepath"
)

type StaticStorage struct {
	root string
}

func NewStaticStorage(root string) (*StaticStorage, error) {
	_, err := os.Stat(root)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		err = os.MkdirAll(root, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	return &StaticStorage{root: root}, nil
}

func (storage *StaticStorage) MoveTemp(_ context.Context, oldpath string, newpath string) error {
	newpath = filepathpkg.Join(storage.root, newpath)

	filedir := filepathpkg.Dir(newpath)

	_, err := os.Stat(filedir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		err = os.MkdirAll(filedir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	err = os.Rename(oldpath, newpath)
	if err != nil {
		return err
	}

	return nil
}

func (storage *StaticStorage) Open(_ context.Context, filepath string) (*os.File, error) {
	filepath = filepathpkg.Join(storage.root, filepath)

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (storage *StaticStorage) Delete(_ context.Context, filepath string) error {
	filepath = filepathpkg.Join(storage.root, filepath)

	err := os.Remove(filepath)
	if err != nil {
		return err
	}

	return nil
}

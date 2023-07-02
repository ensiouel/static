package static

import (
	"context"
	"os"
)

type Storage interface {
	MoveTemp(ctx context.Context, temppath string, filepath string) error
	Open(ctx context.Context, filepath string) (*os.File, error)
	Delete(ctx context.Context, filepath string) error
}

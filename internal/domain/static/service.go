package static

import (
	"context"
	"os"
)

type Service interface {
	Upload(ctx context.Context, file *os.File, hashsum string) (string, error)
	Download(ctx context.Context, filename string) (*os.File, error)
	Delete(ctx context.Context, filename string) error
}

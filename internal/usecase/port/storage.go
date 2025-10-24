package port

import (
	"context"
	"io"
)

type FileStorage interface {
	Save(ctx context.Context, path string, data io.Reader, size int64) error
	Get(ctx context.Context, path string) (io.ReadCloser, error)
	Delete(ctx context.Context, path string) error
}

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

type KeyValueStorage interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, expireSeconds int) error
	Del(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}

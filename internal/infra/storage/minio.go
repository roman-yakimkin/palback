package storage

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type MinioStorage struct {
	client     *minio.Client
	bucketName string
}

func NewMinioStorage(client *minio.Client, bucketName string) *MinioStorage {
	return &MinioStorage{
		client:     client,
		bucketName: bucketName,
	}
}

func (s *MinioStorage) Save(ctx context.Context, path string, data io.Reader, size int64) error {
	_, err := s.client.PutObject(ctx, s.bucketName, path, data, size, minio.PutObjectOptions{})
	return err
}

func (s *MinioStorage) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	obj, err := s.client.GetObject(ctx, s.bucketName, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	_, err = obj.Stat()
	if err != nil {
		obj.Close()
		return nil, err
	}

	return obj, nil
}

func (s *MinioStorage) Delete(ctx context.Context, path string) error {
	return s.client.RemoveObject(ctx, s.bucketName, path, minio.RemoveObjectOptions{})
}

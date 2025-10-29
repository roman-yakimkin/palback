package storage

import (
	"context"
	"errors"

	"github.com/gomodule/redigo/redis"

	"palback/internal/usecase"
)

type RedisStorage struct {
	redisPool *redis.Pool
}

func NewRedisStorage(redisPool *redis.Pool) *RedisStorage {
	return &RedisStorage{
		redisPool: redisPool,
	}
}

func (s *RedisStorage) Set(ctx context.Context, key, value string, expireSeconds int) error {
	conn := s.redisPool.Get()
	defer conn.Close()

	if expireSeconds > 0 {
		_, err := redis.DoContext(conn, ctx, "SET", key, value, "EX", expireSeconds)
		return err
	}

	_, err := redis.DoContext(conn, ctx, "SET", key, value)

	return err
}

func (s *RedisStorage) Get(ctx context.Context, key string) (string, error) {
	conn := s.redisPool.Get()
	defer conn.Close()

	reply, err := redis.DoContext(conn, ctx, "GET", key)

	if err != nil {
		switch {
		case errors.Is(err, redis.ErrNil):
			return "", usecase.ErrKeyNotFound
		default:
			return "", err
		}
	}

	if reply == nil {
		return "", usecase.ErrNoReplyFromKeyValueStorage
	}

	return redis.String(reply, err)
}

func (s *RedisStorage) Del(ctx context.Context, key string) error {
	conn := s.redisPool.Get()
	defer conn.Close()

	_, err := redis.DoContext(conn, ctx, "DEL", key)

	return err
}

func (s *RedisStorage) Exists(ctx context.Context, key string) (bool, error) {
	conn := s.redisPool.Get()
	defer conn.Close()

	return redis.Bool(redis.DoContext(conn, ctx, "EXISTS", key))
}

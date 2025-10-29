package rate

import (
	"context"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

type RedigoRateLimiter struct {
	pool *redis.Pool
}

func NewRedigoRateLimiter(pool *redis.Pool) *RedigoRateLimiter {
	return &RedigoRateLimiter{pool: pool}
}

// Allow возвращает true, если запрос разрешён.
// key — уникальный идентификатор (например, "login:192.168.1.1")
// limit — максимальное число запросов
// windowSeconds — размер временного окна в секундах
func (r *RedigoRateLimiter) Allow(ctx context.Context, key string, limit int, windowSeconds int) (bool, error) {
	conn := r.pool.Get()
	defer conn.Close()

	// Атомарно увеличиваем счётчик
	newCount, err := redis.Int(redis.DoContext(conn, ctx, "INCR", key))
	if err != nil {
		return false, fmt.Errorf("redis INCR failed: %w", err)
	}

	// Если это первый запрос в окне — устанавливаем TTL
	if newCount == 1 {
		_, err := redis.DoContext(conn, ctx, "EXPIRE", key, windowSeconds)
		if err != nil {
			return false, fmt.Errorf("redis EXPIRE failed: %w", err)
		}
	}

	// Проверяем лимит
	return newCount <= limit, nil
}

package redis

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	*redis.Client
}

func New(ctx context.Context, log *slog.Logger, addr string, password string, db int) (*Redis, error) {
	rdb := redis.NewClient(
		&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		},
	)

	// Ping the Redis server to check the connection
	res, err := rdb.Ping(ctx).Result()
	log.Info("Redis - New - Ping: ", res)
	if err != nil {
		return nil, fmt.Errorf("redis - New - Ping: %w", err)
	}

	return &Redis{rdb}, nil
}

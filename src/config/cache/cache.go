package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/suricat89/rinha-2024-q1/src/config"
)

var (
	Rdb *redis.Client
)

func InitCache() {
	cfg := config.Env.Cache

	Rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	})
}

func PingRedis() error {
	result := Rdb.Ping(context.Background())
	return result.Err()
}

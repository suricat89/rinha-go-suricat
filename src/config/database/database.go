package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suricat89/rinha-2024-q1/src/config"
)

var (
	DBPool *pgxpool.Pool
)

func InitDb() error {
	cfg := config.Env.Database
	var err error

  connStr := fmt.Sprintf(
		"dbname=%s user=%s password=%s host=%s port=%d application_name=%s sslmode=%s timezone=%s pool_min_conns=%d pool_min_conns=%d",
		cfg.DB,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.AppName,
		cfg.SSLMode,
		cfg.Timezone,
    cfg.MaxConnections,
    cfg.MaxConnections,
	)

  options, err := pgxpool.ParseConfig(connStr)
  if err != nil {
    return err
  }

	DBPool, err = pgxpool.NewWithConfig(context.Background(), options)
	if err != nil {
		return err
	}

	return nil
}

func PingDB() error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(config.Env.Database.CommandTimeout)*time.Second,
	)
	defer cancel()
	return DBPool.Ping(ctx)
}

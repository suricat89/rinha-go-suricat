package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suricat89/rinha-2024-q1/src/config"
	"github.com/suricat89/rinha-2024-q1/src/interfaces"
)

type PostgresDatabase struct {
	DBPool *pgxpool.Pool
}

func NewPostgresDatabase() interfaces.DatabaseConfig {
	return &PostgresDatabase{}
}

func (d *PostgresDatabase) InitDb() (interface{}, error) {
	cfg := config.Env.Database
	var err error

	connStr := fmt.Sprintf(
		"dbname=%s user=%s password=%s host=%s port=%d application_name=%s sslmode=%s timezone=%s pool_min_conns=%d pool_max_conns=%d",
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
		return nil, err
	}

	d.DBPool, err = pgxpool.NewWithConfig(context.Background(), options)
	if err != nil {
		return nil, err
	}

	return d.DBPool, nil
}

func (d *PostgresDatabase) Close() error {
	d.DBPool.Close()
	return nil
}

func (d *PostgresDatabase) PingDb() error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(config.Env.Database.CommandTimeout)*time.Second,
	)
	defer cancel()
	return d.DBPool.Ping(ctx)
}

package database

import (
	"context"
	"time"

	"github.com/suricat89/rinha-2024-q1/src/config"
	"github.com/suricat89/rinha-2024-q1/src/interfaces"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatabase struct {
	mongoClient *mongo.Client
	Db          *mongo.Database
}

func NewMongoDatabase() interfaces.DatabaseConfig {
	return &MongoDatabase{}
}

func (d *MongoDatabase) InitDb() (interface{}, error) {
	cfg := config.Env.Database
	var err error

	opts := options.Client()
	opts.ApplyURI(cfg.Uri)
	opts.SetConnectTimeout(time.Duration(cfg.ConnectionTimeout) * time.Second)
	opts.SetMinPoolSize(uint64(cfg.MaxConnections))
	opts.SetMaxPoolSize(uint64(cfg.MaxConnections))
	opts.SetAppName(cfg.AppName)
	opts.SetTimeout(time.Duration(cfg.CommandTimeout) * time.Second)

	d.mongoClient, err = mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, err
	}

	d.Db = d.mongoClient.Database(cfg.DB)

	return d.Db, nil
}

func (d *MongoDatabase) Close() error {
	return d.mongoClient.Disconnect(context.Background())
}

func (d *MongoDatabase) PingDb() error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(config.Env.Database.CommandTimeout)*time.Second,
	)
	defer cancel()
	return d.mongoClient.Ping(ctx, nil)
}

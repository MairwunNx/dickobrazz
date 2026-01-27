package storage

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo(ctx context.Context, log *slog.Logger, uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI(uri).
		SetAppName("DickobrazzServer").
		SetTimeout(10*time.Second),
	)
	if err != nil {
		return nil, err
	}

	ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := client.Ping(ctxPing, nil); err != nil {
		return nil, err
	}

	log.Info("MongoDB connection established")
	return client, nil
}

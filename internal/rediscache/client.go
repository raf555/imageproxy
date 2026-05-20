package rediscache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewSingleClient(uri string) (*redis.Client, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	redisOpt, err := redis.ParseURL(uri)
	if err != nil {
		return nil, fmt.Errorf("parse uri: %w", err)
	}

	client := redis.NewClient(redisOpt)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return client, nil
}

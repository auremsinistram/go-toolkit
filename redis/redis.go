package redis

import (
	"context"
	"time"

	"github.com/auremsinistram/go-errors"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func New() *Redis {
	return &Redis{}
}

func (r *Redis) Connect(
	addr string,
	username string,
	password string,
	db int,
) error {
	r.Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err := r.Client.Ping(ctx).Err(); err != nil {
		return errors.Wrap(err, "Redis - Connect - #1")
	}

	return nil
}

func (r *Redis) Close() error {
	if err := r.Client.Close(); err != nil {
		return errors.Wrap(err, "Redis - Close - #1")
	}

	return nil
}

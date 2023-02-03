package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Ctx      context.Context
	Addr     string
	Password string
	rdb      *redis.Client
	DB       int
}

func (r *Redis) Init() {
	r.rdb = redis.NewClient(&redis.Options{
		Addr:     r.Addr,
		Password: r.Password,
		DB:       r.DB,
	})
}

func (r *Redis) GetDB() *redis.Client {
	return r.rdb
}

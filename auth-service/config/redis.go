package config

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var Rdb *redis.Client

func ConnectRedis() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	Rdb = redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if _, err := Rdb.Ping(Ctx).Result(); err != nil {
		// log instead of panic if redis not available for now
		// but as per instructions "panic"
		panic("Redis failed")
	}
}

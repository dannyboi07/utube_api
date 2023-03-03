package redis

import (
	"context"
	"os"
	"strconv"

	redisPkg "github.com/redis/go-redis/v9"
)

var redisClient *redisPkg.Client
var ctx context.Context

func Init() error {
	ctx = context.Background()
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return err
	}

	redisClient = redisPkg.NewClient(&redisPkg.Options{
		Network:     "tcp",
		Addr:        os.Getenv("REDIS_ADDR"),
		Password:    os.Getenv("REDIS_PWD"),
		DB:          db,
		ReadTimeout: 0,
	})

	return nil
}

func Close() error {
	return redisClient.Close()
}

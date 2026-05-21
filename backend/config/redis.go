package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := 1; i <= 10; i++ {
		if err := rdb.Ping(ctx).Err(); err == nil {
			break
		} else if i == 10 {
			log.Fatalf("[Redis] could not connect after 10 attempts: %v", err)
		} else {
			log.Printf("[Redis] attempt %d/10 failed, retrying in 2s...", i)
			time.Sleep(2 * time.Second)
		}
	}

	log.Println("[Redis] connected successfully")
	return rdb
}

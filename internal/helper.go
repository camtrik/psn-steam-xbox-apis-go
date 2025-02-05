package helper

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/camtrik/ebbilogue-backend/internal/global"
	"github.com/redis/go-redis/v9"
)

func InitRedis(cfg *global.Config) *redis.Client {
	redisURL := cfg.RedisUrl
	if !strings.HasPrefix(redisURL, "redis://") {
		redisURL = fmt.Sprintf("redis://%s", redisURL)
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}

	rdb := redis.NewClient(opt)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	return rdb
}

func InitHttpClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 10,
	}
}

package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/camtrik/ebbilogue-backend/internal/global"
	"github.com/camtrik/ebbilogue-backend/internal/models"
	"github.com/redis/go-redis/v9"
)

type RedisPSNCache struct {
	client *redis.Client
}

func NewRedisPSNCache(client *redis.Client) *RedisPSNCache {
	return &RedisPSNCache{
		client: client,
	}
}

func (c *RedisPSNCache) GetUserTitles(ctx context.Context, accountId string) (*models.UserTitlesResponse, error) {
	key := fmt.Sprintf(global.USER_TITLES_KEY, accountId)

	// try fetch from cache
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user titles from cache: %v", err)
	}

	var titles models.UserTitlesResponse
	if err := json.Unmarshal(data, &titles); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user titles: %v", err)
	}
	return &titles, nil
}

func (c *RedisPSNCache) SetUserTitles(ctx context.Context, accountId string, titles *models.UserTitlesResponse) error {
	key := fmt.Sprintf(global.USER_TITLES_KEY, accountId)

	data, err := json.Marshal(titles)
	if err != nil {
		return fmt.Errorf("failed to marshal user titles: %v", err)
	}

	return c.client.Set(ctx, key, data, global.DEFAULT_EXPIRATION).Err()
}

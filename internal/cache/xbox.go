package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/camtrik/psn-steam-api/internal/global"
	"github.com/camtrik/psn-steam-api/internal/models"
	"github.com/redis/go-redis/v9"
)

type RedisXboxCache struct {
	client *redis.Client
}

func NewRedisXboxCache(client *redis.Client) *RedisXboxCache {
	return &RedisXboxCache{
		client: client,
	}
}

func (c *RedisXboxCache) GetPlayerAchievements(ctx context.Context) (*models.XboxGamaAchievements, error) {
	key := fmt.Sprintf(global.XBOX_PLAYER_ACHIEVEMENTS_KEY)

	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get player achievements from cache: %v", err)
	}

	var playerAchievements models.XboxGamaAchievements
	if err := json.Unmarshal(data, &playerAchievements); err != nil {
		return nil, fmt.Errorf("failed to unmarshal player achievements: %v", err)
	}

	return &playerAchievements, nil
}

func (c *RedisXboxCache) SetPlayerAchievements(ctx context.Context, playerAchievements *models.XboxGamaAchievements) error {
	key := fmt.Sprintf(global.XBOX_PLAYER_ACHIEVEMENTS_KEY)
	data, err := json.Marshal(playerAchievements)
	if err != nil {
		return fmt.Errorf("failed to marshal player achievements: %v", err)
	}

	return c.client.Set(ctx, key, data, global.DEFAULT_EXPIRATION).Err()
}

func (c *RedisXboxCache) GetGameStats(ctx context.Context, titleId string) (*models.XboxGameStats, error) {
	key := fmt.Sprintf(global.XBOX_PLAYER_GAME_STATS_KEY, titleId)

	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}

	var gameStats models.XboxGameStats
	if err := json.Unmarshal(data, &gameStats); err != nil {
		return nil, fmt.Errorf("failed to unmarshal game stats: %v", err)
	}

	return &gameStats, nil
}

func (c *RedisXboxCache) SetGameStats(ctx context.Context, titleId string, gameStats *models.XboxGameStats) error {
	key := fmt.Sprintf(global.XBOX_PLAYER_GAME_STATS_KEY, titleId)
	data, err := json.Marshal(gameStats)
	if err != nil {
		return fmt.Errorf("failed to marshal game stats: %v", err)
	}

	return c.client.Set(ctx, key, data, global.DEFAULT_EXPIRATION).Err()
}

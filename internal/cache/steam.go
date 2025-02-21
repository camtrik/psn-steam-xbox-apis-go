package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/camtrik/psn-steam-api/internal/global"
	"github.com/camtrik/psn-steam-api/internal/models"
	"github.com/redis/go-redis/v9"
)

type RedisSteamCache struct {
	client *redis.Client
}

func NewRedisSteamCache(client *redis.Client) *RedisSteamCache {
	return &RedisSteamCache{
		client: client,
	}
}

func (c *RedisSteamCache) GetOwnedGames(ctx context.Context, steamId string) (*models.OwnedGamesResponse, error) {
	key := fmt.Sprintf(global.OWNED_GAMES_KEY, steamId)

	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get owned games from redis: %v", err)
	}

	var ownedGames models.OwnedGamesResponse
	if err := json.Unmarshal(data, &ownedGames); err != nil {
		return nil, fmt.Errorf("failed to unmarshal owned games: %v", err)
	}

	return &ownedGames, nil

}

func (c *RedisSteamCache) SetOwnedGames(ctx context.Context, steamId string, ownedGames *models.OwnedGamesResponse) error {
	key := fmt.Sprintf(global.OWNED_GAMES_KEY, steamId)
	data, err := json.Marshal(ownedGames)
	if err != nil {
		return fmt.Errorf("failed to marshal owned games: %v", err)
	}

	return c.client.Set(ctx, key, data, global.DEFAULT_EXPIRATION).Err()
}

func (c *RedisSteamCache) GetPlayerAchievements(ctx context.Context, steamId string, appId int) (*models.PlayerAchievementsResponse, error) {
	key := fmt.Sprintf(global.PLAYER_ACHIEVEMENTS_KEY, steamId, appId)
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get player achievements from redis: %v", err)
	}
	var playerAchievements models.PlayerAchievementsResponse
	if err := json.Unmarshal(data, &playerAchievements); err != nil {
		return nil, fmt.Errorf("failed to unmarshal player achievements: %v", err)
	}
	return &playerAchievements, nil
}

func (c *RedisSteamCache) SetPlayerAchievements(ctx context.Context, steamId string, appId int, playerAchievements *models.PlayerAchievementsResponse) error {
	key := fmt.Sprintf(global.PLAYER_ACHIEVEMENTS_KEY, steamId, appId)
	data, err := json.Marshal(playerAchievements)
	if err != nil {
		return fmt.Errorf("failed to marshal player achievements: %v", err)
	}
	return c.client.Set(ctx, key, data, global.DEFAULT_EXPIRATION).Err()
}

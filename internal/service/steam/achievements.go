package steam

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/camtrik/ebbilogue-backend/internal/global"
	"github.com/camtrik/ebbilogue-backend/internal/models"
)

func (s *SteamService) GetPlayerAchievements(ctx context.Context, steamId string, appId int) (*models.PlayerAchievementsResponse, error) {
	// cache
	if cached, err := s.cache.GetPlayerAchievements(ctx, steamId, appId); err == nil && cached != nil {
		return cached, nil
	}

	url := fmt.Sprintf("%s/ISteamUserStats/GetPlayerAchievements/v1/?appid=%d&key=%s&steamid=%s", global.STEAM_API_BASE_URL, appId, s.apiKey, steamId)
	resp, err := s.client.Get(url)
	if err != nil {
		s.logger.Error("Failed to fetch steam player achievements for %s %v", steamId, err)
		return nil, err
	}

	var playerAchievementsResponse models.PlayerAchievementsResponse
	if err := json.NewDecoder(resp.Body).Decode(&playerAchievementsResponse); err != nil {
		s.logger.Error("Failed to decode steam player achievements for %s %v", steamId, err)
		return nil, err
	}
	// set cache
	if err := s.cache.SetPlayerAchievements(ctx, steamId, appId, &playerAchievementsResponse); err != nil {
		s.logger.Error("Failed to set steam player achievements cache for %s %v", steamId, err)
	} else {
		s.logger.Info("Set steam player achievements cache for %s", steamId)
	}

	return &playerAchievementsResponse, nil

}

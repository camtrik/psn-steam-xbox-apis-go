package steam

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/camtrik/psn-steam-api/internal/models"
)

func (s *SteamService) GetOwnedGames(ctx context.Context, steamId string) (*models.OwnedGamesResponse, error) {
	// fetch from cache
	if cached, err := s.cache.GetOwnedGames(ctx, steamId); err == nil && cached != nil {
		s.logger.Debug("Cache hit for steam owned games for %s", steamId)
		return cached, nil
	}

	// fetch
	url := fmt.Sprintf("%s/IPlayerService/GetOwnedGames/v0001/?key=%s&steamid=%s&include_appinfo=true&include_played_free_games=true&format=json", ApiBaseURL, s.apiKey, steamId)

	resp, err := s.client.Get(url)
	if err != nil {
		s.logger.Error("Failed to fetch steam owned games for %s %v", steamId, err)
		return nil, err
	}

	var ownedGamesResp models.OwnedGamesResponse
	if err := json.NewDecoder(resp.Body).Decode(&ownedGamesResp); err != nil {
		s.logger.Error("Failed to decode steam owned games for %s %v", steamId, err)
		return nil, err
	}

	// set cache
	if err := s.cache.SetOwnedGames(ctx, steamId, &ownedGamesResp); err != nil {
		s.logger.Error("Failed to set cache for steam owned games for %s %v", steamId, err)
	} else {
		s.logger.Debugf("Cached Owned Games for user %s", steamId)
	}

	return &ownedGamesResp, nil
}

package xbox

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/camtrik/psn-steam-api/internal/models"
)

func (s *XboxService) GetPlayerAchievements(ctx context.Context) (*models.XboxGamaAchievements, error) {
	// cahce
	if cached, err := s.cache.GetPlayerAchievements(ctx); err == nil && cached != nil {
		return cached, nil
	}

	url := fmt.Sprintf("%s/achievements", ApiBaseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		s.logger.Error("Failed to create request for xbox player achievements %v", err)
		return nil, err
	}

	req.Header.Set("accept", "*/*")
	req.Header.Set("x-authorization", s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Error("Failed to send request for xbox player achievements %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var achievements models.XboxGamaAchievements
	if err := json.NewDecoder(resp.Body).Decode(&achievements); err != nil {
		s.logger.Error("Failed to decode xbox player achievements %v", err)
		return nil, err
	}

	// set cache
	if err := s.cache.SetPlayerAchievements(ctx, &achievements); err != nil {
		s.logger.Error("Failed to set xbox player achievements to cache %v", err)
	} else {
		s.logger.Info("Set xbox player achievements to cache")
	}

	return &achievements, nil
}

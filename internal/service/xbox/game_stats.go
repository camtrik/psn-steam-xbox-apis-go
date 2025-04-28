package xbox

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/camtrik/psn-steam-api/internal/models"
)

func (s *XboxService) GetGameStats(ctx context.Context, titleId string) (models.XboxGameStats, error) {
	// cache
	if cached, err := s.cache.GetGameStats(ctx, titleId); err == nil && cached != nil {
		return *cached, nil
	}

	url := fmt.Sprintf("%s/achievements/stats/%s", ApiBaseURL, titleId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		s.logger.Error("Failed to create request for xbox game stats %v", err)
		return models.XboxGameStats{}, err
	}

	req.Header.Set("accept", "*/*")
	req.Header.Set("x-authorization", s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Error("Failed to send request for xbox game stats %v", err)
		return models.XboxGameStats{}, err
	}
	defer resp.Body.Close()

	var stats models.XboxGameStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		s.logger.Error("Failed to decode xbox game stats %v", err)
		return models.XboxGameStats{}, err
	}

	// set cache
	if err := s.cache.SetGameStats(ctx, titleId, &stats); err != nil {
		s.logger.Error("Failed to set xbox game stats to cache %v", err)
	} else {
		s.logger.Info("Set xbox game stats to cache")
	}

	return stats, nil
}

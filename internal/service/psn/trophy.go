// internal/service/psn/trophy.go
package psn

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/camtrik/psn-steam-api/internal/models"
)

func applyPagination(titles []models.TrophyTitle, totalItems int, options *models.GetUserTitlesOptions) []models.TrophyTitle {
	if options == nil {
		return titles
	}

	start := 0
	if options.Offset != nil {
		start = *options.Offset
	}

	end := totalItems
	if options.Limit != nil {
		end = start + *options.Limit
		if end > totalItems {
			end = totalItems
		}
	}

	if start > totalItems {
		start = totalItems
	}

	return titles[start:end]
}

func applyFilter(titles []models.TrophyTitle, filter models.TrophyFilter) []models.TrophyTitle {
	fmt.Println("Checking filter.... ")
	if (filter.MinProgress == 0) && (filter.Platform == "") && (filter.SortBy == "") {
		return titles
	}

	filteredTrophyTitles := make([]models.TrophyTitle, 0)

	for _, trophyTitle := range titles {
		if trophyTitle.Progress >= filter.MinProgress &&
			(filter.Platform == "" || trophyTitle.TrophyTitlePlatform == filter.Platform) {
			filteredTrophyTitles = append(filteredTrophyTitles, trophyTitle)
		}
	}

	switch filter.SortBy {
	case "lastUpdated":
		sort.Slice(filteredTrophyTitles, func(i, j int) bool {
			return filteredTrophyTitles[i].LastUpdatedDateTime.Unix() > filteredTrophyTitles[j].LastUpdatedDateTime.Unix()
		})
	case "progress":
		sort.Slice(filteredTrophyTitles, func(i, j int) bool {
			return filteredTrophyTitles[i].Progress > filteredTrophyTitles[j].Progress
		})
	}

	return filteredTrophyTitles
}

func (s *PSNService) GetUserTitles(ctx context.Context, accountId string, options *models.GetUserTitlesOptions, filter models.TrophyFilter) (*models.UserTitlesResponse, error) {
	// fetch from cache
	startTime := time.Now()
	s.logger.Infof("Fetching trophies for user %s with options: %+v", accountId, options)

	if cached, err := s.cache.GetUserTitles(ctx, accountId); err == nil && cached != nil {
		s.logger.Debugf("Cache hit for user %s", accountId)
		filteredTitles := applyFilter(cached.TrophyTitles, filter)
		pagedTitles := applyPagination(filteredTitles, len(filteredTitles), options)
		return &models.UserTitlesResponse{
			TrophyTitles:   pagedTitles,
			TotalItemCount: len(pagedTitles),
		}, nil
	}

	s.logger.Debug("Cache miss for user %s", accountId)
	// fetch from API
	titles, err := s.fetchUserTitles(ctx, accountId)
	if err != nil {
		s.logger.Errorf("Error fetching trophies for user %s: %v", accountId, err)
		return nil, err
	}

	// set cache
	if err := s.cache.SetUserTitles(ctx, accountId, titles); err != nil {
		s.logger.Errorf("Failed to cache user titles: %v", err)
	} else {
		s.logger.Debugf("Cached trophies for user %s", accountId)
	}

	// apply pagaination and return
	filteredTitles := applyFilter(titles.TrophyTitles, filter)
	pagedTitles := applyPagination(filteredTitles, len(filteredTitles), options)
	s.logger.Infof("Retrieved %d trophies from cache for user %s (took %v)",
		len(pagedTitles), accountId, time.Since(startTime))

	return &models.UserTitlesResponse{
		TrophyTitles:   pagedTitles,
		TotalItemCount: len(filteredTitles),
	}, nil

}

func (s *PSNService) fetchUserTitles(ctx context.Context, accountId string) (*models.UserTitlesResponse, error) {
	// get valid token by refresh token
	auth, err := s.GetValidAuthorization()
	fmt.Println("accesstoken: ", auth.AccessToken)
	if err != nil {
		return nil, err
	}

	baseURL := fmt.Sprintf("%s/v1/users/%s/trophyTitles", TrophyBaseURL, accountId)

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", auth.AccessToken))
	req.Header.Set("Accept", "application/json")

	req = req.WithContext(ctx)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user titles: %v", resp.Status)
	}

	var userTitlesResp models.UserTitlesResponse
	if err := json.NewDecoder(resp.Body).Decode(&userTitlesResp); err != nil {
		return nil, fmt.Errorf("failed to decode user titles response: %v", err)
	}

	return &userTitlesResp, nil
}

package psn

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
)

type GetUserTitlesOptions struct {
	Limit  *int
	Offset *int
}

type UserTitlesResponse struct {
	TrophyTitles   []TrophyTitle `json:"trophyTitles"`
	TotalItemCount int           `json:"totalItemCount"`
}

type TrophyTitle struct {
	NpCommunicationId   string          `json:"npCommunicationId"`
	TrophySetVersion    string          `json:"trophySetVersion"`
	TrophyTitleName     string          `json:"trophyTitleName"`
	TrophyTitlePlatform string          `json:"trophyTitlePlatform"`
	TrophyTitleIconUrl  string          `json:"trophyTitleIconUrl"`
	HasEarnedTrophies   bool            `json:"hasEarnedTrophies"`
	DefinedTrophies     DefinedTrophies `json:"definedTrophies"`
	Progress            int             `json:"progress"`
	EarnedTrophies      EarnedTrophies  `json:"earnedTrophies"`
	LastUpdatedDateTime string          `json:"lastUpdatedDateTime"`
}

type DefinedTrophies struct {
	Bronze   int `json:"bronze"`
	Silver   int `json:"silver"`
	Gold     int `json:"gold"`
	Platinum int `json:"platinum"`
}

type EarnedTrophies struct {
	Bronze   int `json:"bronze"`
	Silver   int `json:"silver"`
	Gold     int `json:"gold"`
	Platinum int `json:"platinum"`
}

type TrophyFilter struct {
	MinProgress int
	Platform    string
	SortBy      string
}

func (s *PSNService) GetMyTitles(option *GetUserTitlesOptions) (*UserTitlesResponse, error) {
	return s.GetUserTitles("me", option)
}

func (s *PSNService) GetMyFilteredTitles(option *GetUserTitlesOptions, filter TrophyFilter) (*UserTitlesResponse, error) {
	return s.GetFilteredUserTitles("me", option, filter)
}

func (s *PSNService) GetUserTitles(accountId string, options *GetUserTitlesOptions) (*UserTitlesResponse, error) {
	auth, err := s.GetValidAuthorization()
	if err != nil {
		return nil, err
	}

	baseURL := fmt.Sprintf("%s/v1/users/%s/trophyTitles", TrophyBaseURL, accountId)

	query := url.Values{}
	if options != nil {
		if options.Limit != nil {
			if *options.Limit < 1 || *options.Limit > 800 {
				return nil, fmt.Errorf("limit must be between 1 and 800")
			}
			query.Add("limit", strconv.Itoa(*options.Limit))
		}
		if options.Offset != nil {
			if *options.Offset < 0 {
				return nil, fmt.Errorf("offset must be greater than or equal to 0")
			}
			query.Add("offset", strconv.Itoa(*options.Offset))
		}
	}

	if len(query) > 0 {
		baseURL += "?" + query.Encode()
	}

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", auth.AccessToken))
	req.Header.Set("Accept", "application/json")

	// send request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user titles: %v", resp.Status)
	}

	var userTitlesResp UserTitlesResponse
	if err := json.NewDecoder(resp.Body).Decode(&userTitlesResp); err != nil {
		return nil, fmt.Errorf("failed to decode user titles: %v", err)
	}

	return &userTitlesResp, nil
}

func (s *PSNService) GetFilteredUserTitles(accountId string, option *GetUserTitlesOptions, filter TrophyFilter) (*UserTitlesResponse, error) {
	resp, err := s.GetUserTitles(accountId, option)
	if err != nil {
		return nil, err
	}

	filteredTrophyTitles := make([]TrophyTitle, 0)

	for _, trophyTitle := range resp.TrophyTitles {
		if trophyTitle.Progress >= filter.MinProgress &&
			(filter.Platform == "" || trophyTitle.TrophyTitlePlatform == filter.Platform) {
			filteredTrophyTitles = append(filteredTrophyTitles, trophyTitle)
		}
	}

	switch filter.SortBy {
	case "lastUpdated":
		sort.Slice(filteredTrophyTitles, func(i, j int) bool {
			return filteredTrophyTitles[i].LastUpdatedDateTime > filteredTrophyTitles[j].LastUpdatedDateTime
		})
	case "progress":
		sort.Slice(filteredTrophyTitles, func(i, j int) bool {
			return filteredTrophyTitles[i].Progress > filteredTrophyTitles[j].Progress
		})
	}

	return &UserTitlesResponse{
		TrophyTitles:   filteredTrophyTitles,
		TotalItemCount: len(filteredTrophyTitles),
	}, nil
}

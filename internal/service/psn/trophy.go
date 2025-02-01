package psn

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

func (s *PSNService) GetUserTitles(accountId string, options *GetUserTitlesOptions) (*UserTitlesResponse, error) {
	auth, err := s.GetValidAuthorization()
	fmt.Println(auth)
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

package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/camtrik/ebbilogue-backend/internal/config"
)

var (
	AuthBaseURL   = "https://ca.account.sony.com/api/authz/v3/oauth"
	TrophyBaseURL = "https://m.np.playstation.com/api/trophy"
)

type TokenData struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int       `json:"expires_in"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type PSNService struct {
	tokenData TokenData
	config    *config.Config
	client    *http.Client
}

type GetUserTitlesOptions struct {
	Limit  *int
	Offset *int
}

type AuthTokensResponse struct {
	AccessToken           string `json:"access_token"`
	ExpiresIn             int    `json:"expires_in"`
	IdToken               string `json:"id_token"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
	Scope                 string `json:"scope"`
	TokenType             string `json:"token_type"`
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

func NewPSNService(config *config.Config) *PSNService {
	return &PSNService{
		config: config,
		client: &http.Client{},
	}
}

func (s *PSNService) GetValidToken() (string, error) {
	auth, err := s.GetValidAuthorization()
	if err != nil {
		return "", err
	}
	return auth.AccessToken, nil
}

func (s *PSNService) GetValidAuthorization() (*TokenData, error) {
	now := time.Now()
	if s.tokenData.AccessToken == "" || s.tokenData.ExpiresAt.Before(now) {
		newAuth, err := s.exchangeRefreshToken()
		if err != nil {
			return nil, fmt.Errorf("failed to refresh token: %v", err)
		}

		s.tokenData = TokenData{
			AccessToken:  newAuth.AccessToken,
			RefreshToken: newAuth.RefreshToken,
			ExpiresIn:    newAuth.ExpiresIn,
			ExpiresAt:    now.Add(time.Duration(newAuth.ExpiresIn) * time.Second),
		}
	}
	return &s.tokenData, nil
}

func (s *PSNService) exchangeRefreshToken() (*AuthTokensResponse, error) {
	formData := url.Values{}
	formData.Set("refresh_token", s.config.PSNRefreshToken)
	formData.Set("grant_type", "refresh_token")
	formData.Set("token_format", "jwt")
	formData.Set("scope", "psn:mobile.v2.core psn:clientapp")

	payload := strings.NewReader(formData.Encode())
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/token", AuthBaseURL), payload)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic MDk1MTUxNTktNzIzNy00MzcwLTliNDAtMzgwNmU2N2MwODkxOnVjUGprYTV0bnRCMktxc1A=")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var authResp AuthTokensResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, err
	}

	return &authResp, nil
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

package psn

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type TokenData struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int       `json:"expires_in"`
	ExpiresAt    time.Time `json:"expires_at"`
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

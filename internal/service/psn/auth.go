// internal/service/psn/auth.go
package psn

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/camtrik/psn-steam-api/internal/global"
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

type NPSSOAuthRequest struct {
	NPSSO string `json:"npsso" binding:"required"`
}

type NPSSOAuthResponse struct {
	AccessToken           string `json:"access_token"`
	RefreshToken          string `json:"refresh_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
	TokenType             string `json:"token_type"`
	Scope                 string `json:"scope"`
}

type AuthorizeResponse struct {
	Location string `json:"location"`
	Code     string `json:"code"`
}

func (s *PSNService) GetValidToken() (string, error) {
	auth, err := s.GetValidAuthorization()
	if err != nil {
		return "", err
	}
	return auth.AccessToken, nil
}

func (s *PSNService) GetTokensFromNPSSO(npsso string) (*NPSSOAuthResponse, error) {
	// 步骤1: 获取access code
	accessCode, err := s.getAccessCodeFromNPSSO(npsso)
	if err != nil {
		s.logger.Errorf("Failed to get access code from NPSSO: %v", err)
		return nil, fmt.Errorf("failed to get access code: %v", err)
	}

	// 步骤2: 使用access code获取tokens
	tokens, err := s.getTokensFromAccessCode(accessCode)
	if err != nil {
		s.logger.Errorf("Failed to get tokens from access code: %v", err)
		return nil, fmt.Errorf("failed to get tokens: %v", err)
	}

	return tokens, nil
}

// Get Refresh Token from NPSSO
func (s *PSNService) getAccessCodeFromNPSSO(npsso string) (string, error) {
	authorizeURL := global.PSN_AUTH_BASE_URL + "/authorize"
	params := url.Values{}
	params.Set("access_type", "offline")
	params.Set("client_id", "09515159-7237-4370-9b40-3806e67c0891")
	params.Set("redirect_uri", "com.scee.psxandroid.scecompcall://redirect")
	params.Set("response_type", "code")
	params.Set("scope", "psn:mobile.v2.core psn:clientapp")

	fullURL := fmt.Sprintf("%s?%s", authorizeURL, params.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return "", err
	}

	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Cookie", fmt.Sprintf("npsso=%s", npsso))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound && resp.StatusCode != http.StatusSeeOther {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	location := resp.Header.Get("Location")
	if location == "" {
		return "", fmt.Errorf("no location header found")
	}

	re := regexp.MustCompile(`code=([^&]+)`)
	matches := re.FindStringSubmatch(location)
	if len(matches) < 2 {
		return "", fmt.Errorf("no access code found in location: %s", location)
	}

	return matches[1], nil
}

func (s *PSNService) getTokensFromAccessCode(accessCode string) (*NPSSOAuthResponse, error) {
	tokenURL := global.PSN_AUTH_BASE_URL + "/token"

	// 准备POST数据
	formData := url.Values{}
	formData.Set("code", accessCode)
	formData.Set("redirect_uri", "com.scee.psxandroid.scecompcall://redirect")
	formData.Set("grant_type", "authorization_code")
	formData.Set("token_format", "jwt")

	// 创建请求
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// 设置headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic MDk1MTUxNTktNzIzNy00MzcwLTliNDAtMzgwNmU2N2MwODkxOnVjUGprYTV0bnRCMktxc1A=")
	req.Header.Set("Accept", "application/json")

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get tokens, status: %d", resp.StatusCode)
	}

	// 解析响应
	var tokenResp NPSSOAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %v", err)
	}

	return &tokenResp, nil
}

func (s *PSNService) UpdateRefreshToken(refreshToken string) {
	s.tokenData.RefreshToken = refreshToken
	s.tokenData.AccessToken = ""
	s.logger.Info("Refresh token updated successfully")
}

// Get Acess Token from Refresh Token
func (s *PSNService) GetValidAuthorization() (*TokenData, error) {
	now := time.Now()
	if s.tokenData.AccessToken == "" || s.tokenData.ExpiresAt.Before(now) {
		newAuth, err := s.exchangeRefreshToken()
		if err != nil {
			s.logger.Errorf("failed to refresh token: %v", err)
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
	formData.Set("refresh_token", s.tokenData.RefreshToken)
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

	if resp.StatusCode != http.StatusOK {
		s.logger.Errorf("failed to get access token: %v", resp.Status)
		return nil, fmt.Errorf("failed to refresh token: %v", resp.Status)
	}
	defer resp.Body.Close()

	var authResp AuthTokensResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, err
	}

	return &authResp, nil
}

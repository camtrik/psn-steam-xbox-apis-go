package psn

import (
	"net/http"
)

var (
	AuthBaseURL   = "https://ca.account.sony.com/api/authz/v3/oauth"
	TrophyBaseURL = "https://m.np.playstation.com/api/trophy"
)

type PSNService struct {
	accountID string
	tokenData TokenData
	client    *http.Client
}

func NewPSNService(accountID, refreshToken string) *PSNService {
	return &PSNService{
		accountID: accountID,
		tokenData: TokenData{
			RefreshToken: refreshToken,
		},
		client: &http.Client{},
	}
}

package psn

import (
	"net/http"

	"github.com/camtrik/ebbilogue-backend/internal/config"
)

var (
	AuthBaseURL   = "https://ca.account.sony.com/api/authz/v3/oauth"
	TrophyBaseURL = "https://m.np.playstation.com/api/trophy"
)

type PSNService struct {
	tokenData TokenData
	config    *config.Config
	client    *http.Client
}

func NewPSNService(config *config.Config) *PSNService {
	return &PSNService{
		config: config,
		client: &http.Client{},
	}
}

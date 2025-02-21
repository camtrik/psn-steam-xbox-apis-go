package psn

import (
	"net/http"

	"github.com/camtrik/psn-steam-api/internal/cache"
	"github.com/camtrik/psn-steam-api/internal/global"
	"github.com/camtrik/psn-steam-api/internal/pkg/logger"
)

var (
	AuthBaseURL   = global.PSN_AUTH_BASE_URL
	TrophyBaseURL = global.PSN_TROPHY_BASE_URL
)

type PSNService struct {
	client    *http.Client
	cache     cache.PSNCache
	tokenData TokenData
	logger    logger.Logger
}

func NewPSNService(client *http.Client, cache cache.PSNCache, logger logger.Logger, refreshToken string) *PSNService {
	return &PSNService{
		client: client,
		cache:  cache,
		logger: logger,
		tokenData: TokenData{
			RefreshToken: refreshToken,
		},
	}
}

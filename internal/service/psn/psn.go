package psn

import (
	"net/http"

	"github.com/camtrik/ebbilogue-backend/internal/cache"
	"github.com/camtrik/ebbilogue-backend/internal/global"
	"github.com/camtrik/ebbilogue-backend/internal/pkg/logger"
)

var (
	AuthBaseURL   = global.AUTH_BASE_URL
	TrophyBaseURL = global.TROPHY_BASE_URL
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

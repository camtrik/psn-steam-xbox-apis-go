package steam

import (
	"net/http"

	"github.com/camtrik/ebbilogue-backend/internal/cache"
	"github.com/camtrik/ebbilogue-backend/internal/global"
	"github.com/camtrik/ebbilogue-backend/internal/pkg/logger"
)

var (
	ApiBaseURL = global.STEAM_API_BASE_URL
)

type SteamService struct {
	client *http.Client
	cache  cache.RedisSteamCache
	logger logger.Logger
	apiKey string
}

func NewSteamService(client *http.Client, cache cache.RedisSteamCache, logger logger.Logger, apiKey string) *SteamService {
	return &SteamService{
		client: client,
		cache:  cache,
		logger: logger,
		apiKey: apiKey,
	}
}

package steam

import (
	"net/http"

	"github.com/camtrik/psn-steam-api/internal/cache"
	"github.com/camtrik/psn-steam-api/internal/global"
	"github.com/camtrik/psn-steam-api/internal/pkg/logger"
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

type GameDetails struct {
	AppId              int
	Name               string
	PlayTime           int
	PlayTime2weeks     int
	Achieved           int
	TotalAchievenments int
	IconUrl            string
	ArtUrl             string
	StoreUrl           string
}

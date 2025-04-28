package xbox

import (
	"net/http"

	"github.com/camtrik/psn-steam-api/internal/cache"
	"github.com/camtrik/psn-steam-api/internal/global"
	"github.com/camtrik/psn-steam-api/internal/pkg/logger"
)

var (
	ApiBaseURL = global.XBOX_API_BASE_URL
)

type XboxService struct {
	client *http.Client
	logger logger.Logger
	cache  cache.XboxCache
	apiKey string
}

func NewXboxService(client *http.Client, logger logger.Logger, cache cache.XboxCache, apiKey string) *XboxService {
	return &XboxService{
		client: client,
		logger: logger,
		cache:  cache,
		apiKey: apiKey,
	}
}

package xbox

import (
	"net/http"

	"github.com/camtrik/psn-steam-api/internal/global"
	"github.com/camtrik/psn-steam-api/internal/pkg/logger"
)

var (
	ApiBaseURL = global.XBOX_API_BASE_URL
)

type XboxService struct {
	client *http.Client
	logger logger.Logger
	apiKey string
}

func NewXboxService(client *http.Client, logger logger.Logger, apiKey string) *XboxService {
	return &XboxService{
		client: client,
		logger: logger,
		apiKey: apiKey,
	}
}


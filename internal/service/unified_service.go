package unified

import (
	"github.com/camtrik/psn-steam-api/internal/pkg/logger"
	"github.com/camtrik/psn-steam-api/internal/service/psn"
	"github.com/camtrik/psn-steam-api/internal/service/steam"
	"github.com/camtrik/psn-steam-api/internal/service/xbox"
)

type UnifiedGameService struct {
	SteamService *steam.SteamService
	PSNService   *psn.PSNService
	XboxService  *xbox.XboxService
	logger       logger.Logger
}

func NewUnifiedGameService(
	steamService *steam.SteamService,
	psnService *psn.PSNService,
	xboxService *xbox.XboxService,
	logger logger.Logger) *UnifiedGameService {
	return &UnifiedGameService{
		SteamService: steamService,
		PSNService:   psnService,
		XboxService:  xboxService,
		logger:       logger,
	}
}

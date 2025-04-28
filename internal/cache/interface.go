package cache

import (
	"context"

	"github.com/camtrik/psn-steam-api/internal/models"
)

type PSNCache interface {
	GetUserTitles(ctx context.Context, accountId string) (*models.UserTitlesResponse, error)
	SetUserTitles(ctx context.Context, accountId string, userTitles *models.UserTitlesResponse) error
}

type SteamCache interface {
	GetOwnedGames(ctx context.Context, steamId string) (*models.OwnedGamesResponse, error)
	SetOwnedGames(ctx context.Context, steamId string, ownedGames *models.OwnedGamesResponse) error
	GetPlayerAchievements(ctx context.Context, steamId string, appId int) (*models.PlayerAchievementsResponse, error)
	SetPlayerAchievements(ctx context.Context, steamId string, appId int, playerAchievements *models.PlayerAchievementsResponse) error
}

type XboxCache interface {
	GetPlayerAchievements(ctx context.Context) (*models.XboxGamaAchievements, error)
	SetPlayerAchievements(ctx context.Context, playerAchievements *models.XboxGamaAchievements) error
	GetGameStats(ctx context.Context, titleId string) (*models.XboxGameStats, error)
	SetGameStats(ctx context.Context, titleId string, gameStats *models.XboxGameStats) error
}

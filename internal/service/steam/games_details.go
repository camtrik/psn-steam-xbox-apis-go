package steam

import (
	"context"
	"fmt"
	"sync"

	"github.com/camtrik/ebbilogue-backend/internal/global"
	"github.com/camtrik/ebbilogue-backend/internal/models"
)

type PlayerGameDetails struct {
	GameCount int
	Games     []GameDetails
}

type GameDetails struct {
	AppId              int
	Name               string
	PlayTime           int
	Achieved           int
	TotalAchievenments int
	IconUrl            string
	ArtUrl             string
}

func (s *SteamService) GetPlayerGameDetails(ctx context.Context, steamId string) (PlayerGameDetails, error) {
	ownedGames, err := s.GetOwnedGames(ctx, steamId)
	if err != nil {
		return PlayerGameDetails{}, err
	}

	// concurrency
	var (
		mu    sync.Mutex
		wg    sync.WaitGroup
		games []GameDetails
	)

	for _, game := range ownedGames.Response.Games {
		wg.Add(1)
		go func(game models.SteamGame) {
			defer wg.Done()
			achieved, totalAchievements, err := s.GetGameAchievements(ctx, steamId, game.AppId)
			if err != nil {
				s.logger.Error("failed to get game achievements for appid %d: %v", game.AppId, err)
				return
			}

			imgIconUrl := fmt.Sprintf("%s/%d/%s.jpg", global.STEAM_ICON_BASE_URL, game.AppId, game.ImgIconUrl)
			imgArtUrl := fmt.Sprintf("%s/%d/%s", global.STEAM_CAPSULE_BASE_URL, game.AppId, global.STEAM_CAPSULE_ART_LARGE)

			mu.Lock()
			games = append(games, GameDetails{
				AppId:              game.AppId,
				Name:               game.Name,
				PlayTime:           game.PlayTimeForever,
				Achieved:           achieved,
				TotalAchievenments: totalAchievements,
				IconUrl:            imgIconUrl,
				ArtUrl:             imgArtUrl,
			})
			mu.Unlock()
		}(game)
	}
	wg.Wait()

	return PlayerGameDetails{
		GameCount: ownedGames.Response.GameCount,
		Games:     games,
	}, nil
}

// return achieved, total achievements
func (s *SteamService) GetGameAchievements(ctx context.Context, steamId string, appId int) (achieved int, totalAchievements int, err error) {
	playerAchievements, err := s.GetPlayerAchievements(ctx, steamId, appId)
	if err != nil {
		return 0, 0, err
	}
	totalAchievements = len(playerAchievements.Playerstats.Achievements)
	for _, achievement := range playerAchievements.Playerstats.Achievements {
		if achievement.Achieved == 1 {
			achieved++
		}
	}
	return achieved, totalAchievements, nil
}

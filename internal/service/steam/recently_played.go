package steam

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/camtrik/psn-steam-api/internal/global"
	"github.com/camtrik/psn-steam-api/internal/models"
)

func (s *SteamService) GetRecentlyPlayedGames(ctx context.Context, steamId string, timeRange int64) ([]models.RecentlyPlayedGame, error) {
	ownedGames, err := s.GetOwnedGames(ctx, steamId)
	if err != nil {
		return nil, err
	}

	timeStamp := time.Now().Unix() - timeRange

	var (
		mu    sync.Mutex
		wg    sync.WaitGroup
		games []models.RecentlyPlayedGame
	)

	for _, game := range ownedGames.Response.Games {
		if game.RTimeLastPlayed < timeStamp {
			continue
		}

		wg.Add(1)
		go func(game models.SteamGame) {
			defer wg.Done()
			achieved, _, err := s.GetGameAchievements(ctx, steamId, game.AppId)
			if err != nil {
				s.logger.Error("Failed to get game achievements %v", err)
				return
			}

			imgArtUrl := fmt.Sprintf("%s/%d/%s", global.STEAM_CAPSULE_BASE_URL, game.AppId, global.STEAM_CAPSULE_ART_LARGE)
			storeUrl := fmt.Sprintf("%s/%d", global.STEAM_STORE_BASE_URL, game.AppId)

			mu.Lock()
			games = append(games, models.RecentlyPlayedGame{
				Name:               game.Name,
				PlayTime:           game.PlayTimeForever,
				LastPlayedTime:     game.RTimeLastPlayed,
				EarnedAchievements: achieved,
				ArtUrl:             imgArtUrl,
				StoreUrl:           storeUrl,
				Platform:           "steam",
			})
			mu.Unlock()
		}(game)
	}
	wg.Wait()

	return games, nil
}

package unified

import (
	"context"
	"sort"
	"sync"

	"github.com/camtrik/psn-steam-api/internal/models"
)

func (s *UnifiedGameService) GetRecentlyPlayedGames(
	ctx context.Context,
	psnAccountId string,
	steamId string,
	timeRange int64,
) ([]models.RecentlyPlayedGame, error) {
	var (
		wg    sync.WaitGroup
		mu    sync.Mutex
		games []models.RecentlyPlayedGame
	)

	wg.Add(3)

	go func() {
		defer wg.Done()
		psnGames, err := s.PSNService.GetRecentlyPlayedGames(ctx, psnAccountId, timeRange)
		if err != nil {
			s.logger.Error("failed to get psn recently played games", "error", err)
			return
		}

		mu.Lock()
		games = append(games, psnGames...)
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		steamGames, err := s.SteamService.GetRecentlyPlayedGames(ctx, steamId, timeRange)
		if err != nil {
			s.logger.Error("failed to get steam recently played games", "error", err)
			return
		}

		mu.Lock()
		games = append(games, steamGames...)
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		xboxGames, err := s.XboxService.GetRecentlyPlayedGames(ctx, timeRange)
		if err != nil {
			s.logger.Error("failed to get xbox recently played games", "error", err)
			return
		}

		mu.Lock()
		games = append(games, xboxGames...)
		mu.Unlock()
	}()

	wg.Wait()

	sort.Slice(games, func(i, j int) bool {
		return games[i].LastPlayedTime > games[j].LastPlayedTime
	})

	return games, nil
}

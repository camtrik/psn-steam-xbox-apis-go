package steam

import (
	"context"
	"fmt"
	"sort"
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
	games := []models.RecentlyPlayedGame{}

	for _, game := range ownedGames.Response.Games {
		if game.RTimeLastPlayed > timeStamp {
			imgIconUrl := fmt.Sprintf("%s/%d/%s.jpg", global.STEAM_ICON_BASE_URL, game.AppId, game.ImgIconUrl)
			imgArtUrl := fmt.Sprintf("%s/%d/%s", global.STEAM_CAPSULE_BASE_URL, game.AppId, global.STEAM_CAPSULE_ART_LARGE)
			storeUrl := fmt.Sprintf("%s/%d", global.STEAM_STORE_BASE_URL, game.AppId)

			games = append(games, models.RecentlyPlayedGame{
				AppId:          game.AppId,
				Name:           game.Name,
				PlayTime:       game.PlayTimeForever,
				LastPlayedTime: game.RTimeLastPlayed,
				IconUrl:        imgIconUrl,
				ArtUrl:         imgArtUrl,
				StoreUrl:       storeUrl,
			})
		}
	}

	sort.Slice(games, func(i, j int) bool {
		return games[i].LastPlayedTime > games[j].LastPlayedTime
	})

	return games, nil
}

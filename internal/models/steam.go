package models

type SteamGame struct {
	AppId           int    `json:"appid"`
	Name            string `json:"name"`
	ImgIconUrl      string `json:"img_icon_url"`
	PlayTimeForever int    `json:"playtime_forever"`
	RTimeLastPlayed int    `json:"rtime_last_played"`
}

type OwnedGamesResponse struct {
	Response struct {
		GameCount int         `json:"game_count"`
		Games     []SteamGame `json:"games"`
	} `json:"response"`
}

type PlayerAchievementsResponse struct {
	Playerstats struct {
		SteamID      string `json:"steamID"`
		GameName     string `json:"gameName"`
		Achievements []struct {
			APIName    string `json:"apiname"`
			Achieved   int    `json:"achieved"`
			UnlockTime int    `json:"unlocktime"`
		} `json:"achievements"`
		Success bool `json:"success"`
	} `json:"playerstats"`
}

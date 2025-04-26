package models

import "time"

type XboxGameTitles struct {
	TitleId      string `json:"titleId"`
	Name         string `json:"name"`
	DisplayImage string `json:"displayImage"`
	Achievement  struct {
		CurrentAchievements int `json:"currentAchievements"`
		TotalAchievements   int `json:"totalAchievements"`
		CurrentGamerscore   int `json:"currentGamerscore"`
		TotalGamerscore     int `json:"totalGamerscore"`
	} `json:"achievement"`
	TitleHistory struct {
		LastTimePlayed time.Time `json:"lastTimePlayed"`
	}
}

type XboxGamaAchievements struct {
	Xuid   string
	Titles []XboxGameTitles `json:"titles"`
}

type XboxGameStats struct {
	StatListsCollection []struct {
		Stats []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"stats"`
	} `json:"statListsCollection"`
}

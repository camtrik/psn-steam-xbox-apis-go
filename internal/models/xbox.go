package models

import "time"

type XboxGamaAchievements struct {
	Xuid   string
	Titles []struct {
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
	} `json:"titles"`
}

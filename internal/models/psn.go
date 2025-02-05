package models

type GetUserTitlesOptions struct {
	Limit  *int
	Offset *int
}

type UserTitlesResponse struct {
	TrophyTitles   []TrophyTitle `json:"trophyTitles"`
	TotalItemCount int           `json:"totalItemCount"`
}

type TrophyTitle struct {
	NpCommunicationId   string          `json:"npCommunicationId"`
	TrophySetVersion    string          `json:"trophySetVersion"`
	TrophyTitleName     string          `json:"trophyTitleName"`
	TrophyTitlePlatform string          `json:"trophyTitlePlatform"`
	TrophyTitleIconUrl  string          `json:"trophyTitleIconUrl"`
	HasEarnedTrophies   bool            `json:"hasEarnedTrophies"`
	DefinedTrophies     DefinedTrophies `json:"definedTrophies"`
	Progress            int             `json:"progress"`
	EarnedTrophies      EarnedTrophies  `json:"earnedTrophies"`
	LastUpdatedDateTime string          `json:"lastUpdatedDateTime"`
}

type DefinedTrophies struct {
	Bronze   int `json:"bronze"`
	Silver   int `json:"silver"`
	Gold     int `json:"gold"`
	Platinum int `json:"platinum"`
}

type EarnedTrophies struct {
	Bronze   int `json:"bronze"`
	Silver   int `json:"silver"`
	Gold     int `json:"gold"`
	Platinum int `json:"platinum"`
}

type TrophyFilter struct {
	MinProgress int
	Platform    string
	SortBy      string
}

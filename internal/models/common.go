package models

// common models for all services

type RecentlyPlayedGame struct {
	Name           string
	PlayTime       int
	LastPlayedTime int64
	ArtUrl         string
	StoreUrl       string
	Platform       string
}

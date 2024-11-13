package ghclient

import "time"

type UserData struct {
	Username    string `json:"login"`
	ReposApiURL string `json:"repos_url"`
	Followers   int    `json:"followers"`
}

type RepoData struct {
	Name            string    `json:"name"`
	LanguagesApiURL string    `json:"languages_url"`
	ForksCount      int       `json:"forks_count"`
	CreatedAt       time.Time `json:"created_at"`
	PushedAt        time.Time `json:"pushed_at"`
}

type LanguageDistribution map[string]float64
type UserActivity map[int]int

type UserFullData struct {
	UserData             UserData
	Repos                []RepoData
	LanguageDistribution LanguageDistribution
}

type UserFormattedData struct {
	Username             string
	Followers            int
	ForksCount           int
	RepoCount            int
	LanguageDistribution LanguageDistribution
	UserActivity         UserActivity
}

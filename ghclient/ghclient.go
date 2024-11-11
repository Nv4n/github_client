package ghclient

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

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

func fetchGithubData[ReturnType UserData | []RepoData | map[string]interface{}](client *http.Client, request *http.Request) ReturnType {
	res, err := client.Do(request)
	if err != nil {
		log.Fatalf("error on doing request: %v\n", err)
	}
	if res == nil {
		log.Fatalf("error not getting any response or hitting rate limit")
	}

	if res.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Fatalf("Error accessing body: %v", err)
			}
		}(res.Body)
	}

	var data ReturnType
	body, _ := io.ReadAll(res.Body)
	jsonErr := json.Unmarshal(body, &data)
	if jsonErr != nil {
		log.Fatalf("error json parsing: %v\n", jsonErr)
	}
	return data
}

func getLanguageApiURLs(repos []RepoData, repoLimit int) []string {
	var langApiList []string
	for i, repoData := range repos {
		if i >= repoLimit {
			break
		}
		langApiList = append(langApiList, repoData.LanguagesApiURL)
	}
	return langApiList
}

func calcLangDistribution(distribution LanguageDistribution, percentageThreshold float64) LanguageDistribution {
	totalLines := float64(0)
	langDistribution := make(map[string]float64)
	for _, lines := range distribution {
		totalLines += lines
	}
	for lang, val := range distribution {
		percentage := val * 100.0 / totalLines
		if percentage > percentageThreshold {
			langDistribution[lang] = percentage
		}
	}
	return langDistribution
}

func calcTotalForksCount(repos []RepoData) int {
	count := 0
	for _, r := range repos {
		count = count + r.ForksCount
	}
	return count
}

func calcUserActivity(repos []RepoData) UserActivity {
	userActivity := make(map[int]int)
	for _, r := range repos {
		pushedAt := r.PushedAt.Year()
		createdAt := r.CreatedAt.Year()
		userActivity[pushedAt] += 1
		if createdAt < pushedAt {
			for y := createdAt; y < pushedAt; y++ {
				userActivity[y] += 1
			}
		}
	}
	return userActivity
}

func GetUserData(username string, repoLimit int, langThreshold float64) UserFormattedData {
	client := http.Client{
		Timeout: time.Second * 10,
	}

	user := UserFullData{}
	userRequest, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.github.com/users/%s", username), nil)
	user.UserData = fetchGithubData[UserData](&client, userRequest)
	reposRequest, _ := http.NewRequest(http.MethodGet, user.UserData.ReposApiURL, nil)
	user.Repos = fetchGithubData[[]RepoData](&client, reposRequest)

	if repoLimit == -1 {
		repoLimit = len(user.Repos)
	}

	langApiList := getLanguageApiURLs(user.Repos, repoLimit)

	languageKBList := make(map[string]float64)
	for _, url := range langApiList {
		languageRequest, _ := http.NewRequest(http.MethodGet, url, nil)

		repoLangUsage := fetchGithubData[map[string]interface{}](&client, languageRequest)
		for lang, val := range repoLangUsage {
			if v, ok := val.(float64); ok {
				languageKBList[lang] = languageKBList[lang] + v/1024.0
			}
		}
	}

	user.LanguageDistribution = calcLangDistribution(languageKBList, langThreshold)
	totalForkCount := calcTotalForksCount(user.Repos)
	userActivity := calcUserActivity(user.Repos)

	return UserFormattedData{
		user.UserData.Username,
		user.UserData.Followers,
		totalForkCount,
		len(user.Repos),
		user.LanguageDistribution,
		userActivity,
	}
}

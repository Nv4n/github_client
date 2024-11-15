package ghclient

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup

func fetchGithubData[ReturnType UserData | []RepoData | map[string]interface{}](client *http.Client, request *http.Request) ReturnType {
	res, err := client.Do(request)
	if err != nil {
		log.Fatalf("error on doing request: %v\nREQUEST: %v\n", err, request)
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

func GetUserData(username string, repoLimit int, langThreshold float64) UserFormattedData {
	client := http.Client{
		Timeout: time.Second * 10,
	}

	token := os.Getenv("GH_TOKEN")

	if token == "" {
		log.Fatal(fmt.Errorf("empty gh_token"))
	}

	user := UserFullData{}
	userRequest, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.github.com/users/%s", username), nil)
	userRequest.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	user.UserData = fetchGithubData[UserData](&client, userRequest)

	reposRequest, _ := http.NewRequest(http.MethodGet, user.UserData.ReposApiURL, nil)
	reposRequest.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	user.Repos = fetchGithubData[[]RepoData](&client, reposRequest)

	if repoLimit == -1 {
		repoLimit = len(user.Repos)
	}

	langApiList := getLanguageApiURLs(user.Repos, repoLimit)

	languageKBList := make(map[string]float64)
	semaphore := make(chan struct{}, 10)
	for _, url := range langApiList {
		wg.Add(1)
		semaphore <- struct{}{}
		languageRequest, _ := http.NewRequest(http.MethodGet, url, nil)
		languageRequest.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		go func() {
			defer wg.Done()
			repoLangUsage := fetchGithubData[map[string]interface{}](&client, languageRequest)
			for lang, val := range repoLangUsage {
				if v, ok := val.(float64); ok {
					languageKBList[lang] = languageKBList[lang] + v/1024.0
				}
			}
			<-semaphore
		}()
	}
	wg.Wait()
	close(semaphore)

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

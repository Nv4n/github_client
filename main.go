package main

import (
	"bufio"
	"fmt"
	gh "ghclient/ghclient"
	pres "ghclient/presenter"
	"github.com/joho/godotenv"
	"log"
	"os"
	"sync"
)

var wg sync.WaitGroup

func readUsernames(file *os.File) []string {
	var usernames []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		username := scanner.Text()
		usernames = append(usernames, username)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error during reading username: %v\n", err)
	}
	return usernames
}
func fetchUsers(usernames []string, repoLimit int, langThreshold float64) []gh.UserFormattedData {
	var users []gh.UserFormattedData
	for _, u := range usernames {
		wg.Add(1)
		go func() {
			defer wg.Done()
			data := gh.GetUserData(u, repoLimit, langThreshold)
			users = append(users, data)
		}()
	}
	wg.Wait()
	return users
}

func main() {
	//TODO
	//	Add flags for filename and repoLimit
	//	Fix language % formula
	//	Add goroutines
	//	Add WaitGroups
	//	Add web representation with e-charts

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		log.Fatal("No filenames in the cli arguments")
	}
	filename := os.Args[1]
	pwd, _ := os.Getwd()
	open, err := os.Open(fmt.Sprintf("%s\\public\\%s", pwd, filename))
	if err != nil {
		log.Fatalf("error opening filename %s: %v", filename, err)
		return
	}

	fmt.Println("Reading usernames...")
	usernames := readUsernames(open)
	fmt.Println("Fetching data...")
	//RepoLimit: -1 FOR NO LIMIT
	//Language Threshold: min percentage to be included
	users := fetchUsers(usernames, -1, 1)
	pres.PresentGhData(users)

	defer func(open *os.File) {
		err := open.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(open)

}

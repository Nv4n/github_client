package main

import (
	"bufio"
	"flag"
	"fmt"
	gh "ghclient/ghclient"
	pres "ghclient/presenter"
	"github.com/joho/godotenv"
	"log"
	"os"
	"sync"
)

var fileDir string
var isCliPres bool
var isWebPres bool
var repoLimit int
var langThreshold float64
var isHelpPres bool

var wg sync.WaitGroup

func init() {
	flag.StringVar(&fileDir, "fileDir", "public\\usernames.txt", ".txt directory of all usernames")
	flag.BoolVar(&isCliPres, "cli", false, "output in isCliPres")
	flag.BoolVar(&isWebPres, "web", false, "output in isWebPres server")
	flag.IntVar(&repoLimit, "repoLimit", -1, "-1 FOR NO LIMIT")
	flag.Float64Var(&langThreshold, "langThreshold", 1, "min percentage to be included in output data")
	flag.BoolVar(&isHelpPres, "help", false, "prints defaults")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
}

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

func fetchUsers(usernames []string) []gh.UserFormattedData {
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
	//	Add isWebPres representation with e-charts

	flag.Parse()

	if isHelpPres {
		flag.PrintDefaults()
		return
	}

	pwd, _ := os.Getwd()
	open, err := os.Open(fmt.Sprintf("%s\\%s", pwd, fileDir))

	if err != nil {
		log.Fatalf("error opening fileDir %s: %v", fileDir, err)
	}

	defer func(open *os.File) {
		err := open.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(open)

	fmt.Println("Reading usernames...")
	usernames := readUsernames(open)
	fmt.Println("Fetching data...")
	users := fetchUsers(usernames)
	pres.PresentGhData(users)
}

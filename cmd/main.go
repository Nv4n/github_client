package main

import (
	"bufio"
	"fmt"
	gh "ghclient_homework/ghclient"
	"github.com/olekukonko/tablewriter"
	"log"
	"os"
	"sort"
	"strings"
)

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
		data := gh.GetUserData(u, repoLimit, langThreshold)
		users = append(users, data)
	}
	return users
}

func formatLangDist(langDist gh.LanguageDistribution) string {
	var builder strings.Builder
	for l, dist := range langDist {
		builder.WriteString(fmt.Sprintf("%s: %.2f%% / ", l, dist))
	}
	return builder.String()
}

func formatUserActivity(userActivity gh.UserActivity) string {
	var builder strings.Builder
	var years []int
	for year := range userActivity {
		years = append(years, year)
	}
	sort.Slice(years, func(l, r int) bool {
		return years[l] > years[r]
	})
	for _, y := range years {
		builder.WriteString(fmt.Sprintf("Y(%v): %v / ", y, userActivity[y]))
	}
	return builder.String()
}

func presentGhData(users []gh.UserFormattedData) {
	columns := []string{"Username", "Followers", "Forks", "Repo Count", "Language usage", "User Activity"}
	tbl := tablewriter.NewWriter(os.Stdout)
	tbl.SetHeader(columns)
	tbl.SetAutoFormatHeaders(true)
	tbl.SetBorder(true)
	tbl.SetRowSeparator("=")
	tbl.SetAutoWrapText(true)

	for _, u := range users {
		langDist := formatLangDist(u.LanguageDistribution)
		userActivity := formatUserActivity(u.UserActivity)
		data := []string{u.Username,
			fmt.Sprintf("%d", u.Followers),
			fmt.Sprintf("%d", u.ForksCount),
			fmt.Sprintf("%d", u.RepoCount),
			langDist,
			userActivity}
		tbl.Append(data)
	}

	tbl.Render()
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("No filenames in the cli arguments")
	}
	filename := os.Args[1]
	pwd, _ := os.Getwd()
	open, err := os.Open(fmt.Sprintf("%s\\cmd\\%s", pwd, filename))
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
	presentGhData(users)

	defer func(open *os.File) {
		err := open.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(open)

}

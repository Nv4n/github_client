package presenter

import (
	"fmt"
	gh "ghclient/ghclient"
	"github.com/olekukonko/tablewriter"
	"os"
	"sort"
	"strings"
)

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

func PresentCli(users []gh.UserFormattedData) {
	columns := []string{"Username", "Followers", "Forks", "Repo Count", "Language usage", "User Activity"}
	tbl := tablewriter.NewWriter(os.Stdout)
	tbl.SetHeader(columns)
	tbl.SetAutoFormatHeaders(true)
	tbl.SetBorder(true)
	tbl.SetRowSeparator("=")
	tbl.SetRowLine(true)
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

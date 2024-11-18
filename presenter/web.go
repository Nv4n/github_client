package presenter

import (
	"fmt"
	"ghclient/ghclient"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"html/template"
	"log"
	"net/http"
)

type HtmlUserData struct {
	Username   string
	Followers  int
	ForksCount int
	RepoCount  int
	Pie        []byte
	Line       []byte
}

func PresentWeb(data []ghclient.UserFormattedData) {
	htmlData := make([]HtmlUserData, len(data))

	for i, user := range data {
		htmlData[i].Username = user.Username
		htmlData[i].Followers = user.Followers
		htmlData[i].ForksCount = user.ForksCount
		htmlData[i].RepoCount = user.RepoCount
		htmlData[i].Pie = getPie(user.LanguageDistribution)
		htmlData[i].Line = getLine(user.UserActivity)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("template.html"))

		_ = tmpl.Execute(w, htmlData)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getLine(data ghclient.UserActivity) []byte {
	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Activity Distribution",
		}),
		charts.WithLegendOpts(opts.Legend{
			Icon: "circle",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Type: "category",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Min: 0,
		}),
	)

	lineData := make([]opts.LineData, 0)
	for k, v := range data {
		lineData = append(lineData, opts.LineData{Value: float64(v), Name: fmt.Sprintf("Year (%d)", k)})
	}

	line.AddSeries("Data", lineData)
	return line.RenderContent()
}

func getPie(data ghclient.LanguageDistribution) []byte {
	pie := charts.NewPie()

	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Language Distribution",
		}),
		charts.WithLegendOpts(opts.Legend{
			Icon: "circle",
		}),
	)

	pieData := make([]opts.PieData, 0)
	for k, v := range data {
		pieData = append(pieData, opts.PieData{Name: k, Value: v})
	}

	pie.AddSeries("Language Distribution", pieData)
	buffer := pie.RenderContent()
	return buffer
}

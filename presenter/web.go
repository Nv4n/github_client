package presenter

import (
	"fmt"
	"ghclient/ghclient"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/render"
	"html/template"
	"log"
	"math"
	"net/http"
	"sort"
)

type HtmlUserData struct {
	Username   string
	Followers  int
	ForksCount int
	RepoCount  int
	Pie        template.HTML
	Line       template.HTML
}

func PresentWeb(data []ghclient.UserFormattedData) {
	htmlData := make([]HtmlUserData, len(data))

	for i, user := range data {
		htmlData[i].Username = user.Username
		htmlData[i].Followers = user.Followers
		htmlData[i].ForksCount = user.ForksCount
		htmlData[i].RepoCount = user.RepoCount
		htmlData[i].Pie = template.HTML(getPie(user.LanguageDistribution).RenderContent())
		htmlData[i].Line = template.HTML(getLine(user.UserActivity).RenderContent())
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("views/template.go.html"))

		_ = tmpl.ExecuteTemplate(w, "Content", htmlData)
	})

	fmt.Println("Listening to localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getLine(data ghclient.UserActivity) render.Renderer {
	line := charts.NewLine()
	var years []int
	for k, _ := range data {
		years = append(years, k)
	}
	sort.Ints(years)

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:  "Activity Distribution",
			Bottom: "Year",
		}),
		charts.WithLegendOpts(opts.Legend{
			Icon: "circle",
			X:    "Year",
			Y:    "Activity",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Type: "category",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Min: 0,
		}),
	)
	line.SetXAxis(years)

	lineData := make([]opts.LineData, 0)
	ind := 0
	for _, k := range years {
		lineData = append(lineData, opts.LineData{YAxisIndex: ind, Value: float64(data[k]), Name: fmt.Sprintf("Activity Y(%v)", k)})
		ind++
	}

	line.AddSeries("Data", lineData)
	return line.Renderer
}

func getPie(data ghclient.LanguageDistribution) render.Renderer {
	pie := charts.NewPie()

	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Language Distribution",
		}),
		charts.WithLegendOpts(opts.Legend{
			Icon: "circle",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Min:     "5rem",
			NameGap: 2,
		}),
	)

	pieData := make([]opts.PieData, 0)
	for k, v := range data {
		pieData = append(pieData, opts.PieData{Name: k, Value: math.Round(v*100) / 100})
	}

	pie.AddSeries("Language Distribution", pieData)
	return pie.Renderer
}

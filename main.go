package main

import (
	"html/template"
	"math/rand"
	"net/http"
	"io"
	"bytes"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/labstack/echo/v4"
)


var (
	itemCntLine = 6
	fruits      = []string{"Apple", "Banana", "Peach ", "Lemon", "Pear", "Cherry"}
)

// generate random data for bar chart
func generateBarItems() []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < 7; i++ {
		items = append(items, opts.BarData{Value: rand.Intn(300)})
	}
	return items
}


type Template struct {
	templates *template.Template
}
      
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type ServiceInfo struct {
	Title string
      }

var serviceInfo = ServiceInfo {
	"サイトのタイトル",
}

func generateLineItems() []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < itemCntLine; i++ {
		items = append(items, opts.LineData{Value: rand.Intn(300)})
	}
	return items
}

func main() {
	line := charts.NewLine()
	line.Renderer = NewSnippetRenderer(line, line.Validate)

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "title and label options",
			Subtitle: "go-echarts is an awesome chart library written in Golang",
		}),
	)

	line.SetXAxis(fruits).
		AddSeries("Category A", generateLineItems()).
		AddSeries("Category B", generateLineItems()).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show: true,
			}),
		)


	var buf bytes.Buffer
	line.Render(&buf)

	e := echo.New()
	temp := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	e.Renderer = temp
	e.GET("/", func(c echo.Context) error {
		data := struct {
			ServiceInfo
			Body template.HTML
		      } {
			ServiceInfo: serviceInfo,
			Body: template.HTML(buf.String()),
		      }
		// return c.Render(http.StatusOK, "page1", data)
		return c.Render(http.StatusOK, "template1", data)
	})
	e.Logger.Fatal(e.Start(":1323"))
}

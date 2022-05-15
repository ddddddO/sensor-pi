package main

import (
	"database/sql"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type Plot struct {
	title                  string
	imageName              string
	yLabel, xLabel         string
	lineColor, pointsColor color.RGBA
}

func newPressurePlot() *Plot {
	pressurePlot := &Plot{
		title:       "pressure @around tama river",
		imageName:   "pressure.png",
		yLabel:      "pressure",
		xLabel:      "date",
		lineColor:   color.RGBA{G: 255, A: 255},
		pointsColor: color.RGBA{R: 255, A: 255},
	}
	return pressurePlot
}

// NOTE: 気温・湿度のグラフを生成する際に以下を編集する
// func newTemperaturePlot() *Plot { return &Plot{} }
// func newHumidityPlot() *Plot    { return &Plot{} }

func (pt *Plot) Save(points plotter.XYs) error {
	p := plot.New()
	p.Title.Text = pt.title
	p.Y.Label.Text = pt.yLabel
	p.X.Label.Text = pt.xLabel
	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04"}
	p.Add(plotter.NewGrid())

	lpLine, lpPoints, err := plotter.NewLinePoints(points)
	if err != nil {
		return err
	}
	lpLine.Color = pt.lineColor
	lpPoints.Shape = draw.PlusGlyph{}
	lpPoints.Color = pt.pointsColor
	p.Add(lpLine, lpPoints)
	// p.Legend.Add("line points", lpLine, lpPoints)

	if err := p.Save(4*vg.Inch, 4*vg.Inch, filepath.Join(baseDir, plotterDir, pt.imageName)); err != nil {
		return err
	}
	return nil
}

// const baseDir = "/home/pi/github.com/ddddddO/sensor-pi/" // raspberry pi
const baseDir = "/mnt/c/DEV/workspace/GO/src/github.com/ddddddO/sensor-pi/" // wsl
const plotterDir = "plotter"

// ref: https://github.com/gonum/plot/wiki/Example-plots#more-detailed-style-settings
func main() {
	environment, err := fetchData()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	points, err := points(environment.pressure)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	pressurePlot := newPressurePlot()
	if err := pressurePlot.Save(points); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type environment struct {
	pressure    []datum
	temperature []datum
	humidity    []datum
}

type datum struct {
	t     time.Time
	value float64
}

func fetchData() (*environment, error) {
	var dsn = filepath.Join(baseDir, "environment.sqlite3")
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	const query = "select date, temperature, pressure, humidity from environment order by date desc limit 10"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	const layout = "2006-01-02 15:04:05"
	var (
		dataP = []datum{} // for pressure
		dataT = []datum{} // for temperature
		dataH = []datum{} // for humidity
	)
	for rows.Next() {
		var (
			d       string
			t, p, h float64
		)
		if err := rows.Scan(&d, &t, &p, &h); err != nil {
			return nil, err
		}

		tm, err := time.Parse(layout, d)
		if err != nil {
			return nil, err
		}

		dataP = append(dataP, datum{t: tm, value: p})
		dataT = append(dataT, datum{t: tm, value: t})
		dataH = append(dataH, datum{t: tm, value: h})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	e := &environment{
		pressure:    dataP,
		temperature: dataT,
		humidity:    dataH,
	}
	return e, nil
}

func points(data []datum) (plotter.XYs, error) {
	pts := make(plotter.XYs, len(data))
	for i := range data {
		pts[i].X = float64(data[i].t.Unix())
		pts[i].Y = data[i].value
	}
	return pts, nil
}

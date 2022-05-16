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

	sq "github.com/Masterminds/squirrel"
)

const baseDir = "/home/pi/github.com/ddddddO/sensor-pi/env-bot/" // raspberry pi
// const baseDir = "/mnt/c/DEV/workspace/GO/src/github.com/ddddddO/sensor-pi/env-bot/" // wsl
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

	const limit = 10
	query := sq.Select("date", "temperature", "pressure", "humidity").
		From("environment").
		OrderBy("date desc").
		Limit(limit)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(sql, args...)
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

type Plot struct {
	*plot.Plot
	imagePath              string
	lineColor, pointsColor color.Color
}

func newPressurePlot() *Plot {
	p := plot.New()
	p.Title.Text = "pressure @around tama river"
	p.Y.Label.Text = "pressure"
	p.X.Label.Text = "date"
	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04"}
	p.Add(plotter.NewGrid())

	return &Plot{
		Plot:        p,
		imagePath:   filepath.Join(baseDir, plotterDir, "pressure.png"),
		lineColor:   color.RGBA{G: 255, A: 255},
		pointsColor: color.RGBA{R: 255, A: 255},
	}
}

// NOTE: 気温・湿度のグラフを生成する際に以下を編集する
// func newTemperaturePlot() *Plot { return &Plot{} }
// func newHumidityPlot() *Plot    { return &Plot{} }

func (p *Plot) Save(points plotter.XYs) error {
	lpLine, lpPoints, err := plotter.NewLinePoints(points)
	if err != nil {
		return err
	}
	lpLine.Color = p.lineColor
	lpPoints.Shape = draw.PlusGlyph{}
	lpPoints.Color = p.pointsColor
	p.Add(lpLine, lpPoints)
	// p.Legend.Add("line points", lpLine, lpPoints)

	return p.Plot.Save(4*vg.Inch, 4*vg.Inch, p.imagePath)
}

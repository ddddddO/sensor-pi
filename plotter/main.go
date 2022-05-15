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

// const baseDir = "/home/pi/github.com/ddddddO/sensor-pi/" // raspberry pi
const baseDir = "/mnt/c/DEV/workspace/GO/src/github.com/ddddddO/sensor-pi/" // wsl
const plotterDir = "plotter"

// ref: https://github.com/gonum/plot/wiki/Example-plots#more-detailed-style-settings
func main() {
	p := plot.New()

	p.Title.Text = "pressure @around tama river"
	p.Y.Label.Text = "pressure"
	p.X.Label.Text = "date"
	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04"}

	p.Add(plotter.NewGrid())

	data, err := fetchData()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	points, err := points(data)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	lpLine, lpPoints, err := plotter.NewLinePoints(points)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	lpLine.Color = color.RGBA{G: 255, A: 255}
	lpPoints.Shape = draw.PlusGlyph{}
	lpPoints.Color = color.RGBA{R: 255, A: 255}

	p.Add(lpLine, lpPoints)
	// p.Legend.Add("line points", lpLine, lpPoints)

	var imageName = "pressure.png"
	if err := p.Save(4*vg.Inch, 4*vg.Inch, filepath.Join(baseDir, plotterDir, imageName)); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type datum struct {
	t     time.Time
	value float64
}

func fetchData() ([]datum, error) {
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
	data := []datum{}
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
		datum := datum{t: tm, value: p} // NOTE: 一先ず気圧だけ
		data = append(data, datum)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func points(data []datum) (plotter.XYs, error) {
	pts := make(plotter.XYs, len(data))
	for i := range data {
		pts[i].X = float64(data[i].t.Unix())
		pts[i].Y = data[i].value
	}
	return pts, nil
}

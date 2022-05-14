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

var dir = "/home/pi/github.com/ddddddO/sensor-pi/" // raspberry pi
// var dir = "/mnt/c/DEV/workspace/GO/src/github.com/ddddddO/sensor-pi/environment.sqlite3" // wsl

// FIXME: 日時の取り扱い
// ref: https://github.com/gonum/plot/wiki/Example-plots#more-detailed-style-settings
func main() {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// TODO: 使い方あってないよね？
	plot.UTCUnixTime = plot.UnixTimeIn(loc)

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

	if err := p.Save(4*vg.Inch, 4*vg.Inch, filepath.Join(dir, "plotter", "pressure.png")); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type datum struct {
	t     time.Time
	value float64
}

func fetchData() ([]datum, error) {
	var dsn = filepath.Join(dir, "environment.sqlite3")
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
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, err
	}
	data := []datum{}
	for rows.Next() {
		var (
			d       string
			t, p, h float64
		)
		if err := rows.Scan(&d, &t, &p, &h); err != nil {
			return nil, err
		}

		// NOTE: 一先ず気圧だけ
		tm, err := time.ParseInLocation(layout, d, loc)
		if err != nil {
			return nil, err
		}
		datum := datum{t: tm, value: p}
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

package main

import (
	"fmt"
	"image/color"
	"os"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// FIXME: 日時の取り扱い
// ref: https://github.com/gonum/plot/wiki/Example-plots#more-detailed-style-settings
func main() {
	p := plot.New()

	p.Title.Text = "pressure @around tama river"
	p.Y.Label.Text = "pressure"
	p.X.Label.Text = "date"
	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04"}
	p.Add(plotter.NewGrid())

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	data := []datum{
		{t: time.Date(2022, 05, 12, 9, 00, 00, 0, loc), value: 1000.11},
		{t: time.Date(2022, 05, 12, 18, 00, 00, 0, loc), value: 992.11},
		{t: time.Date(2022, 05, 13, 9, 00, 00, 0, loc), value: 1013.99},
		{t: time.Date(2022, 05, 13, 18, 00, 00, 0, loc), value: 992.11},
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
	lpLine.Color = color.RGBA{R: 38, G: 205, B: 29}
	lpPoints.Shape = draw.PlusGlyph{}
	lpPoints.Color = color.RGBA{R: 255, A: 255}

	p.Add(lpLine, lpPoints)
	// p.Legend.Add("line points", lpLine, lpPoints)

	if err := p.Save(4*vg.Inch, 4*vg.Inch, "points.png"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type datum struct {
	t     time.Time
	value float64
}

func points(data []datum) (plotter.XYs, error) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, err
	}
	// TODO: 使い方あってないよね？
	plot.UTCUnixTime = plot.UnixTimeIn(loc)

	pts := make(plotter.XYs, len(data))
	for i := range data {
		pts[i].X = float64(data[i].t.Unix())
		pts[i].Y = data[i].value
	}
	return pts, nil
}

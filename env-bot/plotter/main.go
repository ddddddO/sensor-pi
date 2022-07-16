package main

import (
	"fmt"
	"image/color"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

// ref: https://github.com/gonum/plot/wiki/Example-plots#more-detailed-style-settings
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	environment, err := fetchData()
	if err != nil {
		return err
	}

	pressurePoints, err := points(environment.pressure)
	if err != nil {
		return err
	}
	pressurePlot := newPressurePlot()
	if err := pressurePlot.build(pressurePoints); err != nil {
		return err
	}
	if err := pressurePlot.save(); err != nil {
		return err
	}

	temperaturePoints, err := points(environment.temperature)
	if err != nil {
		return err
	}
	temperaturePlot := newTemperaturePlot()
	if err := temperaturePlot.build(temperaturePoints); err != nil {
		return err
	}
	if err := temperaturePlot.save(); err != nil {
		return err
	}

	humidityPoints, err := points(environment.humidity)
	if err != nil {
		return err
	}
	humidityPlot := newHumidityPlot()
	if err := humidityPlot.build(humidityPoints); err != nil {
		return err
	}
	if err := humidityPlot.save(); err != nil {
		return err
	}

	return nil
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

type env struct {
	D string  `db:"date"`
	T float64 `db:"temperature"`
	P float64 `db:"pressure"`
	H float64 `db:"humidity"`
}

func fetchData() (*environment, error) {
	var dsn = os.Getenv("DSN")
	db, err := sqlx.Connect("sqlite3", dsn)
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

	rows, err := db.Queryx(sql, args...)
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
		e := &env{}
		if err := rows.StructScan(e); err != nil {
			return nil, err
		}
		tm, err := time.Parse(layout, e.D)
		if err != nil {
			return nil, err
		}

		dataP = append(dataP, datum{t: tm, value: e.P})
		dataT = append(dataT, datum{t: tm, value: e.T})
		dataH = append(dataH, datum{t: tm, value: e.H})
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

	return &Plot{
		Plot:        p,
		imagePath:   os.Getenv("PRESSURE_IMAGE_PATH"),
		lineColor:   color.RGBA{G: 255, A: 255},
		pointsColor: color.RGBA{R: 255, A: 255},
	}
}

func newTemperaturePlot() *Plot {
	p := plot.New()
	p.Title.Text = "temperature @around tama river"
	p.Y.Label.Text = "temperature"
	p.X.Label.Text = "date"
	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04"}

	return &Plot{
		Plot:        p,
		imagePath:   os.Getenv("TEMPERATURE_IMAGE_PATH"),
		lineColor:   color.RGBA{R: 255, B: 255, A: 255},
		pointsColor: color.RGBA{R: 255, A: 255},
	}
}

func newHumidityPlot() *Plot {
	p := plot.New()
	p.Title.Text = "humidity @around tama river"
	p.Y.Label.Text = "humidity"
	p.X.Label.Text = "date"
	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04"}

	return &Plot{
		Plot:        p,
		imagePath:   os.Getenv("HUMIDITY_IMAGE_PATH"),
		lineColor:   color.RGBA{G: 255, B: 255, A: 255},
		pointsColor: color.RGBA{R: 255, A: 255},
	}
}

func (p *Plot) build(points plotter.XYs) error {
	p.Add(plotter.NewGrid())
	lpLine, lpPoints, err := plotter.NewLinePoints(points)
	if err != nil {
		return err
	}
	lpLine.Color = p.lineColor
	lpPoints.Shape = draw.PlusGlyph{}
	lpPoints.Color = p.pointsColor
	p.Add(lpLine, lpPoints)
	// p.Legend.Add("line points", lpLine, lpPoints)

	return nil
}

func (p *Plot) save() error {
	return p.Save(4*vg.Inch, 4*vg.Inch, p.imagePath)
}

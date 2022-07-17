package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(test)
}

type environment struct {
	Pressure    []datum `json:"pressure"`
	Temperature []datum `json:"temperature"`
	Humidity    []datum `json:"humidity"`
}

type datum struct {
	T     time.Time `json:"date"`
	Value float64   `json:"value"`
}

// TODO: rename
func test(ctx context.Context, snsEvent events.SNSEvent) {
	var r io.Reader
	for _, record := range snsEvent.Records {
		snsRecord := record.SNS
		r = strings.NewReader(snsRecord.Message)
		break
	}

	env := &environment{}
	if err := json.NewDecoder(r).Decode(env); err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := generate(env); err != nil {
		fmt.Println(err.Error())
		return
	}

	encoded, err := encode()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(encoded); err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("%s\n", buf.String())
	fmt.Println("finish!!!!!")
}

func generate(environment *environment) error {
	pressurePoints, err := points(environment.Pressure)
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

	temperaturePoints, err := points(environment.Temperature)
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

	humidityPoints, err := points(environment.Humidity)
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

func points(data []datum) (plotter.XYs, error) {
	pts := make(plotter.XYs, len(data))
	for i := range data {
		pts[i].X = float64(data[i].T.Unix())
		pts[i].Y = data[i].Value
	}
	return pts, nil
}

type Plot struct {
	*plot.Plot
	imagePath              string
	lineColor, pointsColor color.Color
}

const storeDir = "/tmp"

var pressureImagePath = filepath.Join(storeDir, "pressure.png")

func newPressurePlot() *Plot {
	p := plot.New()
	p.Title.Text = "pressure @around tama river"
	p.Y.Label.Text = "pressure"
	p.X.Label.Text = "date"
	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04"}

	return &Plot{
		Plot:        p,
		imagePath:   pressureImagePath,
		lineColor:   color.RGBA{G: 255, A: 255},
		pointsColor: color.RGBA{R: 255, A: 255},
	}
}

var temperatureImagePath = filepath.Join(storeDir, "temperature.png")

func newTemperaturePlot() *Plot {
	p := plot.New()
	p.Title.Text = "temperature @around tama river"
	p.Y.Label.Text = "temperature"
	p.X.Label.Text = "date"
	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04"}

	return &Plot{
		Plot:        p,
		imagePath:   temperatureImagePath,
		lineColor:   color.RGBA{R: 255, B: 255, A: 255},
		pointsColor: color.RGBA{R: 255, A: 255},
	}
}

var humidityImagePath = filepath.Join(storeDir, "humidity.png")

func newHumidityPlot() *Plot {
	p := plot.New()
	p.Title.Text = "humidity @around tama river"
	p.Y.Label.Text = "humidity"
	p.X.Label.Text = "date"
	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04"}

	return &Plot{
		Plot:        p,
		imagePath:   humidityImagePath,
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

type encodedImage struct {
	Pressure    string `json:"pressure"`
	Temperature string `json:"temperature"`
	Humidity    string `json:"humidity"`
}

func encode() (*encodedImage, error) {
	enc := &encodedImage{}
	var err error

	enc.Pressure, err = encodeImageToBase64(pressureImagePath)
	if err != nil {
		return nil, err
	}

	enc.Temperature, err = encodeImageToBase64(temperatureImagePath)
	if err != nil {
		return nil, err
	}

	enc.Humidity, err = encodeImageToBase64(humidityImagePath)
	if err != nil {
		return nil, err
	}

	return enc, nil
}

func encodeImageToBase64(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	body, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(body)
	return encoded, nil
}

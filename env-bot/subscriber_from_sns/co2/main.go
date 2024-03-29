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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
	lambda.Start(test)
}

type environment struct {
	Pressure    []datum `json:"pressure"`
	Temperature []datum `json:"temperature"`
	Humidity    []datum `json:"humidity"`
	CO2         []datum `json:"co2"`
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

	ret := &ret{
		Environment: []*retEnv{
			&retEnv{Type: "co2", Latest: env.CO2[0], Encoded: encoded},
		},
	}

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(ret); err != nil {
		fmt.Println(err.Error())
		return
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	queueURL := "https://sqs.ap-northeast-1.amazonaws.com/820544363308/filedata_to_tweeter"
	_, err = svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Title": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String("The Whistler"),
			},
			"Author": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String("John Grisham"),
			},
			"WeeksOn": &sqs.MessageAttributeValue{
				DataType:    aws.String("Number"),
				StringValue: aws.String("6"),
			},
		},
		MessageBody: aws.String(buf.String()),
		QueueUrl:    &queueURL,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("finish!!!!!")
}

func generate(environment *environment) error {
	co2Points, err := points(environment.CO2)
	if err != nil {
		return err
	}
	co2Plot := newCO2Plot()
	if err := co2Plot.build(co2Points); err != nil {
		return err
	}
	if err := co2Plot.save(); err != nil {
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

var co2ImagePath = filepath.Join(storeDir, "co2.png")

func newCO2Plot() *Plot {
	p := plot.New()
	p.Title.Text = "co2 @around tama river"
	p.Y.Label.Text = "co2"
	p.X.Label.Text = "date"
	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04"}

	return &Plot{
		Plot:        p,
		imagePath:   co2ImagePath,
		lineColor:   color.RGBA{R: 128, G: 128, B: 128, A: 1},
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

type ret struct {
	Environment []*retEnv `json:"environment"`
}

type retEnv struct {
	Type    string `json:"type"`
	Latest  datum  `json:"latest"`
	Encoded string `json:"encoded"`
}

func encode() (string, error) {
	encoded, err := encodeImageToBase64(co2ImagePath)
	if err != nil {
		return "", err
	}
	return encoded, nil
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

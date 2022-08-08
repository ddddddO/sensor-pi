package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"

	sq "github.com/Masterminds/squirrel"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/jmoiron/sqlx"
)

// ref: https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/sns-example-publish.html
func main() {
	env, err := fetchData()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(env); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file. (~/.aws/credentials).
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := sns.New(sess)
	msg := buf.String()
	topicARN := "arn:aws:sns:ap-northeast-1:820544363308:dbdata_to_filegenerator"

	result, err := svc.Publish(&sns.PublishInput{
		Message:  &msg,
		TopicArn: &topicARN,
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(*result.MessageId)
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

type env struct {
	D string  `db:"date"`
	T float64 `db:"temperature"`
	P float64 `db:"pressure"`
	H float64 `db:"humidity"`
	C float64 `db:"co2"`
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

		dataP = append(dataP, datum{T: tm, Value: e.P})
		dataT = append(dataT, datum{T: tm, Value: e.T})
		dataH = append(dataH, datum{T: tm, Value: e.H})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	queryMHZ19 := sq.Select("date", "co2").
		From("mh_z19").
		OrderBy("date desc").
		Limit(limit)
	sqlMHZ19, argsMHZ19, err := queryMHZ19.ToSql()
	if err != nil {
		return nil, err
	}

	rowsMHZ19, err := db.Queryx(sqlMHZ19, argsMHZ19...)
	if err != nil {
		return nil, err
	}
	defer rowsMHZ19.Close()

	var (
		dataC = []datum{} // for co2
	)
	for rowsMHZ19.Next() {
		e := &env{}
		if err := rowsMHZ19.StructScan(e); err != nil {
			return nil, err
		}
		tm, err := time.Parse(layout, e.D)
		if err != nil {
			return nil, err
		}

		dataC = append(dataC, datum{T: tm, Value: e.C})
	}
	if err := rowsMHZ19.Err(); err != nil {
		return nil, err
	}
	e := &environment{
		Pressure:    dataP,
		Temperature: dataT,
		Humidity:    dataH,
		CO2:         dataC,
	}
	return e, nil
}

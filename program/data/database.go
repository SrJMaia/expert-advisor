package data

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/SrJMaia/expert-advisor/program/check"
	"github.com/SrJMaia/expert-advisor/program/table"
	"github.com/jmoiron/sqlx"
)

const (
	host     = "localhost"
	database = "prices"
	user     = "root"
	password = "root"
)

type MyValues struct {
	Time  sql.NullTime `db:"date_hour"`
	Open  float64      `db:"open_prices"`
	High  float64      `db:"high_prices"`
	Low   float64      `db:"low_prices"`
	Close float64      `db:"close_prices"`
}

func FetchData(myData *LayoutData) {

	table.StartHead()

	var start = time.Now()

	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3307)/%s?parseTime=true", user, password, host, database)

	var db, err = sqlx.Open("mysql", connectionString)
	check.MyCheckingError(err)
	defer db.Close()

	var fetchedValues []MyValues

	err = db.Select(&fetchedValues, "SELECT eurusd.date_hour, eurusd.open_prices, eurusd.high_prices, eurusd.low_prices, eurusd.close_prices FROM prices.eurusd")
	check.MyCheckingError(err)

	(*myData).Time = make([]time.Time, len(fetchedValues))
	(*myData).Open = make([]float64, len(fetchedValues))
	(*myData).High = make([]float64, len(fetchedValues))
	(*myData).Low = make([]float64, len(fetchedValues))
	(*myData).Close = make([]float64, len(fetchedValues))

	for i, v := range fetchedValues {
		(*myData).Time[i] = v.Time.Time
		(*myData).Open[i] = v.Open
		(*myData).High[i] = v.High
		(*myData).Low[i] = v.Low
		(*myData).Close[i] = v.Close
	}

	var elapsed = time.Since(start)

	table.StartBody("Successfully Fetched Data from Database.", elapsed)
}

func PullData(pathFile string, datetimeFromat string) {
	/*
		CSV Format must be Datetime, Open, High, Low, Close only, like cleaned metatrader
		Metatrader datetime forma: %%Y-%%m-%%d %%H:%%i:%%s
	*/

	var databaseConection = fmt.Sprintf("%s:%s@tcp(%s:3307)/%s", user, password, host, database)

	var db, err = sql.Open("mysql", databaseConection)
	check.MyCheckingError(err)

	defer db.Close()
	table.AnalysisBody("Database Connected", "white")

	table.StartHead()
	start := time.Now()

	csvFile, err := os.Open(pathFile)
	check.MyCheckingError(err)
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	check.MyCheckingError(err)

	fmt.Println(len(csvLines))

	var stringSize int
	var stringToSlicce string

	var openValue float64
	var highValue float64
	var lowValue float64
	var closeValue float64

	var allValues []string

	for i, v := range csvLines {
		if i < 1 {
			continue
		}

		openValue, err = strconv.ParseFloat(v[1], 64)
		check.MyCheckingError(err)

		highValue, err = strconv.ParseFloat(v[2], 64)
		check.MyCheckingError(err)

		lowValue, err = strconv.ParseFloat(v[3], 64)
		check.MyCheckingError(err)

		closeValue, err = strconv.ParseFloat(v[4], 64)
		check.MyCheckingError(err)

		stringToSlicce = fmt.Sprintf("(STR_TO_DATE('%s', '%s' ) ,%.5f, %.5f, %.5f, %.5f)", v[0], datetimeFromat, openValue, highValue, lowValue, closeValue)

		if stringSize+len(stringToSlicce) >= 4000000 {
			x := fmt.Sprintf("INSERT INTO prices.eurusd(eurusd.date_hour, eurusd.open_prices, eurusd.high_prices, eurusd.low_prices, eurusd.close_prices) VALUES %s", strings.Join(allValues, ","))
			insert, err := db.Query(x)
			check.MyCheckingError(err)
			insert.Close()
			allValues = nil
			stringSize = 0
		}

		allValues = append(allValues, stringToSlicce)
		stringSize += len(stringToSlicce)

		if i == len(csvLines)-1 {
			x := fmt.Sprintf("INSERT INTO prices.eurusd(eurusd.date_hour, eurusd.open_prices, eurusd.high_prices, eurusd.low_prices, eurusd.close_prices) VALUES %s", strings.Join(allValues, ","))
			insert, err := db.Query(x)
			check.MyCheckingError(err)
			insert.Close()
			allValues = nil
			stringSize = 0
		}
	}

	elapsed := time.Since(start)

	table.StartBody("Successfully Pulled Data to Database.", elapsed)
}

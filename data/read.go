package data

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"

	"github.com/SrJMaia/expert-advisor/check"
	"github.com/SrJMaia/expert-advisor/table"
)

func ReadData(filePath string, myData *LayoutData) {
	table.StartHead()
	start := time.Now()

	var timeValue time.Time
	var openValue float64
	var highValue float64
	var lowValue float64
	var closeValue float64

	csvFile, err := os.Open(filePath)
	check.MyCheckingError(err)
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	check.MyCheckingError(err)

	var openSlice, highSlice = make([]float64, len(csvLines)-1), make([]float64, len(csvLines)-1)
	var lowSlice, closeSlice = make([]float64, len(csvLines)-1), make([]float64, len(csvLines)-1)
	var timeSlice = make([]time.Time, len(csvLines)-1)

	for i, v := range csvLines {
		if i < 1 {
			continue
		}

		timeValue, err = time.Parse("2006-01-02 15:04:05", v[0])
		check.MyCheckingError(err)
		timeSlice[i-1] = timeValue

		openValue, err = strconv.ParseFloat(v[1], 64)
		check.MyCheckingError(err)
		openSlice[i-1] = openValue

		highValue, err = strconv.ParseFloat(v[2], 64)
		check.MyCheckingError(err)
		highSlice[i-1] = highValue

		lowValue, err = strconv.ParseFloat(v[3], 64)
		check.MyCheckingError(err)
		lowSlice[i-1] = lowValue

		closeValue, err = strconv.ParseFloat(v[4], 64)
		check.MyCheckingError(err)
		closeSlice[i-1] = closeValue

	}

	(*myData).Time = timeSlice
	(*myData).Open = openSlice
	(*myData).High = highSlice
	(*myData).Low = lowSlice
	(*myData).Close = closeSlice

	elapsed := time.Since(start)

	table.StartBody("Successfully Read CSV file.", elapsed)
}

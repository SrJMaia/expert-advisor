package data

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/SrJMaia/expert-advisor/program/conversion"
	"github.com/SrJMaia/expert-advisor/program/table"
)

type LayoutData struct {
	Time          []time.Time
	TimeHour      []int
	TimeMinutes   []int
	TimeWeekDays  []int
	Open          []float64
	High          []float64
	Low           []float64
	Close         []float64
	PriceTf       []float64
	Tpsl          []float64
	BuyFlag       []bool
	SellFlag      []bool
	SizeBuy       uint32
	SizeSell      uint32
	SizeTimeFrame uint32
}

func CheckLayoutData(data *LayoutData) {
	start := time.Now()
	if data.Open == nil {
		log.Fatal("Data Open is nil")
	} else if data.High == nil {
		log.Fatal("Data High is nil")
	} else if data.Low == nil {
		log.Fatal("Data Low is nil")
	} else if data.Close == nil {
		log.Fatal("Data Close is nil")
	} else if data.BuyFlag == nil {
		log.Fatal("Data BuyFlag is nil")
	} else if data.PriceTf == nil {
		log.Fatal("Data PriceTf is nil")
	} else if data.SellFlag == nil {
		log.Fatal("Data SellFlag is nil")
	} else if data.Time == nil {
		log.Fatal("Data Time is nil")
	} else if data.Tpsl == nil {
		log.Fatal("Data Tpsl is nil")
	}

	for i := range data.Time {
		if data.Time[i].IsZero() {
			log.Fatal("Data Time is Zero ")
		} else if math.IsInf(data.Open[i], 0) || math.IsNaN(data.Open[i]) {
			log.Fatal("Data Open is Inf or NaN ")
		} else if math.IsInf(data.High[i], 0) || math.IsNaN(data.High[i]) {
			log.Fatal("Data High is Inf or NaN ")
		} else if math.IsInf(data.Low[i], 0) || math.IsNaN(data.Low[i]) {
			log.Fatal("Data Low is Inf or NaN ")
		} else if math.IsInf(data.Close[i], 0) || math.IsNaN(data.Close[i]) {
			log.Fatal("Data Close is Inf or NaN ")
		} else if math.IsInf(data.PriceTf[i], 0) || math.IsNaN(data.PriceTf[i]) {
			log.Fatal("Data PriceTf is Inf or NaN ")
		} else if math.IsInf(data.Tpsl[i], 0) || math.IsNaN(data.Tpsl[i]) {
			log.Fatal("Data.Tpsl is Inf or NaN ")
		}
	}
	elapsed := time.Since(start)
	table.StartBody("Successfully Checked Layout Data.", elapsed)
}

func ColumnsNames(x ...float64) string {
	var sliceStrings = make([]string, len(x))
	for i := range x {
		ret := fmt.Sprintf("%.2f", x[i])
		sliceStrings[i] = ret
	}
	var join = strings.Join(sliceStrings, "-")
	return join
}

func NormalizePricesTFAndDateTime(myData *LayoutData, timeFrame string) {
	var start = time.Now()
	(*myData).PriceTf = make([]float64, len((*myData).Open))
	(*myData).TimeHour = make([]int, len((*myData).Open))
	(*myData).TimeMinutes = make([]int, len((*myData).Open))
	(*myData).TimeWeekDays = make([]int, len((*myData).Open))

	for i := range (*myData).Open {
		(*myData).TimeHour[i] = myData.Time[i].Hour()
		(*myData).TimeMinutes[i] = myData.Time[i].Minute()
		(*myData).TimeWeekDays[i] = int(myData.Time[i].Weekday())
		if i == 0 {
			(*myData).PriceTf[i] = (*myData).Open[i]
			(*myData).SizeTimeFrame++
			continue
		}
		if timeFrame == "H1" && (*myData).Time[i].Hour() != (*myData).Time[i-1].Hour() {
			(*myData).PriceTf[i] = (*myData).Open[i]
			(*myData).SizeTimeFrame++
		} else if timeFrame == "D1" && (*myData).Time[i].Day() != (*myData).Time[i-1].Day() {
			(*myData).PriceTf[i] = (*myData).Open[i]
			(*myData).SizeTimeFrame++
		} else {
			(*myData).PriceTf[i] = (*myData).PriceTf[i-1]
		}
	}
	var elapsed = time.Since(start)
	table.StartBody("Successfully Created Price TF and Times.", elapsed)
}

func NormalizeTimeFrameFloat(smallSlice *[]float64, timeValues *[]time.Time, timeFrame string) ([]float64, int) {
	var normalizedSlice = make([]float64, len(*timeValues))
	var size int

	for i := range *smallSlice {
		if i == 0 {
			normalizedSlice[i] = (*smallSlice)[i]
			size++
			continue
		}
		if timeFrame == "H1" && (*timeValues)[i].Hour() != (*timeValues)[i-1].Hour() {
			normalizedSlice[i] = (*smallSlice)[i]
			size++
		} else if timeFrame == "D1" && (*timeValues)[i].Day() != (*timeValues)[i-1].Day() {
			normalizedSlice[i] = (*smallSlice)[i]
			size++
		} else {
			normalizedSlice[i] = normalizedSlice[i-1]
		}
	}

	return normalizedSlice, size
}

func NormalizeTimeFrameBool(smallSlice *[]bool, timeValues *[]time.Time, timeFrame string) []bool {
	var normalizedSlice = make([]bool, len(*timeValues))

	for i := range *smallSlice {
		if i == 0 {
			normalizedSlice[i] = (*smallSlice)[i]
			continue
		}
		if timeFrame == "H1" {
			if (*timeValues)[i].Hour() != (*timeValues)[i-1].Hour() {
				normalizedSlice[i] = (*smallSlice)[i]
			} else {
				normalizedSlice[i] = normalizedSlice[i-1]
			}
		} else if timeFrame == "D1" {
			if (*timeValues)[i].Day() != (*timeValues)[i-1].Day() {
				normalizedSlice[i] = (*smallSlice)[i]
			} else {
				normalizedSlice[i] = normalizedSlice[i-1]
			}
		}
	}

	return normalizedSlice
}

func RawTimeFrame(largerSlice *[]float64, timeValues *[]time.Time, timeFrame string) []float64 {
	var normalizedSlice = make([]float64, len(*timeValues))
	var realIndex = 1

	for i := range *timeValues {
		if i == 0 {
			normalizedSlice[0] = (*largerSlice)[i]
			continue
		}
		if timeFrame == "H1" {
			if (*timeValues)[i].Hour() != (*timeValues)[i-1].Hour() {
				normalizedSlice[realIndex] = (*largerSlice)[i]
				realIndex++
			}
		} else if timeFrame == "D1" {
			if (*timeValues)[i].Day() != (*timeValues)[i-1].Day() {
				normalizedSlice[realIndex] = (*largerSlice)[i]
				realIndex++
			}
		}

	}

	normalizedSlice = conversion.RemoveExcessZeros(normalizedSlice)

	return normalizedSlice
}

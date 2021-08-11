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
	Time              []time.Time
	TimeHour          []int
	TimeMinutes       []int
	TimeWeekDays      []int
	Open              []float64
	High              []float64
	Low               []float64
	Close             []float64
	PriceTf           []float64
	TpslNonFix        []float64
	TestMax           []float64
	TestMin           []float64
	BuyFlag           []bool
	SellFlag          []bool
	TpslFix           float64
	MultiplyTp        float64
	MultiplySl        float64
	ChannelInPips     float64
	StartPointInPips  float64
	MultiplyTpChannel float64
	SizeBuy           uint32
	SizeSell          uint32
	SizeTimeFrame     uint32
	IsFixTpsl         bool
}

func CheckLayoutData(data *LayoutData, optimationType string) {
	start := time.Now()
	if data.Open == nil {
		log.Panic("Data Open is nil")
	} else if data.High == nil {
		log.Panic("Data High is nil")
	} else if data.Low == nil {
		log.Panic("Data Low is nil")
	} else if data.Close == nil {
		log.Panic("Data Close is nil")
	} else if data.BuyFlag == nil && optimationType != "fimathe" && optimationType != "fimathe-with-strategy" {
		log.Panic("Data BuyFlag is nil")
	} else if data.PriceTf == nil {
		log.Panic("Data PriceTf is nil")
	} else if data.SellFlag == nil && optimationType != "fimathe" && optimationType != "fimathe-with-strategy" {
		log.Panic("Data SellFlag is nil")
	} else if data.Time == nil {
		log.Panic("Data Time is nil")
	}
	if data.IsFixTpsl {
		if data.TpslFix == 0 {
			log.Panic("Data TpslFix is zero")
		}
	} else {
		if data.TpslNonFix == nil {
			log.Panic("Data TpslNonFix is nil")
		}
	}

	for i := range data.Time {
		if data.Time[i].IsZero() {
			log.Panic("Data Time is Zero ")
		} else if math.IsInf(data.Open[i], 0) || math.IsNaN(data.Open[i]) {
			log.Panic("Data Open is Inf or NaN ")
		} else if math.IsInf(data.High[i], 0) || math.IsNaN(data.High[i]) {
			log.Panic("Data High is Inf or NaN ")
		} else if math.IsInf(data.Low[i], 0) || math.IsNaN(data.Low[i]) {
			log.Panic("Data Low is Inf or NaN ")
		} else if math.IsInf(data.Close[i], 0) || math.IsNaN(data.Close[i]) {
			log.Panic("Data Close is Inf or NaN ")
		} else if math.IsInf(data.PriceTf[i], 0) || math.IsNaN(data.PriceTf[i]) {
			log.Panic("Data PriceTf is Inf or NaN ")
		}
		if !data.IsFixTpsl {
			if math.IsInf(data.TpslNonFix[i], 0) || math.IsNaN(data.TpslNonFix[i]) {
				log.Panic("Data.TpslNonFix is Inf or NaN ")
			}
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

	for i := range *timeValues {
		if i == 0 {
			normalizedSlice[i] = (*smallSlice)[size]
			size++
			continue
		}
		if timeFrame == "H1" && (*timeValues)[i].Hour() != (*timeValues)[i-1].Hour() {
			normalizedSlice[i] = (*smallSlice)[size]
			size++
		} else if timeFrame == "D1" && (*timeValues)[i].Day() != (*timeValues)[i-1].Day() {
			normalizedSlice[i] = (*smallSlice)[size]
			size++
		} else {
			normalizedSlice[i] = normalizedSlice[i-1]
		}
	}

	if size != len(*smallSlice) {
		log.Panic("Normalize Time Frmae Float Has Wrong Size. Size: ", size, " Slice: ", len(*smallSlice))
	}

	return normalizedSlice, size
}

func NormalizeTimeFrameBool(smallSlice *[]bool, timeValues *[]time.Time, timeFrame string) []bool {
	var normalizedSlice = make([]bool, len(*timeValues))
	var indexSmallSlice uint32

	for i := range *timeValues {
		if i == 0 {
			normalizedSlice[i] = (*smallSlice)[indexSmallSlice]
			indexSmallSlice++
			continue
		}
		if timeFrame == "H1" && (*timeValues)[i].Hour() != (*timeValues)[i-1].Hour() {
			normalizedSlice[i] = (*smallSlice)[indexSmallSlice]
			indexSmallSlice++
		} else if timeFrame == "D1" && (*timeValues)[i].Day() != (*timeValues)[i-1].Day() {
			normalizedSlice[i] = (*smallSlice)[indexSmallSlice]
			indexSmallSlice++
		} else {
			normalizedSlice[i] = normalizedSlice[i-1]
		}
	}

	if int(indexSmallSlice) != len(*smallSlice) {
		log.Panic("Normalize Time Frmae Bool Has Wrong Size. Size: ", indexSmallSlice, " Slice: ", len(*smallSlice))
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

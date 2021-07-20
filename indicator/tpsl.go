package indicator

import (
	"math"
	"time"

	"github.com/SrJMaia/expert-advisor/conversion"
	"github.com/SrJMaia/expert-advisor/data"
	"github.com/SrJMaia/expert-advisor/mystats"
	"github.com/SrJMaia/expert-advisor/table"
)

func FixTpsl(myData *data.LayoutData, value float64, timeFrame string, jpy bool) {
	// Return with the size of dt
	start := time.Now()
	(*myData).Tpsl = make([]float64, len((*myData).Time))
	var realValue float64
	if jpy {
		realValue = conversion.Round(value/100, 3)
	} else {
		realValue = conversion.Round(value/10000, 5)
	}

	for i := range (*myData).Time {
		(*myData).Tpsl[i] = realValue
	}

	elapsed := time.Since(start)
	table.StartBody("Successfully Created TPSL.", elapsed)
}

func AtrTpsl(myData *data.LayoutData, averagePeriod int, timeFrame string, jpy bool) {
	start := time.Now()
	var closeValues []float64 = data.RawTimeFrame(&myData.Close, &myData.Time, timeFrame)
	var highValues = data.RawTimeFrame(&myData.High, &myData.Time, timeFrame)
	var lowValues = data.RawTimeFrame(&myData.Low, &myData.Time, timeFrame)
	var tpsl = make([]float64, len((*myData).Time))
	var realValue = make([]float64, len((*myData).Time))
	var averageValue float64
	var x float64
	var y float64
	var z float64

	for i := 0; i < len(closeValues)-2; i++ {
		averageValue = 0
		x = highValues[i+1] - lowValues[i+1]
		y = math.Abs(highValues[i+1] - closeValues[i+2])
		z = math.Abs(lowValues[i+1] - closeValues[i+2])
		realValue[i] = mystats.Max(x, y, z)
		if i < averagePeriod {
			continue
		}
		for j := 0; j < averagePeriod; j++ {
			averageValue += realValue[i-j]
		}
		averageValue = averageValue / float64(averagePeriod)
		if jpy {
			tpsl[i] = conversion.Round(averageValue, 3)
		} else {
			tpsl[i] = conversion.Round(averageValue, 5)
		}
	}

	var avg float64
	if jpy {
		avg = conversion.Round(mystats.Mean(tpsl), 3)
	} else {
		avg = conversion.Round(mystats.Mean(tpsl), 5)
	}

	tpsl = conversion.FillWithValue(tpsl, 0, avg)

	(*myData).Tpsl, _ = data.NormalizeTimeFrameFloat(&tpsl, &myData.Time, timeFrame)

	elapsed := time.Since(start)
	table.StartBody("Successfully Created TPSL.", elapsed)
}

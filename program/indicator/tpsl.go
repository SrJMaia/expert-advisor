package indicator

import (
	"math"
	"time"

	"github.com/SrJMaia/expert-advisor/program/conversion"
	"github.com/SrJMaia/expert-advisor/program/data"
	"github.com/SrJMaia/expert-advisor/program/mystats"
	"github.com/SrJMaia/expert-advisor/program/table"
)

func GetTPSL(myData *data.LayoutData, tpslType string, fixValue float64, averageAtrPeriod int, timeFrame string, jpy bool) {
	if tpslType == "fix" {
		if jpy {
			fixValue = fixValue / 100
		} else {
			fixValue = fixValue / 10000
		}
		(*myData).TpslFix = conversion.RoundIsJpy(fixValue, jpy)
		(*myData).IsFixTpsl = true
	} else if tpslType == "atr" {
		atrTpsl(myData, averageAtrPeriod, timeFrame, jpy)
		(*myData).IsFixTpsl = false
	}
}

func atrTpsl(myData *data.LayoutData, averagePeriod int, timeFrame string, jpy bool) {
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

		tpsl[i] = conversion.RoundIsJpy(averageValue, jpy)

	}

	var averageAtr = conversion.RoundIsJpy(mystats.Mean(tpsl), jpy)

	tpsl = conversion.FillWithValue(tpsl, 0, averageAtr)

	(*myData).TpslNonFix, _ = data.NormalizeTimeFrameFloat(&tpsl, &myData.Time, timeFrame)

	elapsed := time.Since(start)
	table.StartBody("Successfully Created TPSL.", elapsed)
}

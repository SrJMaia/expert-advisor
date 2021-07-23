package backtest

import (
	"fmt"
	"time"

	"github.com/SrJMaia/expert-advisor/program/data"
	"github.com/SrJMaia/expert-advisor/program/myanalysis"
	"github.com/SrJMaia/expert-advisor/program/table"
)

func stepsLoop(tf string) (float64, float64) {
	var stepOtimization float64
	var limitOtimization float64
	if tf == "H1" {
		stepOtimization = 0.5
		limitOtimization = 10.5
	} else if tf == "D1" {
		stepOtimization = 1.
		limitOtimization = 21.
	}
	return stepOtimization, limitOtimization
}

func Optimize(myData *data.LayoutData, myAnalysisStruct *myanalysis.LayoutAnalysis, optimizeType string, capital float64, jpy bool, saveCsv bool, leverage bool, tf string, indicator ...float64) {
	data.CheckLayoutData(myData)
	var tot []float64
	var buy []float64
	var sell []float64
	var totalBacktested uint32
	var columnsSlice []float64
	if optimizeType == "fimathe" {
		columnsSlice = make([]float64, len(indicator)+3)
		if len(columnsSlice) > 3 {
			for i := range indicator {
				columnsSlice[i+3] = indicator[i]
			}
		}
	} else {
		columnsSlice = make([]float64, len(indicator)+2)
		if len(columnsSlice) > 2 {
			for i := range indicator {
				columnsSlice[i+2] = indicator[i]
			}
		}
	}
	var stepOtimization, limitOtimization = stepsLoop(tf)

	var start = time.Now()
	if optimizeType == "notpsl" {
		var columns = data.ColumnsNames(indicator...)
		tot, buy, sell = NoTpslBacktest(myData, capital, jpy, leverage)
		myanalysis.AnalyseBacktest(tot, buy, sell, columns, myAnalysisStruct, int(myData.SizeTimeFrame), tf)
		if saveCsv {
			SaveBacktest(&tot, &buy, &sell, columns)
		}
	} else if optimizeType == "fimathe" {
		for channelInPips := 0.001; channelInPips < 0.021; channelInPips += 0.001 { // 0.021
			for startPointInPips := -channelInPips; startPointInPips <= channelInPips; startPointInPips += 0.001 {
				for multiplyTpChannel := 1.; multiplyTpChannel < 6.; multiplyTpChannel++ {
					table.LoopStepsPrint(channelInPips, startPointInPips, multiplyTpChannel)
					columnsSlice[0], columnsSlice[1], columnsSlice[2] = channelInPips, startPointInPips, multiplyTpChannel
					var columns = data.ColumnsNames(columnsSlice...)
					tot, buy, sell = Fimathe(myData, channelInPips, startPointInPips, multiplyTpChannel, capital, jpy, leverage)
					myanalysis.AnalyseBacktest(tot, buy, sell, columns, myAnalysisStruct, int(myData.SizeTimeFrame), tf)
					if saveCsv {
						SaveBacktest(&tot, &buy, &sell, columns)
					}
					totalBacktested++
				}
			}
		}

	} else {
		for i := stepOtimization; i < limitOtimization; i += stepOtimization {
			for j := stepOtimization; j <= i; j += stepOtimization {
				columnsSlice[0], columnsSlice[1] = i, j
				var columns = data.ColumnsNames(columnsSlice...)
				if optimizeType == "hedging" {
					tot, buy, sell = HedgingBacktest(myData, i, j, capital, jpy, leverage)
				} else if optimizeType == "netting" {
					tot, buy, sell = NettingBacktest(myData, i, j, capital, jpy, leverage)
				}
				myanalysis.AnalyseBacktest(tot, buy, sell, columns, myAnalysisStruct, int(myData.SizeTimeFrame), tf)
				if saveCsv {
					SaveBacktest(&tot, &buy, &sell, columns)
				}
				totalBacktested++
			}
		}
	}
	var elapsed = time.Since(start)
	if optimizeType == "notpsl" {
		table.StartBody("Successfully Backtested NoTPSL", elapsed)
	} else {
		table.StartBody(fmt.Sprint("Successfully Backtested ", totalBacktested, " times."), elapsed)
	}
}

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
	data.CheckLayoutData(myData, optimizeType)
	var tot []float64
	var buy []float64
	var sell []float64
	var totalBacktested uint32
	var columnsSlice []float64
	columnsSlice = make([]float64, len(indicator)+2)
	if len(columnsSlice) > 2 {
		for i := range indicator {
			columnsSlice[i+2] = indicator[i]
		}
	}

	var stepOtimization, limitOtimization = stepsLoop(tf)

	var start = time.Now()
	if optimizeType == "no-tpsl" {
		var columns = data.ColumnsNames(indicator...)
		tot, buy, sell = BacktestMain(myData, capital, optimizeType, jpy, leverage)
		myanalysis.AnalyseBacktest(tot, buy, sell, columns, myAnalysisStruct, int(myData.SizeTimeFrame), tf)
		if saveCsv {
			SaveBacktest(&tot, &buy, &sell, columns)
		}
	} else if optimizeType == "fimathe" {
		// For the best optimation should the startpoint be 0.0001
		for channelInPips := 0.001; channelInPips < 0.0105; channelInPips += 0.0005 {
			for startPointInPips := 0.000; startPointInPips < channelInPips; startPointInPips += 0.0001 {
				var optimationTime = time.Now()
				columnsSlice[0], columnsSlice[1] = channelInPips, startPointInPips
				var columns = data.ColumnsNames(columnsSlice...)
				(*myData).ChannelInPips = channelInPips
				(*myData).StartPointInPips = startPointInPips
				(*myData).MultiplyTpChannel = 2.
				tot, buy, sell = BacktestMain(myData, capital, optimizeType, jpy, leverage)
				myanalysis.AnalyseBacktest(tot, buy, sell, columns, myAnalysisStruct, int(myData.SizeTimeFrame), tf)
				if saveCsv {
					SaveBacktest(&tot, &buy, &sell, columns)
				}
				totalBacktested++
				var optimationElapsed = time.Since(optimationTime)
				table.LoopStepsPrint(optimationElapsed, "inside", channelInPips, startPointInPips)
			}
		}

	} else if optimizeType == "fimathe-with-strategy" {
		for channelInPips := 0.001; channelInPips < 0.021; channelInPips += 0.001 { // 0.021
			for startPointInPips := 0.000; startPointInPips < channelInPips; startPointInPips += 0.00001 {
				for multiplyTpChannel := 1.; multiplyTpChannel < 6.; multiplyTpChannel++ {
					var optimationTime = time.Now()
					columnsSlice[0], columnsSlice[1], columnsSlice[2] = channelInPips, startPointInPips, multiplyTpChannel
					var columns = data.ColumnsNames(columnsSlice...)
					(*myData).ChannelInPips = channelInPips
					(*myData).StartPointInPips = startPointInPips
					(*myData).MultiplyTpChannel = multiplyTpChannel
					tot, buy, sell = BacktestMain(myData, capital, optimizeType, jpy, leverage)
					myanalysis.AnalyseBacktest(tot, buy, sell, columns, myAnalysisStruct, int(myData.SizeTimeFrame), tf)
					if saveCsv {
						SaveBacktest(&tot, &buy, &sell, columns)
					}
					totalBacktested++
					var optimationElapsed = time.Since(optimationTime)
					table.LoopStepsPrint(optimationElapsed, "inside", channelInPips, startPointInPips, multiplyTpChannel)
				}
			}
		}
	} else {
		for i := stepOtimization; i < limitOtimization; i += stepOtimization {
			for j := stepOtimization; j <= i; j += stepOtimization {
				var optimationTime = time.Now()
				columnsSlice[0], columnsSlice[1] = i, j
				var columns = data.ColumnsNames(columnsSlice...)
				(*myData).MultiplyTp = i
				(*myData).MultiplySl = j
				if optimizeType == "hedging" {
					tot, buy, sell = BacktestMain(myData, capital, optimizeType, jpy, leverage)
				} else if optimizeType == "netting" {
					tot, buy, sell = BacktestMain(myData, capital, optimizeType, jpy, leverage)
				}
				myanalysis.AnalyseBacktest(tot, buy, sell, columns, myAnalysisStruct, int(myData.SizeTimeFrame), tf)
				if saveCsv {
					SaveBacktest(&tot, &buy, &sell, columns)
				}
				totalBacktested++
				var optimationElapsed = time.Since(optimationTime)
				table.LoopStepsPrint(optimationElapsed, "inside")
				table.LoopStepsPrint(optimationElapsed, "inside", i, j)
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

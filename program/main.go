package main

import (
	"time"

	"github.com/SrJMaia/expert-advisor/program/backtest"
	"github.com/SrJMaia/expert-advisor/program/data"
	"github.com/SrJMaia/expert-advisor/program/indicator"
	"github.com/SrJMaia/expert-advisor/program/myanalysis"
	"github.com/SrJMaia/expert-advisor/program/table"
	_ "github.com/go-sql-driver/mysql"
)

const (
	jpy             = false
	leverage        = false
	saveCsvBacktest = false
	saveCsvAnalysis = false
	tf              = "H1"
	capital         = 1000.
)

func main() {

	var myAnalysisStruct = myanalysis.LayoutAnalysis{}
	var myData = data.LayoutData{}

	data.FetchData(&myData)

	data.NormalizePricesTFAndDateTime(&myData, tf)

	indicator.GetTPSL(&myData, "fix", 10, 14, tf, jpy)

	var start = time.Now()
	/*
		Maybe expand the range of period to 100?
		or maybe decrease to 20 ou 30?

		Invert hpperiod1 and hpperiod2
	*/

	for period := 2; period < 52; period += 2 {
		for cutoff := 30; cutoff < 390; cutoff += 30 {
			var optimizationTime = time.Now()
			indicator.HigherLower(&myData, float64(period), float64(cutoff), tf, jpy)
			if myData.SizeBuy < 2 || myData.SizeSell < 2 {
				continue
			}
			backtest.Optimize(&myData, &myAnalysisStruct, "netting", capital, jpy, saveCsvBacktest, leverage, tf, float64(period), float64(cutoff))
			var optimizationElapsed = time.Since(optimizationTime)
			table.LoopStepsPrint(optimizationElapsed, "outside", float64(period), float64(cutoff))
		}
	}

	// for slow := 60; slow < 300; slow += 60 {
	// 	indicator.StrategyOneDecyclerPriceCross(&myData, float64(slow), tf, jpy)
	// 	if myData.SizeBuy < 2 || myData.SizeSell < 2 {
	// 		continue
	// 	}
	// 	backtest.Optimize(&myData, &myAnalysisStruct, "netting", capital, jpy, saveCsvBacktest, leverage, tf, float64(slow))
	// }

	// backtest.Optimize(&myData, &myAnalysisStruct, "fimathe", capital, jpy, saveCsvBacktest, leverage, tf)

	var elapsed = time.Since(start)
	table.StartBody("Successfully Backtested.", elapsed)
	myanalysis.FilterOtimization(&myAnalysisStruct, capital, float64(myData.SizeTimeFrame), saveCsvAnalysis)
}

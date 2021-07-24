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
	// for slow := 60; slow < 390; slow += 30 {
	// 	//for fast := 30; fast < slow; fast += 30 {
	// 	//for i := 90; i < 180; i += 30 {
	// 	table.LoopStepsPrint(float64(slow))
	// 	indicator.StrategyOneDecyclerPriceCross(&myData, float64(slow), tf, jpy)
	// 	if myData.SizeBuy < 2 || myData.SizeSell < 2 {
	// 		continue
	// 	}
	// 	backtest.Optimize(&myData, &myAnalysisStruct, "hedging", capital, jpy, saveCsvBacktest, leverage, tf, float64(slow))
	// 	//}
	// 	//}
	// }

	for slow := 60; slow < 300; slow += 60 {
		indicator.StrategyOneDecyclerPriceCross(&myData, float64(slow), tf, jpy)
		if myData.SizeBuy < 2 || myData.SizeSell < 2 {
			continue
		}
		backtest.Optimize(&myData, &myAnalysisStruct, "fimathe", capital, jpy, saveCsvBacktest, leverage, tf, float64(slow))
	}

	// indicator.StrategyOneDecyclerPriceCross(&myData, 30., tf, jpy)
	// backtest.Optimize(&myData, &myAnalysisStruct, "fimathe", capital, jpy, saveCsvBacktest, leverage, tf)

	var elapsed = time.Since(start)
	table.StartBody("Successfully Backtested.", elapsed)
	myanalysis.FilterOtimization(&myAnalysisStruct, capital, float64(myData.SizeTimeFrame), saveCsvAnalysis)
}

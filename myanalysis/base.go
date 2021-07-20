package myanalysis

import (
	"fmt"
	"math"
	"os"

	"github.com/SrJMaia/expert-advisor/conversion"
	"github.com/SrJMaia/expert-advisor/mystats"
	"github.com/SrJMaia/expert-advisor/table"
)

type LayoutAnalysis struct {
	Paramaters              []string
	TotalTrades             []float64
	TotalLongTrades         []float64
	TotalShortTrades        []float64
	AvgTradesYear           []float64
	FinalResult             []float64
	NetProfit               []float64
	GrossProfit             []float64
	GrossLoss               []float64
	PercentReturn           []float64
	PercentAnnualizedReturn []float64
	MaxDrawdownPercent      []float64
	MaxDrawdownValue        []float64
	WinRate                 []float64
	WinCount                []float64
	AvgWin                  []float64
	LargestWin              []float64
	SmallestWin             []float64
	LossRate                []float64
	LossCount               []float64
	AvgLoss                 []float64
	LargestLoss             []float64
	SmallestLoss            []float64
	RiskRuin                []float64
	ShaperRatio             []float64
}

type filterVariables struct {
	indexSelection []int
	befoneCut      int
	afterCut       int
}

// Futures implemetations
// fmt.Println("Paramaters:", myAnalysisStruct.Paramaters)
// fmt.Println("Gross Profit:", myAnalysisStruct.GrossProfit)
// fmt.Println("Gross Loss:", myAnalysisStruct.GrossLoss)
// fmt.Println("Percent Return:", myAnalysisStruct.PercentReturn)

// fmt.Println("Annualized Return:", myAnalysisStruct.PercentAnnualizedReturn)

// fmt.Println("Maximum Drawdown Percent:", myAnalysisStruct.MaxDrawDrownPercent)
// fmt.Println("Maximum Drawdown Value:", myAnalysisStruct.MaxDrawDrownValue)
// fmt.Println("Win Rate:", myAnalysisStruct.WinRate)
// fmt.Println("Win Count:", myAnalysisStruct.WinCount)
// fmt.Println("Avg Win:", myAnalysisStruct.AvgWin)
// fmt.Println("Largest Win:", myAnalysisStruct.LargestWin)
// fmt.Println("Smallest Win:", myAnalysisStruct.SmallestWin)
// fmt.Println("Loss Rate:", myAnalysisStruct.LossRate)
// fmt.Println("Loss Count:", myAnalysisStruct.LossCount)
// fmt.Println("Avg Loss:", myAnalysisStruct.AvgLoss)
// fmt.Println("Largest Loss:", myAnalysisStruct.LargestLoss)
// fmt.Println("Smallest Loss:", myAnalysisStruct.SmallestLoss)
// fmt.Println("Risk Of Ruin:", myAnalysisStruct.RiskRuin)
// fmt.Println("Shaper Ratio:", myAnalysisStruct.ShaperRatio)

func FilterOtimization(results *LayoutAnalysis, capital float64, sizeTf float64, save bool) {
	if save {
		SaveOtimization(results, "FullOtimization.csv", "Before Cut")
	}

	var mainVariables = filterVariables{}

	table.AnalysisHead()

	printResults(results)

	minimumTrades(&mainVariables, results, &sizeTf)

	UpdateAnalysisStruct(results, mainVariables.indexSelection)

	stopAnalysis(len(mainVariables.indexSelection))

	printResults(results)

	zeroNetProfit(&mainVariables, results)

	UpdateAnalysisStruct(results, mainVariables.indexSelection)

	stopAnalysis(len(mainVariables.indexSelection))

	printResults(results)

	annualizedRetOverX(&mainVariables, results, 5.)

	UpdateAnalysisStruct(results, mainVariables.indexSelection)

	stopAnalysis(len(mainVariables.indexSelection))

	printResults(results)

	if save {
		SaveOtimization(results, "CutedOtimization.csv", "After Cut")
	}
}

func stopAnalysis(leng int) {
	if leng == 0 {
		table.AnalysisBody("The Strategy did not pass.", "red")
		os.Exit(2)
	}
}

func printResults(results *LayoutAnalysis) {

	printDistribution("NET PROFIT", (*results).NetProfit)
	printDistribution("TOTAL TRADES", (*results).TotalTrades)
	printDistribution("ANNUALIZED RETURN", (*results).PercentAnnualizedReturn)
}

func annualizedRetOverX(mainVariables *filterVariables, results *LayoutAnalysis, cut float64) {
	var newIndexes []int
	for i, v := range (*results).PercentAnnualizedReturn {
		if v > cut {
			newIndexes = append(newIndexes, i)
		}
	}
	(*mainVariables).indexSelection = newIndexes
	(*mainVariables).afterCut = len((*mainVariables).indexSelection)
	table.AnalysisBody(
		fmt.Sprint("Annualized Return Cut - From ",
			(*mainVariables).befoneCut, " indexes, ",
			(*mainVariables).afterCut, " will continue."),
		"white")
	(*mainVariables).befoneCut = (*mainVariables).afterCut
}

func zeroNetProfit(mainVariables *filterVariables, results *LayoutAnalysis) {
	var cutZeroNetProfit = 0.
	var newIndexes []int
	for i, v := range (*results).NetProfit {
		if v > cutZeroNetProfit {
			newIndexes = append(newIndexes, i)
		}
	}
	(*mainVariables).indexSelection = newIndexes
	(*mainVariables).afterCut = len((*mainVariables).indexSelection)
	table.AnalysisBody(
		fmt.Sprint("Net Profit Below Zero Cut - From ",
			(*mainVariables).befoneCut, " indexes, ",
			(*mainVariables).afterCut, " will continue."),
		"white")
	(*mainVariables).befoneCut = (*mainVariables).afterCut
}

func minimumTrades(mainVariables *filterVariables, results *LayoutAnalysis, sizeTf *float64) {
	(*mainVariables).befoneCut = len(results.Paramaters)
	var cutTotalTrades = math.Sqrt(*sizeTf)
	for i, v := range results.TotalTrades {
		if v > cutTotalTrades {
			(*mainVariables).indexSelection = append((*mainVariables).indexSelection, i)
		}
	}
	(*mainVariables).afterCut = len((*mainVariables).indexSelection)
	table.AnalysisBody(
		fmt.Sprint("Minimum Trades Cut - From ", (*mainVariables).befoneCut,
			" indexes, ", (*mainVariables).afterCut, " will continue."),
		"white")
	(*mainVariables).befoneCut = (*mainVariables).afterCut
}

func printDistribution(name string, array []float64) {
	var minValue = conversion.Round(mystats.Min(array...), 2)
	var maxValue = conversion.Round(mystats.Max(array...), 2)
	var medianValue = conversion.Round(mystats.Median(array), 2)
	var averageValue = conversion.Round(mystats.Mean(array), 2)
	var q1Value, q3Value = mystats.Quartile(array)
	q1Value = conversion.Round(q1Value, 2)
	q3Value = conversion.Round(q3Value, 2)
	table.Distribution(name,
		minValue,
		q1Value,
		medianValue,
		averageValue,
		q3Value,
		maxValue,
		len(array))
}

func AnalyseBacktest(totalTrades []float64, buyTrades []float64, sellTrades []float64, paramaters string, myAnalysisStruct *LayoutAnalysis, lenPeriod int, tf string) {
	(*myAnalysisStruct).Paramaters = append((*myAnalysisStruct).Paramaters, paramaters)
	(*myAnalysisStruct).TotalTrades = append((*myAnalysisStruct).TotalTrades, float64(len(totalTrades)))
	var avgTradesY = AvgTradesYear(totalTrades, lenPeriod)
	(*myAnalysisStruct).AvgTradesYear = append((*myAnalysisStruct).AvgTradesYear, avgTradesY)
	(*myAnalysisStruct).TotalLongTrades = append((*myAnalysisStruct).TotalLongTrades, float64(len(buyTrades)))
	(*myAnalysisStruct).TotalShortTrades = append((*myAnalysisStruct).TotalShortTrades, float64(len(sellTrades)))
	(*myAnalysisStruct).FinalResult = append((*myAnalysisStruct).FinalResult, (totalTrades)[len(totalTrades)-1])
	var netProf = NetProfit(totalTrades)
	(*myAnalysisStruct).NetProfit = append((*myAnalysisStruct).NetProfit, netProf)
	var grossProfit, grossLoss = GrossProfitLoss(totalTrades)
	(*myAnalysisStruct).GrossProfit = append((*myAnalysisStruct).GrossProfit, grossProfit)
	(*myAnalysisStruct).GrossLoss = append((*myAnalysisStruct).GrossLoss, grossLoss)
	var perRet = PercentReturnCalc(totalTrades)
	(*myAnalysisStruct).PercentReturn = append((*myAnalysisStruct).PercentReturn, perRet)
	var annRet = AnnualizedReturn(totalTrades, lenPeriod, tf)
	(*myAnalysisStruct).PercentAnnualizedReturn = append((*myAnalysisStruct).PercentAnnualizedReturn, annRet)
	var maxDDP, maxDDV = MaximumDrawdown(totalTrades)
	(*myAnalysisStruct).MaxDrawdownPercent = append((*myAnalysisStruct).MaxDrawdownPercent, maxDDP)
	(*myAnalysisStruct).MaxDrawdownValue = append((*myAnalysisStruct).MaxDrawdownValue, maxDDV)
	var winRate, winCount = WinRateMean(totalTrades)
	(*myAnalysisStruct).WinRate = append((*myAnalysisStruct).WinRate, winRate)
	(*myAnalysisStruct).WinCount = append((*myAnalysisStruct).WinCount, winCount)
	var averageWin = AvgWin(totalTrades)
	(*myAnalysisStruct).AvgWin = append((*myAnalysisStruct).AvgWin, averageWin)
	var largerWin = LargeWin(totalTrades)
	(*myAnalysisStruct).LargestWin = append((*myAnalysisStruct).LargestWin, largerWin)
	var smallerWin = SmallWin(totalTrades)
	(*myAnalysisStruct).SmallestWin = append((*myAnalysisStruct).SmallestWin, smallerWin)
	var lossRate, lossCount = LossRateMean(totalTrades)
	(*myAnalysisStruct).LossRate = append((*myAnalysisStruct).LossRate, lossRate)
	(*myAnalysisStruct).LossCount = append((*myAnalysisStruct).LossCount, lossCount)
	var averageLoss = AvgLoss(totalTrades)
	(*myAnalysisStruct).AvgLoss = append((*myAnalysisStruct).AvgLoss, averageLoss)
	var largerLoss = LargeLoss(totalTrades)
	(*myAnalysisStruct).LargestLoss = append((*myAnalysisStruct).LargestLoss, largerLoss)
	var smallerLoss = SmallLoss(totalTrades)
	(*myAnalysisStruct).SmallestLoss = append((*myAnalysisStruct).SmallestLoss, smallerLoss)
	var riskRuin = RiskOfRuin(totalTrades)
	(*myAnalysisStruct).RiskRuin = append((*myAnalysisStruct).RiskRuin, riskRuin)
	var shaper = ShaperRatio(totalTrades, tf)
	(*myAnalysisStruct).ShaperRatio = append((*myAnalysisStruct).ShaperRatio, shaper)
}

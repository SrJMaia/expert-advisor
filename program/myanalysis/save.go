package myanalysis

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/SrJMaia/expert-advisor/program/check"
	"github.com/SrJMaia/expert-advisor/program/mystats"
	"github.com/SrJMaia/expert-advisor/program/table"
)

func SaveOtimization(results *LayoutAnalysis, fileName string, stage string) {

	start := time.Now()
	var avgLoss = make([]string, len((*results).AvgLoss)+1)
	var avgTradesYear = make([]string, len((*results).AvgTradesYear)+1)
	var avgWin = make([]string, len((*results).AvgWin)+1)
	var finalResult = make([]string, len((*results).FinalResult)+1)
	var grossLoss = make([]string, len((*results).GrossLoss)+1)
	var grossProfit = make([]string, len((*results).GrossProfit)+1)
	var largerLoss = make([]string, len((*results).LargestLoss)+1)
	var largerWin = make([]string, len((*results).LargestWin)+1)
	var lossCount = make([]string, len((*results).LossCount)+1)
	var lossRate = make([]string, len((*results).LossRate)+1)
	var maxDrawdownPercent = make([]string, len((*results).MaxDrawdownPercent)+1)
	var maxDrawdownValue = make([]string, len((*results).MaxDrawdownValue)+1)
	var netProfit = make([]string, len((*results).NetProfit)+1)
	var percentAnnualizedReturn = make([]string, len((*results).PercentAnnualizedReturn)+1)
	var percentReturn = make([]string, len((*results).PercentReturn)+1)
	var riskRuin = make([]string, len((*results).RiskRuin)+1)
	var shaperRatio = make([]string, len((*results).ShaperRatio)+1)
	var smallestLoss = make([]string, len((*results).SmallestLoss)+1)
	var smallestWin = make([]string, len((*results).SmallestWin)+1)
	var totalLongTrades = make([]string, len((*results).TotalLongTrades)+1)
	var totalShortTrades = make([]string, len((*results).TotalShortTrades)+1)
	var totalTrades = make([]string, len((*results).TotalTrades)+1)
	var winCount = make([]string, len((*results).WinCount)+1)
	var winRate = make([]string, len((*results).WinRate)+1)
	var paramaters = make([]string, len((*results).Paramaters)+1)

	var _, errCheck = os.Stat("C:/Users/johnk/Google Drive/Programming/Dados/" + fileName)
	if os.IsNotExist(errCheck) {
		var file, err = os.OpenFile("C:/Users/johnk/Google Drive/Programming/Dados/"+fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		var writer = csv.NewWriter(file)
		paramaters[0] = "0"
		for i := 1; i < len(paramaters); i++ {
			paramaters[i] = (*results).Paramaters[i-1]
		}

		err = writer.Write(paramaters)
		check.MyCheckingError(err)
		writer.Flush()
		file.Close()
	}

	var file, err = os.OpenFile("C:/Users/johnk/Google Drive/Programming/Dados/"+fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	check.MyCheckingError(err)
	defer file.Close()

	var writer = csv.NewWriter(file)
	defer writer.Flush()

	avgLoss[0] = "avgLoss"
	avgTradesYear[0] = "avgTradesYear"
	avgWin[0] = "avgWin"
	finalResult[0] = "finalResult"
	grossLoss[0] = "grossLoss"
	grossProfit[0] = "grossProfit"
	largerLoss[0] = "largerLoss"
	largerWin[0] = "largerWin"
	lossCount[0] = "lossCount"
	lossRate[0] = "lossRate"
	maxDrawdownPercent[0] = "maxDrawdownPercent"
	maxDrawdownValue[0] = "maxDrawdownValue"
	netProfit[0] = "netProfit"
	percentAnnualizedReturn[0] = "percentAnnualizedReturn"
	percentReturn[0] = "percentReturn"
	riskRuin[0] = "riskRuin"
	shaperRatio[0] = "shaperRatio"
	smallestLoss[0] = "smallestLoss"
	smallestWin[0] = "smallestWin"
	totalLongTrades[0] = "totalLongTrades"
	totalShortTrades[0] = "totalShortTrades"
	totalTrades[0] = "totalTrades"
	winCount[0] = "winCount"
	winRate[0] = "winRate"

	for i := 1; i < len(avgLoss); i++ {
		avgLoss[i] = strconv.FormatFloat((*results).AvgLoss[i-1], 'E', -1, 64)
		avgTradesYear[i] = strconv.FormatFloat((*results).AvgTradesYear[i-1], 'E', -1, 64)
		avgWin[i] = strconv.FormatFloat((*results).AvgWin[i-1], 'E', -1, 64)
		finalResult[i] = strconv.FormatFloat((*results).FinalResult[i-1], 'E', -1, 64)
		grossLoss[i] = strconv.FormatFloat((*results).GrossLoss[i-1], 'E', -1, 64)
		grossProfit[i] = strconv.FormatFloat((*results).GrossProfit[i-1], 'E', -1, 64)
		largerLoss[i] = strconv.FormatFloat((*results).LargestLoss[i-1], 'E', -1, 64)
		largerWin[i] = strconv.FormatFloat((*results).LargestWin[i-1], 'E', -1, 64)
		lossCount[i] = strconv.FormatFloat((*results).LossCount[i-1], 'E', -1, 64)
		lossRate[i] = strconv.FormatFloat((*results).LossRate[i-1], 'E', -1, 64)
		maxDrawdownPercent[i] = strconv.FormatFloat((*results).MaxDrawdownPercent[i-1], 'E', -1, 64)
		maxDrawdownValue[i] = strconv.FormatFloat((*results).MaxDrawdownValue[i-1], 'E', -1, 64)
		netProfit[i] = strconv.FormatFloat((*results).NetProfit[i-1], 'E', -1, 64)
		percentAnnualizedReturn[i] = strconv.FormatFloat((*results).PercentAnnualizedReturn[i-1], 'E', -1, 64)
		percentReturn[i] = strconv.FormatFloat((*results).PercentReturn[i-1], 'E', -1, 64)
		riskRuin[i] = strconv.FormatFloat((*results).RiskRuin[i-1], 'E', -1, 64)
		shaperRatio[i] = strconv.FormatFloat((*results).ShaperRatio[i-1], 'E', -1, 64)
		smallestLoss[i] = strconv.FormatFloat((*results).SmallestLoss[i-1], 'E', -1, 64)
		smallestWin[i] = strconv.FormatFloat((*results).SmallestWin[i-1], 'E', -1, 64)
		totalLongTrades[i] = strconv.FormatFloat((*results).TotalLongTrades[i-1], 'E', -1, 64)
		totalShortTrades[i] = strconv.FormatFloat((*results).TotalShortTrades[i-1], 'E', -1, 64)
		totalTrades[i] = strconv.FormatFloat((*results).TotalTrades[i-1], 'E', -1, 64)
		winCount[i] = strconv.FormatFloat((*results).WinCount[i-1], 'E', -1, 64)
		winRate[i] = strconv.FormatFloat((*results).WinRate[i-1], 'E', -1, 64)
	}

	err = writer.Write(avgLoss)
	check.MyCheckingError(err)
	err = writer.Write(avgTradesYear)
	check.MyCheckingError(err)
	err = writer.Write(avgWin)
	check.MyCheckingError(err)
	err = writer.Write(finalResult)
	check.MyCheckingError(err)
	err = writer.Write(grossLoss)
	check.MyCheckingError(err)
	err = writer.Write(grossProfit)
	check.MyCheckingError(err)
	err = writer.Write(largerLoss)
	check.MyCheckingError(err)
	err = writer.Write(largerWin)
	check.MyCheckingError(err)
	err = writer.Write(lossCount)
	check.MyCheckingError(err)
	err = writer.Write(lossRate)
	check.MyCheckingError(err)
	err = writer.Write(maxDrawdownPercent)
	check.MyCheckingError(err)
	err = writer.Write(maxDrawdownValue)
	check.MyCheckingError(err)
	err = writer.Write(netProfit)
	check.MyCheckingError(err)
	err = writer.Write(percentAnnualizedReturn)
	check.MyCheckingError(err)
	err = writer.Write(percentReturn)
	check.MyCheckingError(err)
	err = writer.Write(riskRuin)
	check.MyCheckingError(err)
	err = writer.Write(shaperRatio)
	check.MyCheckingError(err)
	err = writer.Write(smallestLoss)
	check.MyCheckingError(err)
	err = writer.Write(smallestWin)
	check.MyCheckingError(err)
	err = writer.Write(totalLongTrades)
	check.MyCheckingError(err)
	err = writer.Write(totalShortTrades)
	check.MyCheckingError(err)
	err = writer.Write(totalTrades)
	check.MyCheckingError(err)
	err = writer.Write(winCount)
	check.MyCheckingError(err)
	err = writer.Write(winRate)
	check.MyCheckingError(err)

	elapsed := time.Since(start)
	table.StartBody(fmt.Sprint("Successfully Saved Analysis", stage, "in CSV file."), elapsed)
}

func UpdateAnalysisStruct(myAnalysisStruct *LayoutAnalysis, indexes []int) {
	(*myAnalysisStruct).Paramaters = mystats.GetValuesByIndexString(&myAnalysisStruct.Paramaters, indexes)
	(*myAnalysisStruct).TotalTrades = mystats.GetValuesByIndexFloat(&myAnalysisStruct.TotalTrades, indexes)
	(*myAnalysisStruct).TotalLongTrades = mystats.GetValuesByIndexFloat(&myAnalysisStruct.TotalLongTrades, indexes)
	(*myAnalysisStruct).TotalShortTrades = mystats.GetValuesByIndexFloat(&myAnalysisStruct.TotalShortTrades, indexes)
	(*myAnalysisStruct).AvgTradesYear = mystats.GetValuesByIndexFloat(&myAnalysisStruct.AvgTradesYear, indexes)
	(*myAnalysisStruct).FinalResult = mystats.GetValuesByIndexFloat(&myAnalysisStruct.FinalResult, indexes)
	(*myAnalysisStruct).NetProfit = mystats.GetValuesByIndexFloat(&myAnalysisStruct.NetProfit, indexes)
	(*myAnalysisStruct).GrossProfit = mystats.GetValuesByIndexFloat(&myAnalysisStruct.GrossProfit, indexes)
	(*myAnalysisStruct).GrossLoss = mystats.GetValuesByIndexFloat(&myAnalysisStruct.GrossLoss, indexes)
	(*myAnalysisStruct).PercentReturn = mystats.GetValuesByIndexFloat(&myAnalysisStruct.PercentReturn, indexes)
	(*myAnalysisStruct).PercentAnnualizedReturn = mystats.GetValuesByIndexFloat(&myAnalysisStruct.PercentAnnualizedReturn, indexes)
	(*myAnalysisStruct).MaxDrawdownPercent = mystats.GetValuesByIndexFloat(&myAnalysisStruct.MaxDrawdownPercent, indexes)
	(*myAnalysisStruct).MaxDrawdownValue = mystats.GetValuesByIndexFloat(&myAnalysisStruct.MaxDrawdownValue, indexes)
	(*myAnalysisStruct).WinRate = mystats.GetValuesByIndexFloat(&myAnalysisStruct.WinRate, indexes)
	(*myAnalysisStruct).WinCount = mystats.GetValuesByIndexFloat(&myAnalysisStruct.WinCount, indexes)
	(*myAnalysisStruct).AvgWin = mystats.GetValuesByIndexFloat(&myAnalysisStruct.AvgWin, indexes)
	(*myAnalysisStruct).LargestWin = mystats.GetValuesByIndexFloat(&myAnalysisStruct.LargestWin, indexes)
	(*myAnalysisStruct).SmallestWin = mystats.GetValuesByIndexFloat(&myAnalysisStruct.SmallestWin, indexes)
	(*myAnalysisStruct).LossRate = mystats.GetValuesByIndexFloat(&myAnalysisStruct.LossRate, indexes)
	(*myAnalysisStruct).LossCount = mystats.GetValuesByIndexFloat(&myAnalysisStruct.LossCount, indexes)
	(*myAnalysisStruct).AvgLoss = mystats.GetValuesByIndexFloat(&myAnalysisStruct.AvgLoss, indexes)
	(*myAnalysisStruct).LargestLoss = mystats.GetValuesByIndexFloat(&myAnalysisStruct.LargestLoss, indexes)
	(*myAnalysisStruct).SmallestLoss = mystats.GetValuesByIndexFloat(&myAnalysisStruct.SmallestLoss, indexes)
	(*myAnalysisStruct).RiskRuin = mystats.GetValuesByIndexFloat(&myAnalysisStruct.RiskRuin, indexes)
	(*myAnalysisStruct).ShaperRatio = mystats.GetValuesByIndexFloat(&myAnalysisStruct.ShaperRatio, indexes)
}

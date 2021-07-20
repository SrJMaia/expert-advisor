package myanalysis

import (
	"math"

	"github.com/SrJMaia/expert-advisor/conversion"
	"github.com/SrJMaia/expert-advisor/mystats"
)

const constRisk = 0.5

func AvgTradesYear(trades []float64, lenPeriod int) float64 {
	// Here i'm using a full year instead of year working days +- 261
	var period = float64(lenPeriod) / 60 / 24
	var result = float64(len(trades)) / period * 365
	return conversion.Round(result, 2)
}

func NetProfit(trades []float64) float64 {
	return conversion.Round(trades[len(trades)-1]-trades[0], 2)
}

func GrossProfitLoss(trades []float64) (float64, float64) {
	var grossProfit float64
	var grossLoss float64
	trades = mystats.Diff(trades)
	for i := range trades {
		if trades[i] < 0 {
			grossLoss += trades[i]
		} else if trades[i] > 0 {
			grossProfit += trades[i]
		}
	}
	grossProfit = conversion.Round(grossProfit, 2)
	grossLoss = conversion.Round(grossLoss, 2)
	return grossProfit, grossLoss
}

func PercentReturnCalc(trades []float64) float64 {
	// Total return in percent %
	return conversion.Round(((trades[len(trades)-1]/trades[0])-1)*100, 2)
}

func WinRateMean(trades []float64) (float64, float64) {
	// The average win positions and the total win positions
	var WinRtValue float64
	trades = mystats.Diff(trades)
	for i := range trades {
		if (trades)[i] > 0 {
			WinRtValue++
		}
	}
	return conversion.Round((WinRtValue/float64(len(trades)))*100, 2), WinRtValue
}

func AvgWin(trades []float64) float64 {
	// The average win price
	trades = mystats.Diff(trades)
	var sumWin float64
	var sumInd float64
	for i := range trades {
		if trades[i] > 0 {
			sumWin += trades[i]
			sumInd++
		}
	}
	sumWin = sumWin / sumInd
	return conversion.Round(sumWin, 2)
}

func LargeWin(trades []float64) float64 {
	// The largest win return
	trades = mystats.Diff(trades)
	return conversion.Round(mystats.Max(trades...), 2)
}

func SmallWin(trades []float64) float64 {
	// The smallest win return
	trades = mystats.Diff(trades)
	return conversion.Round(mystats.MinAboveZero(trades...), 2)
}

func LossRateMean(trades []float64) (float64, float64) {
	// Avera of loss trades and the total loss trades
	var LossRtValue float64
	trades = mystats.Diff(trades)
	for i := range trades {
		if trades[i] < 0 {
			LossRtValue++
		}
	}
	return conversion.Round((LossRtValue/float64(len(trades)))*100, 2), LossRtValue
}

func AvgLoss(trades []float64) float64 {
	// Average loss price
	trades = mystats.Diff(trades)
	var sumLoss float64
	var sumInd float64
	for i := range trades {
		if trades[i] < 0 {
			sumLoss += (trades)[i]
			sumInd++
		}
	}
	sumLoss = sumLoss / sumInd
	return conversion.Round(sumLoss, 2)
}

func LargeLoss(trades []float64) float64 {
	// The largest loss
	trades = mystats.Diff(trades)
	return conversion.Round(mystats.Min(trades...), 2)
}

func SmallLoss(trades []float64) float64 {
	// The smallest loss
	trades = mystats.Diff(trades)
	return conversion.Round(mystats.MaxBelowZero(trades...), 2)
}

func AnnualizedReturn(trades []float64, lenPeriod int, tf string) float64 {
	var period float64
	var ret float64
	if tf == "H1" {
		period = 365 / (float64(lenPeriod) / 24)
	} else if tf == "D1" {
		period = 365 / float64(lenPeriod)
	}
	if trades[len(trades)-1] < 0 {
		ret = math.Abs(trades[len(trades)-1]-trades[0]) / trades[0]
		return conversion.Round((math.Pow(ret, period)-1)*-100, 2)
	} else if trades[len(trades)-1] >= 0 && trades[len(trades)-1] <= trades[0] {
		ret = 2 - (trades[len(trades)-1] / trades[0])
		return conversion.Round((math.Pow(ret, period)-1)*-100, 2)
	} else {
		ret = trades[len(trades)-1] / trades[0]
		return conversion.Round((math.Pow(ret, period)-1)*100, 2)
	}
}

func MaximumDrawdown(trades []float64) (float64, float64) {
	var maxAcc = mystats.MaximumAccumulate(trades)
	var maxDDSlice = make([]float64, len(trades))
	var maxDDPercent float64
	var maxDDValue float64
	var index int
	for i := range trades {
		maxDDSlice[i] = (maxAcc[i] - trades[i]) / maxAcc[i]
	}
	maxDDPercent = mystats.Max(maxDDSlice...)
	index = mystats.FindFirstIndexValue(&maxDDSlice, maxDDPercent)
	if index == 0 || len(trades) == 1 {
		return 999., 999.
	}
	maxDDValue = trades[index-1] - trades[index]
	return conversion.Round(maxDDPercent*-100, 2), conversion.Round(maxDDValue, 2)
}

func RiskOfRuin(trades []float64) float64 {
	var result float64
	var avgLoss = AvgLoss(trades)
	var probWin, _ = WinRateMean(trades)
	var probLoss, _ = LossRateMean(trades)
	probWin = probWin / 100
	probLoss = probLoss / 100
	var myRisk = trades[0] * constRisk / math.Abs(avgLoss)
	if probWin > probLoss {
		result = math.Pow((1-(probWin-probLoss))/(1+(probWin-probLoss)), myRisk)
	} else {
		result = math.Pow((1 + (probWin-probLoss)/(1-(probWin-probLoss))), myRisk)
	}
	return conversion.Round(result*100, 5)
}

func ShaperRatio(trades []float64, tf string) float64 {
	var period float64
	if tf == "H1" {
		period = math.Sqrt(365 * 24)
	} else if tf == "D1" {
		period = math.Sqrt(365)
	}
	trades = mystats.PctChange(trades)
	var avg = mystats.Mean(trades)
	var std = mystats.StandardDeviationSample(trades)
	var result = (avg / 100) / std
	return conversion.Round(result*period, 3)
}

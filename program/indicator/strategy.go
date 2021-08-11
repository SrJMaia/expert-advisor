package indicator

import (
	"fmt"
	"time"

	"github.com/SrJMaia/expert-advisor/program/data"
	"github.com/SrJMaia/expert-advisor/program/mystats"
	"github.com/SrJMaia/expert-advisor/program/table"
)

func HigherLower(myData *data.LayoutData, period float64, cutoff float64, tf string, jpy bool) {
	var start = time.Now()
	var newPrice = data.RawTimeFrame(&myData.Open, &myData.Time, tf)
	var newHigh = data.RawTimeFrame(&myData.High, &myData.Time, tf)
	var newLow = data.RawTimeFrame(&myData.Low, &myData.Time, tf)
	var buy = make([]bool, len(newPrice))
	var sell = make([]bool, len(newPrice))
	var decycler = Decycler(&newPrice, cutoff, jpy)
	var testMax = make([]float64, len(newPrice))
	var testMin = make([]float64, len(newPrice))
	var sizeBuy uint32
	var sizeSell uint32
	var lowestPrice float64
	var highestPrice float64
	for i := range newPrice {
		if i < int(period) {
			buy[i] = false
			sell[i] = false
			testMax[i] = 0.
			testMin[i] = 0.
			continue
		}
		highestPrice, _ = mystats.MaxMin(newHigh[i-int(period) : i])
		_, lowestPrice = mystats.MaxMin(newLow[i-int(period) : i])
		testMax[i] = highestPrice
		testMin[i] = lowestPrice
		if newPrice[i] < lowestPrice && newPrice[i-1] > lowestPrice && decycler[i] < newPrice[i] {
			sell[i] = true
			sizeSell++
		} else {
			sell[i] = false
		}
		if newPrice[i] > highestPrice && newPrice[i-1] < highestPrice && decycler[i] > newPrice[i] {
			buy[i] = true
			sizeBuy++
		} else {
			buy[i] = false
		}
	}
	(*myData).BuyFlag = data.NormalizeTimeFrameBool(&buy, &myData.Time, tf)
	(*myData).SellFlag = data.NormalizeTimeFrameBool(&sell, &myData.Time, tf)
	(*myData).TestMax, _ = data.NormalizeTimeFrameFloat(&testMax, &myData.Time, tf)
	(*myData).TestMin, _ = data.NormalizeTimeFrameFloat(&testMin, &myData.Time, tf)
	(*myData).SizeBuy = sizeBuy + 1
	(*myData).SizeSell = sizeSell + 1
	var elapsed = time.Since(start)
	table.StartBody(fmt.Sprint("Successfully Created Strategy - Higher Lower - ", "Buy: ", sizeBuy, " ", "Sell: ", sizeSell), elapsed)
}

func StrategyOneDecyclerPriceSingleCross(myData *data.LayoutData, cutoff float64, tf string, jpy bool) {
	// Heging - Fix TPSL 10  - TPSL = 19 2.5 - slow i=240 fast j=210 Mostraram resultados promissores
	// Hedging - ATR TPSL  shows more interesting
	// Netting didnt show anyhting good
	var start = time.Now()
	var newPrice = data.RawTimeFrame(&myData.Open, &myData.Time, tf)
	var decycler = Decycler(&newPrice, cutoff, jpy)
	var buy = make([]bool, len(newPrice))
	var sell = make([]bool, len(newPrice))
	var sizeBuy uint32
	var sizeSell uint32
	for i := range newPrice {
		if i < 1 {
			buy[i] = false
			sell[i] = false
			continue
		}
		if decycler[i] > newPrice[i] && decycler[i-1] < newPrice[i-1] {
			buy[i] = true
			sizeBuy++
		} else if decycler[i] < newPrice[i] && decycler[i-1] > newPrice[i-1] {
			sell[i] = true
			sizeSell++
		}
	}
	var newBuy = data.NormalizeTimeFrameBool(&buy, &myData.Time, tf)
	var newSell = data.NormalizeTimeFrameBool(&sell, &myData.Time, tf)
	(*myData).BuyFlag = newBuy
	(*myData).SellFlag = newSell
	(*myData).SizeBuy = sizeBuy + 1
	(*myData).SizeSell = sizeSell + 1
	var elapsed = time.Since(start)
	table.StartBody(fmt.Sprint("Successfully Created Strategy - Decycler Price Single Cross - ", "Buy: ", sizeBuy, " ", "Sell: ", sizeSell), elapsed)
}

func StrategyOneDecyclerPriceCross(myData *data.LayoutData, cutoff float64, tf string, jpy bool) {
	// Heging - Fix TPSL 10  - TPSL = 19 2.5 - slow i=240 fast j=210 Mostraram resultados promissores
	// Hedging - ATR TPSL  shows more interesting
	// Netting didnt show anyhting good
	var start = time.Now()
	var newPrice = data.RawTimeFrame(&myData.Open, &myData.Time, tf)
	var decyclerFast = Decycler(&newPrice, cutoff, jpy)
	var buy = make([]bool, len(newPrice))
	var sell = make([]bool, len(newPrice))
	var sizeBuy uint32
	var sizeSell uint32
	for i := range newPrice {
		if decyclerFast[i] > newPrice[i] {
			buy[i] = true
			sizeBuy++
		} else if decyclerFast[i] < newPrice[i] {
			sell[i] = true
			sizeSell++
		}
	}
	var newBuy = data.NormalizeTimeFrameBool(&buy, &myData.Time, tf)
	var newSell = data.NormalizeTimeFrameBool(&sell, &myData.Time, tf)
	(*myData).BuyFlag = newBuy
	(*myData).SellFlag = newSell
	(*myData).SizeBuy = sizeBuy + 1
	(*myData).SizeSell = sizeSell + 1
	var elapsed = time.Since(start)
	table.StartBody(fmt.Sprint("Successfully Created Strategy - Decycler Price Cross - ", "Buy: ", sizeBuy, " ", "Sell: ", sizeSell), elapsed)
}

func StrategyTwoDecyclerCrossWithDecyclerOscilator(myData *data.LayoutData, cutoffFast float64, cutoffSlow float64, hpperiod1 float64, hpperiod2 float64, tf string, jpy bool) {
	// Caso faça otimização em loop, cuidar pois quando os hpperiod sao iguais, a flag sera 0
	var start = time.Now()
	var newPrice = data.RawTimeFrame(&myData.Open, &myData.Time, tf)
	var decyclerFast = Decycler(&newPrice, cutoffFast, jpy)
	var decyclerSlow = Decycler(&newPrice, cutoffSlow, jpy)
	var oscilator = DecyclerOscilator(&newPrice, hpperiod1, hpperiod2, jpy)
	var buy = make([]bool, len(newPrice))
	var sell = make([]bool, len(newPrice))
	var sizeBuy uint32
	var sizeSell uint32
	for i := range newPrice {
		if decyclerFast[i] > decyclerSlow[i] && oscilator[i] > 0. {
			buy[i] = true
			sizeBuy++
		} else if decyclerFast[i] < decyclerSlow[i] && oscilator[i] < 0 {
			sell[i] = true
			sizeSell++
		}
	}

	var newBuy = data.NormalizeTimeFrameBool(&buy, &myData.Time, tf)
	var newSell = data.NormalizeTimeFrameBool(&sell, &myData.Time, tf)
	(*myData).BuyFlag = newBuy
	(*myData).SellFlag = newSell
	(*myData).SizeBuy = sizeBuy + 1
	(*myData).SizeSell = sizeSell + 1
	var elapsed = time.Since(start)

	table.StartBody(fmt.Sprint("Successfully Created Strategy - Slow & Fast Decycler Cross w/ Oscilator - ", "Buy: ", sizeBuy, " ", "Sell: ", sizeSell), elapsed)
}

func StrategyTwoDecyclerHL(myData *data.LayoutData, cutoffFast float64, cutoffSlow float64, period int, tf string, jpy bool) {
	var start = time.Now()
	var newPrice = data.RawTimeFrame(&myData.Open, &myData.Time, tf)
	var high = data.RawTimeFrame(&myData.High, &myData.Time, tf)
	var low = data.RawTimeFrame(&myData.Low, &myData.Time, tf)
	var highPeriod = mystats.MaxValueOverPeriod(&high, period)
	var lowPeriod = mystats.MinValueOverPeriod(&low, period)
	var decyclerFast = Decycler(&newPrice, cutoffFast, jpy)
	var decyclerSlow = Decycler(&newPrice, cutoffSlow, jpy)
	var buy = make([]bool, len(newPrice))
	var sell = make([]bool, len(newPrice))
	var sizeBuy uint32
	var sizeSell uint32
	for i := range newPrice {
		if i < 1 {
			buy[i] = false
			sell[i] = false
			continue
		}
		if decyclerFast[i] > decyclerSlow[i] && high[i] > highPeriod[i] {
			buy[i] = true
			sizeBuy++
		} else if decyclerFast[i] < decyclerSlow[i] && low[i] < lowPeriod[i] {
			sell[i] = true
			sizeSell++
		}
	}
	var newBuy = data.NormalizeTimeFrameBool(&buy, &myData.Time, tf)
	var newSell = data.NormalizeTimeFrameBool(&sell, &myData.Time, tf)
	(*myData).BuyFlag = newBuy
	(*myData).SellFlag = newSell
	(*myData).SizeBuy = sizeBuy + 1
	(*myData).SizeSell = sizeSell + 1
	var elapsed = time.Since(start)
	table.StartBody(fmt.Sprint("Successfully Created Strategy - Slow & Fast Decycler HL - ", "Buy: ", sizeBuy, " ", "Sell: ", sizeSell), elapsed)
}

func StrategyTwoDecyclerSingleCross(myData *data.LayoutData, cutoffFast float64, cutoffSlow float64, tf string, jpy bool) {
	var start = time.Now()
	var newPrice = data.RawTimeFrame(&myData.Open, &myData.Time, tf)
	var decyclerFast = Decycler(&newPrice, cutoffFast, jpy)
	var decyclerSlow = Decycler(&newPrice, cutoffSlow, jpy)
	var buy = make([]bool, len(newPrice))
	var sell = make([]bool, len(newPrice))
	var sizeBuy uint32
	var sizeSell uint32
	for i := range newPrice {
		if i < 1 {
			buy[i] = false
			sell[i] = false
			continue
		}
		if decyclerFast[i] > decyclerSlow[i] && decyclerFast[i-1] < decyclerSlow[i-1] {
			buy[i] = true
			sizeBuy++
		} else if decyclerFast[i] < decyclerSlow[i] && decyclerFast[i-1] > decyclerSlow[i-1] {
			sell[i] = true
			sizeSell++
		}
	}
	var newBuy = data.NormalizeTimeFrameBool(&buy, &myData.Time, tf)
	var newSell = data.NormalizeTimeFrameBool(&sell, &myData.Time, tf)
	(*myData).BuyFlag = newBuy
	(*myData).SellFlag = newSell
	(*myData).SizeBuy = sizeBuy + 1
	(*myData).SizeSell = sizeSell + 1
	var elapsed = time.Since(start)
	table.StartBody(fmt.Sprint("Successfully Created Strategy - Slow & Fast Decycler Cross - ", "Buy: ", sizeBuy, " ", "Sell: ", sizeSell), elapsed)
}

func StrategyTwoDecyclerCross(myData *data.LayoutData, cutoffFast float64, cutoffSlow float64, tf string, jpy bool) {
	// Heging - Fix TPSL 10  - TPSL = 19 2.5 - slow i=240 fast j=210 Mostraram resultados promissores
	// Hedging - ATR TPSL  shows more interesting
	// Netting didnt show anyhting good
	var start = time.Now()
	var newPrice = data.RawTimeFrame(&myData.Open, &myData.Time, tf)
	var decyclerFast = Decycler(&newPrice, cutoffFast, jpy)
	var decyclerSlow = Decycler(&newPrice, cutoffSlow, jpy)
	var buy = make([]bool, len(newPrice))
	var sell = make([]bool, len(newPrice))
	var sizeBuy uint32
	var sizeSell uint32
	for i := range newPrice {
		if decyclerFast[i] > decyclerSlow[i] {
			buy[i] = true
			sizeBuy++
		} else if decyclerFast[i] < decyclerSlow[i] {
			sell[i] = true
			sizeSell++
		}
	}
	var newBuy = data.NormalizeTimeFrameBool(&buy, &myData.Time, tf)
	var newSell = data.NormalizeTimeFrameBool(&sell, &myData.Time, tf)
	(*myData).BuyFlag = newBuy
	(*myData).SellFlag = newSell
	(*myData).SizeBuy = sizeBuy + 1
	(*myData).SizeSell = sizeSell + 1
	var elapsed = time.Since(start)
	table.StartBody(fmt.Sprint("Successfully Created Strategy - Slow & Fast Decycler Cross - ", "Buy: ", sizeBuy, " ", "Sell: ", sizeSell), elapsed)
}

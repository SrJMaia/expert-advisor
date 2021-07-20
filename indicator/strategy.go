package indicator

import (
	"fmt"
	"time"

	"github.com/SrJMaia/expert-advisor/data"
	"github.com/SrJMaia/expert-advisor/mystats"
	"github.com/SrJMaia/expert-advisor/table"
)

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

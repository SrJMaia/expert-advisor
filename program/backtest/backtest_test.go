package backtest

import (
	"fmt"
	"testing"

	"github.com/SrJMaia/expert-advisor/program/conversion"
	"github.com/SrJMaia/expert-advisor/program/data"
)

func TestHedgingBacktest(t *testing.T) {
	/*
		Because float precision, probably, there is a difference between expected
		and the reality, but the difference is really small
		So I cant predict the real result, so instead using the absolute value
		I'll use the percent and my hypothesis is:
		p = 0.01
		Null Hypothesis: The backtest does not work
		If the difference is less than p: We reject the null hypothesis and we accept
		that the backtest works
		If the difference is higher than p: We assume the null hypothesis, we will need to
		do another backtest from scratch or find the error
	*/
	// EUR HEDGING
	var myData = data.LayoutData{
		Open:       []float64{1.001, 1.002, 1.003, 1.004, 1.005, 1.006},
		High:       []float64{1.003, 1.004, 1.005, 1.006, 1.007, 1.008},
		Low:        []float64{1.000, 1.001, 1.002, 1.003, 1.004, 1.005},
		Close:      []float64{1.002, 1.003, 1.004, 1.005, 1.006, 1.007},
		PriceTf:    []float64{1.001, 1.002, 1.002, 1.002, 1.002, 1.002},
		TpslNonFix: []float64{0.001, 0.001, 0.001, 0.001, 0.001, 0.001},
		BuyFlag:    []bool{true, true, true, false, false, false},
		SellFlag:   []bool{false, true, true, true, false, false},
		SizeBuy:    3,
		SizeSell:   3,
	}
	var capital = 100000.
	var pValue = 0.0001 // Even 0.01c can make a huge difference in higher interests
	var tot, buy, sell = HedgingBacktest(&myData, 3., 1., capital, false, false)

	var previsionTot = conversion.Round(0.003*capital/1.002-(capital/1000*0.07)+capital, 2)
	var difference = tot[1] / previsionTot
	if difference < 1 {
		difference = 1 - difference
	} else {
		difference = difference - 1
	}
	if difference > pValue {
		t.Error("EUR - HEDGING - Tot - Null Hypothesis Accepted! Result:", difference, "P Value:", pValue, "Expected:", previsionTot, "Got:", tot[1])
	}
	var previsionBuy = conversion.Round(0.003*capital/1.002-(capital/1000*0.07)+capital, 2)
	difference = buy[1] / previsionBuy
	if difference < 1 {
		difference = 1 - difference
	} else {
		difference = difference - 1
	}
	if difference > pValue {
		t.Error("EUR - HEDGING - Buy - Null Hypothesis Accepted! Result:", difference, "P Value:", pValue, "Expected:", previsionBuy, "Got:", buy[1])
	}
	var previsionSell = conversion.Round(-0.001*capital/1.002-(capital/1000*0.07)+capital, 2)
	difference = sell[1] / previsionSell
	if difference < 1 {
		difference = 1 - difference
	} else {
		difference = difference - 1
	}
	if difference > pValue {
		t.Error("SEUR - HEDGING - ell - Null Hypothesis Accepted! Result:", difference, "P Value:", pValue, "Expected:", previsionSell, "Got:", sell[1])
	}

	// EUR NETTING
	tot, buy, sell = NettingBacktest(&myData, 3., 1., capital, false, false)

	previsionTot = conversion.Round(0.003*capital/1.002-(capital/1000*0.07)+capital, 2)
	difference = tot[1] / previsionTot
	if difference < 1 {
		difference = 1 - difference
	} else {
		difference = difference - 1
	}
	if difference > pValue {
		t.Error("EUR - NETTING - Tot - Null Hypothesis Accepted! Result:", difference, "P Value:", pValue, "Expected:", previsionTot, "Got:", tot[1])
	}
	previsionBuy = conversion.Round(0.003*capital/1.002-(capital/1000*0.07)+capital, 2)
	difference = buy[1] / previsionBuy
	if difference < 1 {
		difference = 1 - difference
	} else {
		difference = difference - 1
	}
	if difference > pValue {
		t.Error("EUR - NETTING - Buy - Null Hypothesis Accepted! Result:", difference, "P Value:", pValue, "Expected:", previsionBuy, "Got:", buy[1])
	}
	if len(sell) > 1 {
		t.Error("Expected:", 1, "Got:", sell)
	}

	// EUR NoTpSL
	myData = data.LayoutData{
		Open:       []float64{1.001, 1.002, 1.003, 1.004, 1.005, 1.006},
		High:       []float64{1.003, 1.004, 1.005, 1.006, 1.007, 1.008},
		Low:        []float64{1.000, 1.001, 1.002, 1.003, 1.004, 1.005},
		Close:      []float64{1.002, 1.003, 1.004, 1.005, 1.006, 1.007},
		PriceTf:    []float64{1.001, 1.002, 1.002, 1.003, 1.003, 1.004},
		TpslNonFix: []float64{0.001, 0.001, 0.001, 0.001, 0.001, 0.001},
		BuyFlag:    []bool{false, true, false, false, false, true},
		SellFlag:   []bool{false, false, false, true, false, false},
		SizeBuy:    3,
		SizeSell:   3,
	}

	_, buy, sell = NoTpslBacktest(&myData, capital, false, false)

	previsionBuy = conversion.Round(0.001*capital/1.003-(capital/1000*0.07)+capital, 2)
	difference = buy[1] / previsionBuy
	if difference < 1 {
		difference = 1 - difference
	} else {
		difference = difference - 1
	}
	if difference > pValue {
		t.Error("EUR - NOTPSL - Buy - Null Hypothesis Accepted! Result:", difference, "P Value:", pValue, "Expected:", previsionBuy, "Got:", buy[1])
	}
	previsionSell = conversion.Round(-0.001*capital/1.003-(capital/1000*0.07)+capital, 2)
	difference = sell[1] / previsionSell
	if difference < 1 {
		difference = 1 - difference
	} else {
		difference = difference - 1
	}
	if difference > pValue {
		t.Error("EUR - NOTPSL - Sell - Null Hypothesis Accepted! Result:", difference, "P Value:", pValue, "Expected:", previsionSell, "Got:", sell[1])
	}

	// JPY HEDGING
	myData = data.LayoutData{
		Open:       []float64{123.1, 123.2, 123.3, 123.4, 123.5, 123.6},
		High:       []float64{123.3, 123.4, 123.5, 123.6, 123.7, 123.8},
		Low:        []float64{123.0, 123.1, 123.2, 123.3, 123.4, 123.5},
		Close:      []float64{123.2, 123.3, 123.4, 123.5, 123.6, 123.7},
		PriceTf:    []float64{123.1, 123.2, 123.2, 123.2, 123.2, 123.2},
		TpslNonFix: []float64{0.1, 0.1, 0.1, 0.1, 0.1, 0.1},
		BuyFlag:    []bool{true, true, true, false, false, false},
		SellFlag:   []bool{false, true, true, true, false, false},
		SizeBuy:    3,
		SizeSell:   3,
	}
	tot, buy, sell = HedgingBacktest(&myData, 3., 1., capital, true, false)

	previsionTot = conversion.Round(0.3*capital/123.2-(capital/1000*0.07)+capital, 2)
	difference = tot[1] / previsionTot
	if difference < 1 {
		difference = 1 - difference
	} else {
		difference = difference - 1
	}
	if difference > pValue {
		t.Error("Tot - JPY - Null Hypothesis Accepted! Result:", difference, "P Value:", pValue, "Expected:", previsionTot, "Got:", tot[1])
	}
	previsionBuy = conversion.Round(0.3*capital/123.2-(capital/1000*0.07)+capital, 2)
	difference = buy[1] / previsionBuy
	if difference < 1 {
		difference = 1 - difference
	} else {
		difference = difference - 1
	}
	if difference > pValue {
		t.Error("Buy - JPY - Null Hypothesis Accepted! Result:", difference, "P Value:", pValue, "Expected:", previsionBuy, "Got:", buy[1])
	}
	previsionSell = conversion.Round(-0.1*capital/123.2-(capital/1000*0.07)+capital, 2)
	difference = sell[1] / previsionSell
	if difference < 1 {
		difference = 1 - difference
	} else {
		difference = difference - 1
	}
	if difference > pValue {
		t.Error("Sell - JPY - Null Hypothesis Accepted! Result:", difference, "P Value:", pValue, "Expected:", previsionSell, "Got:", sell[1])
	}

	// JPY NETTING
	tot, buy, sell = NettingBacktest(&myData, 3., 1., capital, true, false)
	previsionTot = conversion.Round(0.3*capital/123.2-(capital/1000*0.07)+capital, 2)
	difference = tot[1] / previsionTot
	if difference < 1 {
		difference = 1 - difference
	} else {
		difference = difference - 1
	}
	if difference > pValue {
		t.Error("Tot - JPY - Null Hypothesis Accepted! Result:", difference, "P Value:", pValue, "Expected:", previsionTot, "Got:", tot[1])
	}
	previsionBuy = conversion.Round(0.3*capital/123.2-(capital/1000*0.07)+capital, 2)
	difference = buy[1] / previsionBuy
	if difference < 1 {
		difference = 1 - difference
	} else {
		difference = difference - 1
	}
	if difference > pValue {
		t.Error("Buy - JPY - Null Hypothesis Accepted! Result:", difference, "P Value:", pValue, "Expected:", previsionBuy, "Got:", buy[1])
	}
	if len(sell) > 1 {
		t.Error("Expected:", 1, "Got:", sell)
	}

	// JPY NoTpSL
	myData = data.LayoutData{
		Open:       []float64{123.1, 123.2, 123.3, 123.4, 123.5, 123.6},
		High:       []float64{123.3, 123.4, 123.5, 123.6, 123.7, 123.8},
		Low:        []float64{123.0, 123.1, 123.2, 123.3, 123.4, 123.5},
		Close:      []float64{123.2, 123.3, 123.4, 123.5, 123.6, 123.7},
		PriceTf:    []float64{123.1, 123.2, 123.2, 123.3, 123.3, 123.4},
		TpslNonFix: []float64{0.1, 0.1, 0.1, 0.1, 0.1, 0.1},
		BuyFlag:    []bool{false, true, false, false, false, true},
		SellFlag:   []bool{false, false, false, true, false, false},
		SizeBuy:    3,
		SizeSell:   3,
	}

	tot, buy, sell = NoTpslBacktest(&myData, capital, false, false)
	fmt.Println(tot)
	fmt.Println(buy)
	fmt.Println(sell)

	previsionBuy = conversion.Round(0.1*capital/123.3-(capital/1000*0.07)+capital, 2)
	difference = buy[1] / previsionBuy
	if difference < 1 {
		difference = 1 - difference
	} else {
		difference = difference - 1
	}
	if difference > pValue {
		t.Error("EUR - NOTPSL - Buy - Null Hypothesis Accepted! Result:", difference, "P Value:", pValue, "Expected:", previsionBuy, "Got:", buy[1])
	}
	previsionSell = conversion.Round(-0.1*capital/123.3-(capital/1000*0.07)+capital, 2)
	difference = sell[1] / previsionSell
	if difference < 1 {
		difference = 1 - difference
	} else {
		difference = difference - 1
	}
	if difference > pValue {
		t.Error("EUR - NOTPSL - Sell - Null Hypothesis Accepted! Result:", difference, "P Value:", pValue, "Expected:", previsionSell, "Got:", sell[1])
	}
}

func TestCheckTpslExit(t *testing.T) {
	// Buy | tp = true e sl = false
	var open = 1.123
	var high = 1.125
	var low = 1.122
	var close = 1.124
	var tpTrue = 1.124
	var tpFalse = 1.126
	var slTrue = 1.123
	var slFalse = 1.11
	var tp = checkTpslExit(open, high, low, close, tpTrue, true)
	if !tp {
		t.Error("EUR - BUY - TP - Expected:", true, "Got:", tp)
	}
	tp = checkTpslExit(open, high, low, close, tpFalse, true)
	if tp {
		t.Error("EUR - BUY - TP - Expected:", false, "Got:", tp)
	}
	var sl = checkTpslExit(open, high, low, close, slTrue, false)
	if !sl {
		t.Error("EUR - BUY - SL - Expected:", true, "Got:", sl)
	}
	sl = checkTpslExit(open, high, low, close, slFalse, false)
	if sl {
		t.Error("EUR - BUY - SL - Expected:", false, "Got:", sl)
	}
	// Sell | tp = false e sl = true
	tpTrue = 1.1225
	tpFalse = 1.121
	slTrue = 1.1245
	slFalse = 1.13
	tp = checkTpslExit(open, high, low, close, tpTrue, false)
	if !tp {
		t.Error("EUR - SELL - TP - Expected:", true, "Got:", tp)
	}
	tp = checkTpslExit(open, high, low, close, tpFalse, false)
	if tp {
		t.Error("EUR - SELL - TP - Expected:", false, "Got:", tp)
	}
	sl = checkTpslExit(open, high, low, close, slTrue, true)
	if !sl {
		t.Error("EUR - SELL - SL - Expected:", true, "Got:", sl)
	}
	sl = checkTpslExit(open, high, low, close, slFalse, true)
	if sl {
		t.Error("EUR - SELL - SL - Expected:", false, "Got:", sl)
	}

}

func TestTpslCalculation(t *testing.T) {
	var backtestMain = backtestHeart(1000, 10, 10)
	var multiplyValue = 0.001
	var price = 1.123
	var multiplyTp = 2.5
	var multiplySl = 1.25
	var jpy = false
	var buy = true
	tpslCalculation(&backtestMain, price, multiplyTp, multiplySl, multiplyValue, jpy, buy)
	if backtestMain.tpBuy != 1.1255 {
		t.Error("EUR - BUY - Expected:", 1.1255, "Got:", backtestMain.tpBuy)
	} else if backtestMain.slBuy != 1.12175 {
		t.Error("EUR - BUY - Expected:", 1.12175, "Got:", backtestMain.slBuy)
	}
	price = 1.123
	multiplyTp = 2.5
	multiplySl = 1.25
	jpy = false
	buy = false
	tpslCalculation(&backtestMain, price, multiplyTp, multiplySl, multiplyValue, jpy, buy)
	if backtestMain.tpSell != 1.1205 {
		t.Error("EUR - SELL - Expected:", 1.1205, "Got:", backtestMain.tpSell)
	} else if backtestMain.slSell != 1.12425 {
		t.Error("EUR - SELL - Expected:", 1.12425, "Got:", backtestMain.slSell)
	}

	multiplyValue = 0.1
	price = 123.456
	multiplyTp = 2.5
	multiplySl = 1.25
	jpy = true
	buy = true
	tpslCalculation(&backtestMain, price, multiplyTp, multiplySl, multiplyValue, jpy, buy)
	if backtestMain.tpBuy != 123.706 {
		t.Error("JPY - BUY - Expected:", 123.706, "Got:", backtestMain.tpBuy)
	} else if backtestMain.slBuy != 123.331 {
		t.Error("JPY - BUY - Expected:", 123.331, "Got:", backtestMain.slBuy)
	}
	price = 123.456
	multiplyTp = 2.5
	multiplySl = 1.25
	jpy = true
	buy = false
	tpslCalculation(&backtestMain, price, multiplyTp, multiplySl, multiplyValue, jpy, buy)
	if backtestMain.tpSell != 123.206 {
		t.Error("JPY - SELL - Expected:", 123.206, "Got:", backtestMain.tpSell)
	} else if backtestMain.slSell != 123.581 {
		t.Error("JPY - SELL - Expected:", 123.581, "Got:", backtestMain.slSell)
	}
}

func TestFinanceCalculation(t *testing.T) {
	var leverage = false

	var backtestMain = backtestHeart(1000., 10, 10)
	var jpy = false

	backtestMain.tpBuy = 1.124
	backtestMain.buyPrice = 1.123
	financeCalculation(&backtestMain, leverage, jpy, "buy-tp")
	if backtestMain.buyResult != 0.81968 {
		t.Error("EUR BUY TP - Expected:", 0.81968, "Got:", backtestMain.buyResult)
	}
	if backtestMain.capital != 1000.82 {
		t.Error("EUR BUY TP - Expected:", 1000.82, "Got:", backtestMain.capital)
	}

	backtestMain = backtestHeart(1000., 10, 10)
	backtestMain.slBuy = 1.122
	backtestMain.buyPrice = 1.123
	financeCalculation(&backtestMain, leverage, jpy, "buy-sl")
	if backtestMain.buyResult != -0.96127 {
		t.Error("EUR BUY SL - Expected:", -0.96127, "Got:", backtestMain.buyResult)
	}
	if backtestMain.capital != 999.04 {
		t.Error("EUR BUY SL - Expected:", 999.04, "Got:", backtestMain.capital)
	}

	backtestMain = backtestHeart(1000., 10, 10)
	backtestMain.tpSell = 1.122
	backtestMain.sellPrice = 1.123
	financeCalculation(&backtestMain, leverage, jpy, "sell-tp")
	if backtestMain.sellResult != 0.82127 {
		t.Error("EUR SELL TP - Expected:", 0.82127, "Got:", backtestMain.sellResult)
	}
	if backtestMain.capital != 1000.82 {
		t.Error("EUR SELL TP - Expected:", 1000.82, "Got:", backtestMain.capital)
	}

	backtestMain = backtestHeart(1000., 10, 10)
	backtestMain.slSell = 1.124
	backtestMain.sellPrice = 1.123
	financeCalculation(&backtestMain, leverage, jpy, "sell-sl")
	if backtestMain.sellResult != -0.95968 {
		t.Error("EUR SELL SL - Expected:", -0.95968, "Got:", backtestMain.sellResult)
	}
	if backtestMain.capital != 999.04 {
		t.Error("EUR SELL SL - Expected:", 999.04, "Got:", backtestMain.capital)
	}

	// JPY
	backtestMain = backtestHeart(1000., 10, 10)
	backtestMain.tpBuy = 127.204
	backtestMain.buyPrice = 127.124
	jpy = true
	financeCalculation(&backtestMain, leverage, jpy, "buy-tp")
	if backtestMain.buyResult != 0.559 {
		t.Error("JPY BUY TP - Expected:", 0.559, "Got:", backtestMain.buyResult)
	}
	if backtestMain.capital != 1000.56 {
		t.Error("JPY BUY TP - Expected:", 1000.56, "Got:", backtestMain.capital)
	}

	backtestMain = backtestHeart(1000., 10, 10)
	backtestMain.slBuy = 127.05
	backtestMain.buyPrice = 127.1
	jpy = true
	financeCalculation(&backtestMain, leverage, jpy, "buy-sl")
	if backtestMain.buyResult != -0.464 {
		t.Error("JPY BUY SL - Expected:", -0.464, "Got:", backtestMain.buyResult)
	}
	if backtestMain.capital != 999.54 {
		t.Error("JPY BUY SL - Expected:", 999.54, "Got:", backtestMain.capital)
	}

	backtestMain = backtestHeart(1000., 10, 10)
	backtestMain.tpSell = 127.05
	backtestMain.sellPrice = 127.1
	jpy = true
	financeCalculation(&backtestMain, leverage, jpy, "sell-tp")
	if backtestMain.sellResult != 0.324 {
		t.Error("JPY SELL TP - Expected:", 0.324, "Got:", backtestMain.sellResult)
	}
	if backtestMain.capital != 1000.32 {
		t.Error("JPY SELL TP - Expected:", 1000.32, "Got:", backtestMain.capital)
	}

	backtestMain = backtestHeart(1000., 10, 10)
	backtestMain.slSell = 127.1
	backtestMain.sellPrice = 127.05
	jpy = true
	financeCalculation(&backtestMain, leverage, jpy, "sell-sl")
	if backtestMain.sellResult != -0.463 {
		t.Error("JPY SELL SL - Expected:", -0.463, "Got:", backtestMain.sellResult)
	}
	if backtestMain.capital != 999.54 {
		t.Error("JPY SELL SL - Expected:", 999.54, "Got:", backtestMain.capital)
	}

}

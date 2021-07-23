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
		Open:     []float64{1.001, 1.002, 1.003, 1.004, 1.005, 1.006},
		High:     []float64{1.003, 1.004, 1.005, 1.006, 1.007, 1.008},
		Low:      []float64{1.000, 1.001, 1.002, 1.003, 1.004, 1.005},
		Close:    []float64{1.002, 1.003, 1.004, 1.005, 1.006, 1.007},
		PriceTf:  []float64{1.001, 1.002, 1.002, 1.002, 1.002, 1.002},
		Tpsl:     []float64{0.001, 0.001, 0.001, 0.001, 0.001, 0.001},
		BuyFlag:  []bool{true, true, true, false, false, false},
		SellFlag: []bool{false, true, true, true, false, false},
		SizeBuy:  3,
		SizeSell: 3,
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
		Open:     []float64{1.001, 1.002, 1.003, 1.004, 1.005, 1.006},
		High:     []float64{1.003, 1.004, 1.005, 1.006, 1.007, 1.008},
		Low:      []float64{1.000, 1.001, 1.002, 1.003, 1.004, 1.005},
		Close:    []float64{1.002, 1.003, 1.004, 1.005, 1.006, 1.007},
		PriceTf:  []float64{1.001, 1.002, 1.002, 1.003, 1.003, 1.004},
		Tpsl:     []float64{0.001, 0.001, 0.001, 0.001, 0.001, 0.001},
		BuyFlag:  []bool{false, true, false, false, false, true},
		SellFlag: []bool{false, false, false, true, false, false},
		SizeBuy:  3,
		SizeSell: 3,
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
		Open:     []float64{123.1, 123.2, 123.3, 123.4, 123.5, 123.6},
		High:     []float64{123.3, 123.4, 123.5, 123.6, 123.7, 123.8},
		Low:      []float64{123.0, 123.1, 123.2, 123.3, 123.4, 123.5},
		Close:    []float64{123.2, 123.3, 123.4, 123.5, 123.6, 123.7},
		PriceTf:  []float64{123.1, 123.2, 123.2, 123.2, 123.2, 123.2},
		Tpsl:     []float64{0.1, 0.1, 0.1, 0.1, 0.1, 0.1},
		BuyFlag:  []bool{true, true, true, false, false, false},
		SellFlag: []bool{false, true, true, true, false, false},
		SizeBuy:  3,
		SizeSell: 3,
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
		Open:     []float64{123.1, 123.2, 123.3, 123.4, 123.5, 123.6},
		High:     []float64{123.3, 123.4, 123.5, 123.6, 123.7, 123.8},
		Low:      []float64{123.0, 123.1, 123.2, 123.3, 123.4, 123.5},
		Close:    []float64{123.2, 123.3, 123.4, 123.5, 123.6, 123.7},
		PriceTf:  []float64{123.1, 123.2, 123.2, 123.3, 123.3, 123.4},
		Tpsl:     []float64{0.1, 0.1, 0.1, 0.1, 0.1, 0.1},
		BuyFlag:  []bool{false, true, false, false, false, true},
		SellFlag: []bool{false, false, false, true, false, false},
		SizeBuy:  3,
		SizeSell: 3,
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

func TestTpslCalculationNonFix(t *testing.T) {
	var multiplyValue = 0.001
	var price = 1.123
	var multiplyTp = 2.5
	var multiplySl = 1.25
	var jpy = false
	var buy = true
	var tp, sl = tpslCalculationNonFix(price, multiplyTp, multiplySl, multiplyValue, jpy, buy)
	if tp != 1.1255 {
		t.Error("EUR - BUY - Expected:", 1.1255, "Got:", tp)
	} else if sl != 1.12175 {
		t.Error("EUR - BUY - Expected:", 1.12175, "Got:", sl)
	}
	price = 1.123
	multiplyTp = 2.5
	multiplySl = 1.25
	jpy = false
	buy = false
	tp, sl = tpslCalculationNonFix(price, multiplyTp, multiplySl, multiplyValue, jpy, buy)
	if tp != 1.1205 {
		t.Error("EUR - SELL - Expected:", 1.1205, "Got:", tp)
	} else if sl != 1.12425 {
		t.Error("EUR - SELL - Expected:", 1.12425, "Got:", sl)
	}

	multiplyValue = 0.1
	price = 123.456
	multiplyTp = 2.5
	multiplySl = 1.25
	jpy = true
	buy = true
	tp, sl = tpslCalculationNonFix(price, multiplyTp, multiplySl, multiplyValue, jpy, buy)
	if tp != 123.706 {
		t.Error("JPY - BUY - Expected:", 123.706, "Got:", tp)
	} else if sl != 123.331 {
		t.Error("JPY - BUY - Expected:", 123.331, "Got:", sl)
	}
	price = 123.456
	multiplyTp = 2.5
	multiplySl = 1.25
	jpy = true
	buy = false
	tp, sl = tpslCalculationNonFix(price, multiplyTp, multiplySl, multiplyValue, jpy, buy)
	if tp != 123.206 {
		t.Error("JPY - SELL - Expected:", 123.206, "Got:", tp)
	} else if sl != 123.581 {
		t.Error("JPY - SELL - Expected:", 123.581, "Got:", sl)
	}
}

func TestTpslCalculationFix(t *testing.T) {
	var price = 1.123
	var multiplyTp = 0.001 * 2.5
	var multiplySl = 0.001 * 1.25
	var jpy = false
	var buy = true
	var tp, sl = tpslCalculationFix(price, multiplyTp, multiplySl, jpy, buy)
	if tp != 1.1255 {
		t.Error("EUR - BUY - Expected:", 1.1255, "Got:", tp)
	} else if sl != 1.12175 {
		t.Error("EUR - BUY - Expected:", 1.12175, "Got:", sl)
	}
	price = 1.123
	multiplyTp = 0.001 * 2.5
	multiplySl = 0.001 * 1.25
	jpy = false
	buy = false
	tp, sl = tpslCalculationFix(price, multiplyTp, multiplySl, jpy, buy)
	if tp != 1.1205 {
		t.Error("EUR - SELL - Expected:", 1.1205, "Got:", tp)
	} else if sl != 1.12425 {
		t.Error("EUR - SELL - Expected:", 1.12425, "Got:", sl)
	}

	price = 123.456
	multiplyTp = 0.1 * 2.5
	multiplySl = 0.1 * 1.25
	jpy = true
	buy = true
	tp, sl = tpslCalculationFix(price, multiplyTp, multiplySl, jpy, buy)
	if tp != 123.706 {
		t.Error("JPY - BUY - Expected:", 123.706, "Got:", tp)
	} else if sl != 123.331 {
		t.Error("JPY - BUY - Expected:", 123.331, "Got:", sl)
	}
	price = 123.456
	multiplyTp = 0.1 * 2.5
	multiplySl = 0.1 * 1.25
	jpy = true
	buy = false
	tp, sl = tpslCalculationFix(price, multiplyTp, multiplySl, jpy, buy)
	if tp != 123.206 {
		t.Error("JPY - SELL - Expected:", 123.206, "Got:", tp)
	} else if sl != 123.581 {
		t.Error("JPY - SELL - Expected:", 123.581, "Got:", sl)
	}

}

func TestFinanceCalculation(t *testing.T) {
	var balance = 1000.
	var initialBalance = 1000.
	var leverage = false
	var initialPrice = 1.124
	var finalPrice = 1.123
	var eurPrice = 1.124
	var total, result = financeCalculation(balance, initialPrice, finalPrice, eurPrice, initialBalance, leverage)
	if result != 0.81968 {
		t.Error("EUR - Expected:", 0.81968, "Got:", result)
	} else if total != 1000.82 {
		t.Error("EUR - Expected:", 1000.82, "Got:", total)
	}
	initialPrice = 127.204
	finalPrice = 127.124
	eurPrice = 113.550
	total, result = financeCalculation(balance, initialPrice, finalPrice, eurPrice, initialBalance, leverage)
	if result != 0.63454 {
		t.Error("JPY - Expected:", 0.63454, "Got:", result)
	} else if total != 1000.63 {
		t.Error("JPY - Expected:", 1000.63, "Got:", total)
	}

}

package backtest

import (
	"testing"

	"github.com/SrJMaia/expert-advisor/program/conversion"
	"github.com/SrJMaia/expert-advisor/program/data"
)

func TestFimatheBacktest(t *testing.T) {
	// EUR
	var myData = data.LayoutData{
		Open:              []float64{1.001, 1.002, 1.003, 1.004, 1.005, 1.006, 1.005, 1.003, 1.001, 0.999, 0.997},
		High:              []float64{1.003, 1.004, 1.005, 1.006, 1.007, 1.008, 1.007, 1.005, 1.003, 1.001, 0.999},
		Low:               []float64{1.000, 1.001, 1.002, 1.003, 1.004, 1.005, 1.004, 1.002, 1.0, 0.998, 0.996},
		Close:             []float64{1.002, 1.003, 1.004, 1.005, 1.006, 1.007, 1.006, 1.004, 1.002, 1.0, 0.998},
		PriceTf:           []float64{1.001, 1.002, 1.002, 1.004, 1.004, 1.003, 1.003, 1.003, 1.001, 0.999, 0.999},
		BuyFlag:           []bool{true, true, true, false, false, false, false},
		SellFlag:          []bool{false, true, true, true, false, false, false},
		SizeBuy:           4,
		SizeSell:          4,
		IsFixTpsl:         true,
		TpslFix:           0.001,
		MultiplyTp:        3.,
		MultiplySl:        2.,
		ChannelInPips:     0.002,
		StartPointInPips:  0.0,
		MultiplyTpChannel: 2.,
	}
	var balance = 1000.
	var jpy = false
	var leverage = false
	var tot, buy, sell = BacktestMain(&myData, balance, "fimathe", jpy, leverage)

	var backtestVariables = backtestHeart(balance, myData.SizeBuy, myData.SizeSell, "", jpy, leverage)
	backtestVariables.buyPrice = 1.004
	var isBuy = true
	backtestVariables.upperChannel = 1.001 + (myData.ChannelInPips / 2)
	takeProfitFimathe(backtestVariables, myData.MultiplyTpChannel, myData.ChannelInPips, isBuy)
	financeCalculation(backtestVariables, "buy-tp")
	backtestVariables.sellPrice = 0.999
	isBuy = false
	backtestVariables.lowerChannel = 1.001 - (myData.ChannelInPips / 2)
	takeProfitFimathe(backtestVariables, myData.MultiplyTpChannel, myData.ChannelInPips, isBuy)
	financeCalculation(backtestVariables, "sell-tp")
	var totalTrades = conversion.RemoveExcessZeros(backtestVariables.totalTrades)
	var buyTrades = conversion.RemoveExcessZeros(backtestVariables.buyTrades)
	var sellTrades = conversion.RemoveExcessZeros(backtestVariables.sellTrades)
	for i := range tot {
		if tot[i] != totalTrades[i] {
			t.Error("EUR Total Trades - Expected:", totalTrades, "Got:", tot)
		}
	}
	for i := range buy {
		if buy[i] != buyTrades[i] {
			t.Error("EUR Buy Trades - Expected:", buyTrades, "Got:", buy)
		}
	}
	for i := range sell {
		if sell[i] != sellTrades[i] {
			t.Error("EUR Sell Trades - Expected:", sellTrades, "Got:", sell)
		}
	}

	// JPY
	myData = data.LayoutData{
		Open:              []float64{100.1, 100.2, 100.3, 100.4, 100.5, 100.6, 100.5, 100.3, 100.1, 99.9, 99.7},
		High:              []float64{100.3, 100.4, 100.5, 100.6, 100.7, 100.8, 100.7, 100.5, 100.3, 100.1, 99.9},
		Low:               []float64{100.0, 100.1, 100.2, 100.3, 100.4, 100.5, 100.4, 100.2, 100.0, 99.8, 99.6},
		Close:             []float64{100.2, 100.3, 100.4, 100.5, 100.6, 100.7, 100.6, 100.4, 100.2, 100.0, 99.8},
		PriceTf:           []float64{100.1, 100.2, 100.2, 100.4, 100.4, 100.3, 100.3, 100.3, 100.1, 99.9, 99.9},
		BuyFlag:           []bool{true, true, true, false, false, false, false},
		SellFlag:          []bool{false, true, true, true, false, false, false},
		SizeBuy:           4,
		SizeSell:          4,
		IsFixTpsl:         true,
		TpslFix:           0.001,
		MultiplyTp:        3.,
		MultiplySl:        2.,
		ChannelInPips:     0.2,
		StartPointInPips:  0.0,
		MultiplyTpChannel: 2.,
	}
	balance = 1000.
	jpy = true
	leverage = false
	tot, buy, sell = BacktestMain(&myData, balance, "fimathe", jpy, leverage)

	backtestVariables = backtestHeart(balance, myData.SizeBuy, myData.SizeSell, "", jpy, leverage)
	backtestVariables.buyPrice = 100.4
	isBuy = true
	backtestVariables.upperChannel = 100.1 + (myData.ChannelInPips / 2)
	takeProfitFimathe(backtestVariables, myData.MultiplyTpChannel, myData.ChannelInPips, isBuy)
	financeCalculation(backtestVariables, "buy-tp")
	backtestVariables.sellPrice = 99.9
	isBuy = false
	backtestVariables.lowerChannel = 100.1 - (myData.ChannelInPips / 2)
	takeProfitFimathe(backtestVariables, myData.MultiplyTpChannel, myData.ChannelInPips, isBuy)
	financeCalculation(backtestVariables, "sell-tp")

	totalTrades = conversion.RemoveExcessZeros(backtestVariables.totalTrades)
	buyTrades = conversion.RemoveExcessZeros(backtestVariables.buyTrades)
	sellTrades = conversion.RemoveExcessZeros(backtestVariables.sellTrades)
	for i := range tot {
		if tot[i] != totalTrades[i] {
			t.Error("EUR Total Trades - Expected:", totalTrades, "Got:", tot)
		}
	}
	for i := range buy {
		if buy[i] != buyTrades[i] {
			t.Error("EUR Buy Trades - Expected:", buyTrades, "Got:", buy)
		}
	}
	for i := range sell {
		if sell[i] != sellTrades[i] {
			t.Error("EUR Sell Trades - Expected:", sellTrades, "Got:", sell)
		}
	}
}

func TestFimatheWithStrategyBacktest(t *testing.T) {
	// EUR
	var myData = data.LayoutData{
		Open:              []float64{1.001, 1.002, 1.003, 1.004, 1.005, 1.006, 1.005, 1.003, 1.001, 0.999, 0.997},
		High:              []float64{1.003, 1.004, 1.005, 1.006, 1.007, 1.008, 1.007, 1.005, 1.003, 1.001, 0.999},
		Low:               []float64{1.000, 1.001, 1.002, 1.003, 1.004, 1.005, 1.004, 1.002, 1.0, 0.998, 0.996},
		Close:             []float64{1.002, 1.003, 1.004, 1.005, 1.006, 1.007, 1.006, 1.004, 1.002, 1.0, 0.998},
		PriceTf:           []float64{1.001, 1.002, 1.002, 1.004, 1.004, 1.003, 1.003, 1.003, 1.001, 0.999, 0.999},
		BuyFlag:           []bool{true, true, true, false, false, false, false, false, false, false, false},
		SellFlag:          []bool{false, true, true, true, false, false, false, false, false, true, false},
		SizeBuy:           4,
		SizeSell:          4,
		IsFixTpsl:         true,
		TpslFix:           0.001,
		MultiplyTp:        3.,
		MultiplySl:        2.,
		ChannelInPips:     0.002,
		StartPointInPips:  0.0,
		MultiplyTpChannel: 2.,
	}
	var balance = 1000.
	var jpy = false
	var leverage = false
	var tot, buy, sell = BacktestMain(&myData, balance, "fimathe-with-strategy", jpy, leverage)

	var backtestVariables = backtestHeart(balance, myData.SizeBuy, myData.SizeSell, "", jpy, leverage)
	backtestVariables.sellPrice = 0.999
	var isBuy = false
	backtestVariables.lowerChannel = 1.001 - (myData.ChannelInPips / 2)
	takeProfitFimathe(backtestVariables, myData.MultiplyTpChannel, myData.ChannelInPips, isBuy)
	financeCalculation(backtestVariables, "sell-tp")
	var totalTrades = conversion.RemoveExcessZeros(backtestVariables.totalTrades)
	var buyTrades = conversion.RemoveExcessZeros(backtestVariables.buyTrades)
	var sellTrades = conversion.RemoveExcessZeros(backtestVariables.sellTrades)
	if len(tot) != len(totalTrades) || len(buy) != len(buyTrades) || len(sell) != len(sellTrades) {
		t.Error("EUR Total Trades - Expected:", totalTrades, "Got:", tot)
		t.Error("EUR Buy Trades - Expected:", buyTrades, "Got:", buy)
		t.Error("EUR Sell Trades - Expected:", sellTrades, "Got:", sell)
	}
	for i := range tot {
		if tot[i] != totalTrades[i] {
			t.Error("EUR Total Trades - Expected:", totalTrades, "Got:", tot)
		}
	}
	for i := range buy {
		if buy[i] != buyTrades[i] {
			t.Error("EUR Buy Trades - Expected:", buyTrades, "Got:", buy)
		}
	}
	for i := range sell {
		if sell[i] != sellTrades[i] {
			t.Error("EUR Sell Trades - Expected:", sellTrades, "Got:", sell)
		}
	}

	// JPY
	myData = data.LayoutData{
		Open:              []float64{100.1, 100.2, 100.3, 100.4, 100.5, 100.6, 100.5, 100.3, 100.1, 99.9, 99.7},
		High:              []float64{100.3, 100.4, 100.5, 100.6, 100.7, 100.8, 100.7, 100.5, 100.3, 100.1, 99.9},
		Low:               []float64{100.0, 100.1, 100.2, 100.3, 100.4, 100.5, 100.4, 100.2, 100.0, 99.8, 99.6},
		Close:             []float64{100.2, 100.3, 100.4, 100.5, 100.6, 100.7, 100.6, 100.4, 100.2, 100.0, 99.8},
		PriceTf:           []float64{100.1, 100.2, 100.2, 100.4, 100.4, 100.3, 100.3, 100.3, 100.1, 99.9, 99.9},
		BuyFlag:           []bool{true, true, true, true, false, false, false, false, false, false, false, false},
		SellFlag:          []bool{false, true, true, true, false, false, false, false, false, false, false, false},
		SizeBuy:           4,
		SizeSell:          4,
		IsFixTpsl:         true,
		TpslFix:           0.001,
		MultiplyTp:        3.,
		MultiplySl:        2.,
		ChannelInPips:     0.2,
		StartPointInPips:  0.0,
		MultiplyTpChannel: 2.,
	}
	balance = 1000.
	jpy = true
	leverage = false
	tot, buy, sell = BacktestMain(&myData, balance, "fimathe-with-strategy", jpy, leverage)

	backtestVariables = backtestHeart(balance, myData.SizeBuy, myData.SizeSell, "", jpy, leverage)
	backtestVariables.buyPrice = 100.4
	isBuy = true
	backtestVariables.upperChannel = 100.1 + (myData.ChannelInPips / 2)
	takeProfitFimathe(backtestVariables, myData.MultiplyTpChannel, myData.ChannelInPips, isBuy)
	financeCalculation(backtestVariables, "buy-tp")

	totalTrades = conversion.RemoveExcessZeros(backtestVariables.totalTrades)
	buyTrades = conversion.RemoveExcessZeros(backtestVariables.buyTrades)
	sellTrades = conversion.RemoveExcessZeros(backtestVariables.sellTrades)
	if len(tot) != len(totalTrades) || len(buy) != len(buyTrades) || len(sell) != len(sellTrades) {
		t.Error("EUR Total Trades - Expected:", totalTrades, "Got:", tot)
		t.Error("EUR Buy Trades - Expected:", buyTrades, "Got:", buy)
		t.Error("EUR Sell Trades - Expected:", sellTrades, "Got:", sell)
	}
	for i := range tot {
		if tot[i] != totalTrades[i] {
			t.Error("EUR Total Trades - Expected:", totalTrades, "Got:", tot)
		}
	}
	for i := range buy {
		if buy[i] != buyTrades[i] {
			t.Error("EUR Buy Trades - Expected:", buyTrades, "Got:", buy)
		}
	}
	for i := range sell {
		if sell[i] != sellTrades[i] {
			t.Error("EUR Sell Trades - Expected:", sellTrades, "Got:", sell)
		}
	}
}

func TestNoTpslBacktest(t *testing.T) {
	// EUR
	var myData = data.LayoutData{
		Open:       []float64{1.001, 1.002, 1.003, 1.004, 1.005, 1.006, 1.007, 1.008, 1.009},
		High:       []float64{1.003, 1.004, 1.005, 1.006, 1.007, 1.008, 1.009, 1.01, 1.011},
		Low:        []float64{1.000, 1.001, 1.002, 1.003, 1.004, 1.005, 1.006, 1.007, 1.008},
		Close:      []float64{1.002, 1.003, 1.004, 1.005, 1.006, 1.007, 1.008, 1.009, 1.01},
		PriceTf:    []float64{1.001, 1.002, 1.002, 1.003, 1.003, 1.003, 1.004, 1.004, 1.004},
		BuyFlag:    []bool{true, true, true, false, false, false, true, true, true},
		SellFlag:   []bool{false, true, true, true, false, false, false, false, false},
		SizeBuy:    3,
		SizeSell:   3,
		IsFixTpsl:  true,
		TpslFix:    0.001,
		MultiplyTp: 3.,
		MultiplySl: 2.,
	}
	var balance = 1000.
	var jpy = false
	var leverage = false
	var tot, buy, sell = BacktestMain(&myData, balance, "no-tpsl", jpy, leverage)

	var backtestVariables = backtestHeart(balance, myData.SizeBuy, myData.SizeSell, "", jpy, leverage)
	backtestVariables.buyPrice = 1.002
	backtestVariables.tpBuy = 1.003
	financeCalculation(backtestVariables, "buy-tp")
	backtestVariables.sellPrice = 1.003
	backtestVariables.tpSell = 1.004
	financeCalculation(backtestVariables, "sell-tp")
	var totalTrades = conversion.RemoveExcessZeros(backtestVariables.totalTrades)
	var buyTrades = conversion.RemoveExcessZeros(backtestVariables.buyTrades)
	var sellTrades = conversion.RemoveExcessZeros(backtestVariables.sellTrades)
	for i := range tot {
		if tot[i] != totalTrades[i] {
			t.Error("EUR Total Trades - Expected:", totalTrades, "Got:", tot)
		}
	}
	for i := range buy {
		if buy[i] != buyTrades[i] {
			t.Error("EUR Buy Trades - Expected:", buyTrades, "Got:", buy)
		}
	}
	for i := range sell {
		if sell[i] != sellTrades[i] {
			t.Error("EUR Sell Trades - Expected:", sellTrades, "Got:", sell)
		}
	}

	// JPY
	myData = data.LayoutData{
		Open:       []float64{100.1, 100.2, 100.3, 100.4, 100.5, 100.6, 100.7, 100.8, 100.9},
		High:       []float64{100.3, 100.4, 100.5, 100.6, 100.7, 100.8, 100.9, 101., 101.1},
		Low:        []float64{100.0, 100.1, 100.2, 100.3, 100.4, 100.5, 100.6, 100.7, 100.8},
		Close:      []float64{100.2, 100.3, 100.4, 100.5, 100.6, 100.7, 100.8, 100.9, 101.},
		PriceTf:    []float64{100.1, 100.2, 100.2, 100.3, 100.3, 100.3, 100.4, 100.4, 100.4},
		BuyFlag:    []bool{true, true, true, false, false, false, true, true, true},
		SellFlag:   []bool{false, true, true, true, false, false, false, false, false},
		SizeBuy:    3,
		SizeSell:   3,
		IsFixTpsl:  true,
		TpslFix:    0.1,
		MultiplyTp: 3.,
		MultiplySl: 2.,
	}
	balance = 1000.
	jpy = true
	leverage = false
	tot, buy, sell = BacktestMain(&myData, balance, "no-tpsl", jpy, leverage)

	backtestVariables = backtestHeart(balance, myData.SizeBuy, myData.SizeSell, "", jpy, leverage)
	backtestVariables.buyPrice = 1.002
	backtestVariables.tpBuy = 1.003
	financeCalculation(backtestVariables, "buy-tp")
	backtestVariables.sellPrice = 1.003
	backtestVariables.tpSell = 1.004
	financeCalculation(backtestVariables, "sell-tp")
	totalTrades = conversion.RemoveExcessZeros(backtestVariables.totalTrades)
	buyTrades = conversion.RemoveExcessZeros(backtestVariables.buyTrades)
	sellTrades = conversion.RemoveExcessZeros(backtestVariables.sellTrades)
	for i := range tot {
		if tot[i] != totalTrades[i] {
			t.Error("JPY Total Trades - Expected:", totalTrades, "Got:", tot)
		}
	}
	for i := range buy {
		if buy[i] != buyTrades[i] {
			t.Error("JPY Buy Trades - Expected:", buyTrades, "Got:", buy)
		}
	}
	for i := range sell {
		if sell[i] != sellTrades[i] {
			t.Error("JPY Sell Trades - Expected:", sellTrades, "Got:", sell)
		}
	}
}

func TestNettingBacktest(t *testing.T) {
	// EUR
	var myData = data.LayoutData{
		Open:       []float64{1.001, 1.002, 1.003, 1.004, 1.005, 1.006, 1.007},
		High:       []float64{1.003, 1.004, 1.005, 1.006, 1.007, 1.008, 1.009},
		Low:        []float64{1.000, 1.001, 1.002, 1.003, 1.004, 1.005, 1.006},
		Close:      []float64{1.002, 1.003, 1.004, 1.005, 1.006, 1.007, 1.008},
		PriceTf:    []float64{1.001, 1.002, 1.002, 1.003, 1.003, 1.003, 1.003},
		BuyFlag:    []bool{true, true, true, false, false, false, false},
		SellFlag:   []bool{false, true, true, true, false, false, false},
		SizeBuy:    4,
		SizeSell:   4,
		IsFixTpsl:  true,
		TpslFix:    0.001,
		MultiplyTp: 3.,
		MultiplySl: 2.,
	}
	var balance = 1000.
	var jpy = false
	var leverage = false
	var tot, buy, sell = BacktestMain(&myData, balance, "netting", jpy, leverage)

	var backtestVariables = backtestHeart(balance, myData.SizeBuy, myData.SizeSell, "", jpy, leverage)
	var priceBuy = 1.002
	var isBuy = true
	tpslCalculation(backtestVariables, priceBuy, myData.MultiplyTp, myData.MultiplySl, myData.TpslFix, isBuy)
	financeCalculation(backtestVariables, "buy-tp")
	var priceSell = 1.003
	isBuy = false
	tpslCalculation(backtestVariables, priceSell, myData.MultiplyTp, myData.MultiplySl, myData.TpslFix, isBuy)
	financeCalculation(backtestVariables, "sell-sl")
	var totalTrades []float64 = conversion.RemoveExcessZeros(backtestVariables.totalTrades)
	var buyTrades []float64 = conversion.RemoveExcessZeros(backtestVariables.buyTrades)
	var sellTrades []float64 = conversion.RemoveExcessZeros(backtestVariables.sellTrades)
	for i := range tot {
		if tot[i] != totalTrades[i] {
			t.Error("EUR Total Trades - Expected:", totalTrades, "Got:", tot)
			t.Error("EUR Buy Trades - Expected:", buyTrades, "Got:", buy)

			t.Error("EUR Sell Trades - Expected:", sellTrades, "Got:", sell)
		}
	}
	for i := range buy {
		if buy[i] != buyTrades[i] {
			t.Error("EUR Buy Trades - Expected:", buyTrades, "Got:", buy)
		}
	}
	for i := range sell {
		if sell[i] != sellTrades[i] {
			t.Error("EUR Sell Trades - Expected:", sellTrades, "Got:", sell)
		}
	}

	// JPY
	myData = data.LayoutData{
		Open:       []float64{100.1, 100.2, 100.3, 100.4, 100.5, 100.6, 100.7},
		High:       []float64{100.3, 100.4, 100.5, 100.6, 100.7, 100.8, 100.9},
		Low:        []float64{100.0, 100.1, 100.2, 100.3, 100.4, 100.5, 100.6},
		Close:      []float64{100.2, 100.3, 100.4, 100.5, 100.6, 100.7, 100.8},
		PriceTf:    []float64{100.1, 100.2, 100.3, 100.3, 100.3, 100.3, 100.4},
		BuyFlag:    []bool{true, true, true, false, false, false, false},
		SellFlag:   []bool{false, true, true, true, false, false, false},
		SizeBuy:    4,
		SizeSell:   4,
		IsFixTpsl:  true,
		TpslFix:    0.1,
		MultiplyTp: 3.,
		MultiplySl: 2.,
	}
	balance = 1000.
	jpy = true
	leverage = false
	tot, buy, sell = BacktestMain(&myData, balance, "netting", jpy, leverage)

	backtestVariables = backtestHeart(balance, myData.SizeBuy, myData.SizeSell, "", jpy, leverage)
	priceBuy = 100.2
	isBuy = true
	tpslCalculation(backtestVariables, priceBuy, myData.MultiplyTp, myData.MultiplySl, myData.TpslFix, isBuy)
	financeCalculation(backtestVariables, "buy-tp")
	priceBuy = 100.3
	isBuy = false
	tpslCalculation(backtestVariables, priceBuy, myData.MultiplyTp, myData.MultiplySl, myData.TpslFix, isBuy)
	financeCalculation(backtestVariables, "sell-sl")
	totalTrades = conversion.RemoveExcessZeros(backtestVariables.totalTrades)
	buyTrades = conversion.RemoveExcessZeros(backtestVariables.buyTrades)
	sellTrades = conversion.RemoveExcessZeros(backtestVariables.sellTrades)
	for i := range tot {
		if tot[i] != totalTrades[i] {
			t.Error("JPY Total Trades - Expected:", totalTrades, "Got:", tot)
		}
	}
	for i := range buy {
		if buy[i] != buyTrades[i] {
			t.Error("JPY Buy Trades - Expected:", buyTrades, "Got:", buy)
		}
	}
	for i := range sell {
		if sell[i] != sellTrades[i] {
			t.Error("JPY Sell Trades - Expected:", sellTrades, "Got:", sell)
		}
	}
}

func TestHedgingBacktest(t *testing.T) {
	// EUR
	var myData = data.LayoutData{
		Open:       []float64{1.001, 1.002, 1.003, 1.004, 1.005, 1.006},
		High:       []float64{1.003, 1.004, 1.005, 1.006, 1.007, 1.008},
		Low:        []float64{1.000, 1.001, 1.002, 1.003, 1.004, 1.005},
		Close:      []float64{1.002, 1.003, 1.004, 1.005, 1.006, 1.007},
		PriceTf:    []float64{1.001, 1.002, 1.002, 1.002, 1.002, 1.002},
		BuyFlag:    []bool{true, true, true, false, false, false},
		SellFlag:   []bool{false, true, true, true, false, false},
		SizeBuy:    3,
		SizeSell:   3,
		IsFixTpsl:  true,
		TpslFix:    0.001,
		MultiplyTp: 3.,
		MultiplySl: 2.,
	}
	var balance = 1000.
	var jpy = false
	var leverage = false
	var tot, buy, sell = BacktestMain(&myData, balance, "hedging", jpy, leverage)

	var backtestVariables = backtestHeart(balance, myData.SizeBuy, myData.SizeSell, "", jpy, leverage)
	var priceBuy = 1.002
	var isBuy = true
	tpslCalculation(backtestVariables, priceBuy, myData.MultiplyTp, myData.MultiplySl, myData.TpslFix, isBuy)
	financeCalculation(backtestVariables, "buy-tp")
	isBuy = false
	tpslCalculation(backtestVariables, priceBuy, myData.MultiplyTp, myData.MultiplySl, myData.TpslFix, isBuy)
	financeCalculation(backtestVariables, "sell-sl")
	var totalTrades []float64 = conversion.RemoveExcessZeros(backtestVariables.totalTrades)
	var buyTrades []float64 = conversion.RemoveExcessZeros(backtestVariables.buyTrades)
	var sellTrades []float64 = conversion.RemoveExcessZeros(backtestVariables.sellTrades)
	for i := range tot {
		if tot[i] != totalTrades[i] {
			t.Error("EUR Total Trades - Expected:", totalTrades, "Got:", tot)
		}
	}
	for i := range buy {
		if buy[i] != buyTrades[i] {
			t.Error("EUR Buy Trades - Expected:", buyTrades, "Got:", buy)
		}
	}
	for i := range sell {
		if sell[i] != sellTrades[i] {
			t.Error("EUR Sell Trades - Expected:", sellTrades, "Got:", sell)
		}
	}

	// JPY
	myData = data.LayoutData{
		Open:       []float64{100.1, 100.2, 100.3, 100.4, 100.5, 100.6},
		High:       []float64{100.3, 100.4, 100.5, 100.6, 100.7, 100.8},
		Low:        []float64{100.0, 100.1, 100.2, 100.3, 100.4, 100.5},
		Close:      []float64{100.2, 100.3, 100.4, 100.5, 100.6, 100.7},
		PriceTf:    []float64{100.1, 100.2, 100.2, 100.2, 100.2, 100.2},
		BuyFlag:    []bool{true, true, true, false, false, false},
		SellFlag:   []bool{false, true, true, true, false, false},
		SizeBuy:    3,
		SizeSell:   3,
		IsFixTpsl:  true,
		TpslFix:    0.1,
		MultiplyTp: 3.,
		MultiplySl: 2.,
	}
	balance = 1000.
	jpy = true
	leverage = false
	tot, buy, sell = BacktestMain(&myData, balance, "hedging", jpy, leverage)

	backtestVariables = backtestHeart(balance, myData.SizeBuy, myData.SizeSell, "", jpy, leverage)
	priceBuy = 100.2
	isBuy = true
	tpslCalculation(backtestVariables, priceBuy, myData.MultiplyTp, myData.MultiplySl, myData.TpslFix, isBuy)
	financeCalculation(backtestVariables, "buy-tp")
	isBuy = false
	tpslCalculation(backtestVariables, priceBuy, myData.MultiplyTp, myData.MultiplySl, myData.TpslFix, isBuy)
	financeCalculation(backtestVariables, "sell-sl")
	totalTrades = conversion.RemoveExcessZeros(backtestVariables.totalTrades)
	buyTrades = conversion.RemoveExcessZeros(backtestVariables.buyTrades)
	sellTrades = conversion.RemoveExcessZeros(backtestVariables.sellTrades)
	for i := range tot {
		if tot[i] != totalTrades[i] {
			t.Error("JPY Total Trades - Expected:", totalTrades, "Got:", tot)
		}
	}
	for i := range buy {
		if buy[i] != buyTrades[i] {
			t.Error("JPY Buy Trades - Expected:", buyTrades, "Got:", buy)
		}
	}
	for i := range sell {
		if sell[i] != sellTrades[i] {
			t.Error("JPY Sell Trades - Expected:", sellTrades, "Got:", sell)
		}
	}
}

func TestTakeProfitFimathe(t *testing.T) {
	var myData = data.LayoutData{
		SizeBuy:  3,
		SizeSell: 3,
	}

	var balance = 1000.
	var jpy = false
	var leverage = false
	var backtestVariables = backtestHeart(balance, myData.SizeBuy, myData.SizeSell, "", jpy, leverage)

	backtestVariables.upperChannel = 1.123
	myData.MultiplyTpChannel = 2.
	myData.ChannelInPips = 0.001
	var buy = true

	takeProfitFimathe(backtestVariables, myData.MultiplyTpChannel, myData.ChannelInPips, buy)
	if backtestVariables.tpBuy != 1.125 {
		t.Error("Buy Non JPY - Expected:", 1.125, "Got:", backtestVariables.tpBuy)
	}

	backtestVariables.lowerChannel = 1.123
	buy = false
	takeProfitFimathe(backtestVariables, myData.MultiplyTpChannel, myData.ChannelInPips, buy)
	if backtestVariables.tpSell != 1.121 {
		t.Error("Sell Non JPY - Expected:", 1.121, "Got:", backtestVariables.tpSell)
	}

	// JPY
	backtestVariables.upperChannel = 111.5
	myData.MultiplyTpChannel = 2.
	myData.ChannelInPips = 0.1
	jpy = true
	buy = true

	takeProfitFimathe(backtestVariables, myData.MultiplyTpChannel, myData.ChannelInPips, buy)
	if backtestVariables.tpBuy != 111.7 {
		t.Error("Buy JPY - Expected:", 111.7, "Got:", backtestVariables.tpBuy)
	}

	backtestVariables.lowerChannel = 111.5
	buy = false
	takeProfitFimathe(backtestVariables, myData.MultiplyTpChannel, myData.ChannelInPips, buy)
	if backtestVariables.tpSell != 111.3 {
		t.Error("Sell JPY - Expected:", 111.3, "Got:", backtestVariables.tpSell)
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

func TestUpdateFimatheChannels(t *testing.T) {
	var jpy = false
	var leverage = false
	var backtestVariables = backtestHeart(1000, 10, 10, "", jpy, leverage)

	var price = 1.120
	var channelInPips = 0.005
	backtestVariables.upperChannel = 1.115
	backtestVariables.lowerChannel = 1.105
	updateFimatheChannels(backtestVariables, price, channelInPips)
	if backtestVariables.upperChannel != 1.125 {
		t.Error("EUR - UpperChannel - Expected:", 1.125, "Got:", backtestVariables.upperChannel)
	}
	if backtestVariables.lowerChannel != 1.12 {
		t.Error("EUR - LowerChannel - Expected:", 1.12, "Got:", backtestVariables.lowerChannel)
	}
	price = 1.100
	channelInPips = 0.005
	backtestVariables.upperChannel = 1.115
	backtestVariables.lowerChannel = 1.105
	updateFimatheChannels(backtestVariables, price, channelInPips)
	if backtestVariables.upperChannel != 1.100 {
		t.Error("EUR - UpperChannel - Expected:", 1.100, "Got:", backtestVariables.upperChannel)
	}
	if backtestVariables.lowerChannel != 1.095 {
		t.Error("EUR - LowerChannel - Expected:", 1.095, "Got:", backtestVariables.lowerChannel)
	}

	jpy = true
	leverage = false
	backtestVariables = backtestHeart(1000, 10, 10, "", jpy, leverage)

	price = 112.0
	channelInPips = 0.5
	backtestVariables.upperChannel = 111.5
	backtestVariables.lowerChannel = 110.5
	updateFimatheChannels(backtestVariables, price, channelInPips)
	if backtestVariables.upperChannel != 112.5 {
		t.Error("EUR - UpperChannel - Expected:", 112.5, "Got:", backtestVariables.upperChannel)
	}
	if backtestVariables.lowerChannel != 112. {
		t.Error("EUR - LowerChannel - Expected:", 112., "Got:", backtestVariables.lowerChannel)
	}
	price = 110.0
	channelInPips = 000.5
	backtestVariables.upperChannel = 111.5
	backtestVariables.lowerChannel = 110.5
	updateFimatheChannels(backtestVariables, price, channelInPips)
	if backtestVariables.upperChannel != 110.0 {
		t.Error("EUR - UpperChannel - Expected:", 110.0, "Got:", backtestVariables.upperChannel)
	}
	if backtestVariables.lowerChannel != 109.5 {
		t.Error("EUR - LowerChannel - Expected:", 109.5, "Got:", backtestVariables.lowerChannel)
	}

}

func TestTpslCalculation(t *testing.T) {
	var jpy = false
	var leverage = false
	var backtestVariables = backtestHeart(1000, 10, 10, "", jpy, leverage)

	var multiplyValue = 0.001
	var price = 1.123
	var multiplyTp = 2.5
	var multiplySl = 1.25
	var buy = true
	tpslCalculation(backtestVariables, price, multiplyTp, multiplySl, multiplyValue, buy)
	if backtestVariables.tpBuy != 1.1255 {
		t.Error("EUR - BUY - Expected:", 1.1255, "Got:", backtestVariables.tpBuy)
	} else if backtestVariables.slBuy != 1.12175 {
		t.Error("EUR - BUY - Expected:", 1.12175, "Got:", backtestVariables.slBuy)
	}
	price = 1.123
	multiplyTp = 2.5
	multiplySl = 1.25
	jpy = false
	buy = false
	tpslCalculation(backtestVariables, price, multiplyTp, multiplySl, multiplyValue, buy)
	if backtestVariables.tpSell != 1.1205 {
		t.Error("EUR - SELL - Expected:", 1.1205, "Got:", backtestVariables.tpSell)
	} else if backtestVariables.slSell != 1.12425 {
		t.Error("EUR - SELL - Expected:", 1.12425, "Got:", backtestVariables.slSell)
	}

	multiplyValue = 0.1
	price = 123.456
	multiplyTp = 2.5
	multiplySl = 1.25
	jpy = true
	buy = true
	tpslCalculation(backtestVariables, price, multiplyTp, multiplySl, multiplyValue, buy)
	if backtestVariables.tpBuy != 123.706 {
		t.Error("JPY - BUY - Expected:", 123.706, "Got:", backtestVariables.tpBuy)
	} else if backtestVariables.slBuy != 123.331 {
		t.Error("JPY - BUY - Expected:", 123.331, "Got:", backtestVariables.slBuy)
	}
	price = 123.456
	multiplyTp = 2.5
	multiplySl = 1.25
	jpy = true
	buy = false
	tpslCalculation(backtestVariables, price, multiplyTp, multiplySl, multiplyValue, buy)
	if backtestVariables.tpSell != 123.206 {
		t.Error("JPY - SELL - Expected:", 123.206, "Got:", backtestVariables.tpSell)
	} else if backtestVariables.slSell != 123.581 {
		t.Error("JPY - SELL - Expected:", 123.581, "Got:", backtestVariables.slSell)
	}
}

func TestFinanceCalculation(t *testing.T) {
	var leverage = false
	var jpy = false

	var backtestVariables = backtestHeart(1000., 10, 10, "", jpy, leverage)

	backtestVariables.tpBuy = 1.124
	backtestVariables.buyPrice = 1.123
	financeCalculation(backtestVariables, "buy-tp")
	if backtestVariables.buyResult != 0.81968 {
		t.Error("EUR BUY TP - Expected:", 0.81968, "Got:", backtestVariables.buyResult)
	}
	if backtestVariables.capital != 1000.82 {
		t.Error("EUR BUY TP - Expected:", 1000.82, "Got:", backtestVariables.capital)
	}

	backtestVariables = backtestHeart(1000., 10, 10, "", jpy, leverage)
	backtestVariables.slBuy = 1.122
	backtestVariables.buyPrice = 1.123
	financeCalculation(backtestVariables, "buy-sl")
	if backtestVariables.buyResult != -0.96127 {
		t.Error("EUR BUY SL - Expected:", -0.96127, "Got:", backtestVariables.buyResult)
	}
	if backtestVariables.capital != 999.04 {
		t.Error("EUR BUY SL - Expected:", 999.04, "Got:", backtestVariables.capital)
	}

	backtestVariables = backtestHeart(1000., 10, 10, "", jpy, leverage)
	backtestVariables.tpSell = 1.122
	backtestVariables.sellPrice = 1.123
	financeCalculation(backtestVariables, "sell-tp")
	if backtestVariables.sellResult != 0.82127 {
		t.Error("EUR SELL TP - Expected:", 0.82127, "Got:", backtestVariables.sellResult)
	}
	if backtestVariables.capital != 1000.82 {
		t.Error("EUR SELL TP - Expected:", 1000.82, "Got:", backtestVariables.capital)
	}

	backtestVariables = backtestHeart(1000., 10, 10, "", jpy, leverage)
	backtestVariables.slSell = 1.124
	backtestVariables.sellPrice = 1.123
	financeCalculation(backtestVariables, "sell-sl")
	if backtestVariables.sellResult != -0.95968 {
		t.Error("EUR SELL SL - Expected:", -0.95968, "Got:", backtestVariables.sellResult)
	}
	if backtestVariables.capital != 999.04 {
		t.Error("EUR SELL SL - Expected:", 999.04, "Got:", backtestVariables.capital)
	}

	// JPY
	jpy = true
	backtestVariables = backtestHeart(1000., 10, 10, "", jpy, leverage)
	backtestVariables.tpBuy = 127.204
	backtestVariables.buyPrice = 127.124
	jpy = true
	financeCalculation(backtestVariables, "buy-tp")
	if backtestVariables.buyResult != 0.559 {
		t.Error("JPY BUY TP - Expected:", 0.559, "Got:", backtestVariables.buyResult)
	}
	if backtestVariables.capital != 1000.56 {
		t.Error("JPY BUY TP - Expected:", 1000.56, "Got:", backtestVariables.capital)
	}

	backtestVariables = backtestHeart(1000., 10, 10, "", jpy, leverage)
	backtestVariables.slBuy = 127.05
	backtestVariables.buyPrice = 127.1
	jpy = true
	financeCalculation(backtestVariables, "buy-sl")
	if backtestVariables.buyResult != -0.464 {
		t.Error("JPY BUY SL - Expected:", -0.464, "Got:", backtestVariables.buyResult)
	}
	if backtestVariables.capital != 999.54 {
		t.Error("JPY BUY SL - Expected:", 999.54, "Got:", backtestVariables.capital)
	}

	backtestVariables = backtestHeart(1000., 10, 10, "", jpy, leverage)
	backtestVariables.tpSell = 127.05
	backtestVariables.sellPrice = 127.1
	jpy = true
	financeCalculation(backtestVariables, "sell-tp")
	if backtestVariables.sellResult != 0.324 {
		t.Error("JPY SELL TP - Expected:", 0.324, "Got:", backtestVariables.sellResult)
	}
	if backtestVariables.capital != 1000.32 {
		t.Error("JPY SELL TP - Expected:", 1000.32, "Got:", backtestVariables.capital)
	}

	backtestVariables = backtestHeart(1000., 10, 10, "", jpy, leverage)
	backtestVariables.slSell = 127.1
	backtestVariables.sellPrice = 127.05
	jpy = true
	financeCalculation(backtestVariables, "sell-sl")
	if backtestVariables.sellResult != -0.463 {
		t.Error("JPY SELL SL - Expected:", -0.463, "Got:", backtestVariables.sellResult)
	}
	if backtestVariables.capital != 999.54 {
		t.Error("JPY SELL SL - Expected:", 999.54, "Got:", backtestVariables.capital)
	}

}

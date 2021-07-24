package backtest

import "github.com/SrJMaia/expert-advisor/program/conversion"

type backtestVariables struct {
	buyPrice       float64
	sellPrice      float64
	tpSell         float64
	slSell         float64
	tpBuy          float64
	slBuy          float64
	buyResult      float64
	sellResult     float64
	capital        float64
	initialCapital float64
	buyFlag        bool
	sellFlag       bool
	universalFlag  bool
	updateBuy      bool
	updateSell     bool
	iBuy           uint32
	iSell          uint32
	iTotalTrades   uint32
	totalTrades    []float64
	buyTrades      []float64
	sellTrades     []float64
}

func backtestHeart(balance float64, sizeBuy uint32, sizeSell uint32) backtestVariables {
	var backtestMain = backtestVariables{
		buyPrice:       0.,
		sellPrice:      0.,
		tpSell:         0.,
		slSell:         0.,
		tpBuy:          0.,
		slBuy:          0.,
		buyResult:      0.,
		sellResult:     0.,
		capital:        balance,
		initialCapital: balance,
		buyFlag:        true,
		sellFlag:       true,
		universalFlag:  true,
		iBuy:           1,
		iSell:          1,
		iTotalTrades:   1,
		totalTrades:    make([]float64, sizeBuy+sizeSell),
		buyTrades:      make([]float64, sizeBuy),
		sellTrades:     make([]float64, sizeSell),
	}

	backtestMain.totalTrades[0] = balance
	backtestMain.buyTrades[0] = balance
	backtestMain.sellTrades[0] = balance

	return backtestMain
}

func updateBacktestHeart(backtestMain *backtestVariables, action string) {
	if action == "buy" {
		(*backtestMain).totalTrades[(*backtestMain).iTotalTrades] = (*backtestMain).capital
		(*backtestMain).buyTrades[(*backtestMain).iBuy] = conversion.Round((*backtestMain).buyTrades[(*backtestMain).iBuy-1]+(*backtestMain).buyResult, 2)
		(*backtestMain).iTotalTrades++
		(*backtestMain).iBuy++
		(*backtestMain).buyFlag = true
		(*backtestMain).universalFlag = true
		(*backtestMain).updateBuy = false
		(*backtestMain).tpBuy = 0.
		(*backtestMain).slBuy = 0.
		(*backtestMain).buyPrice = 0.
	} else if action == "sell" {
		(*backtestMain).totalTrades[(*backtestMain).iTotalTrades] = (*backtestMain).capital
		(*backtestMain).sellTrades[(*backtestMain).iSell] = conversion.Round((*backtestMain).sellTrades[(*backtestMain).iSell-1]+(*backtestMain).sellResult, 2)
		(*backtestMain).iTotalTrades++
		(*backtestMain).iSell++
		(*backtestMain).sellFlag = true
		(*backtestMain).universalFlag = true
		(*backtestMain).updateSell = false
		(*backtestMain).tpSell = 0.
		(*backtestMain).slSell = 0.
		(*backtestMain).sellPrice = 0.
	}
}

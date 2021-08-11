package backtest

import (
	"log"

	"github.com/SrJMaia/expert-advisor/program/conversion"
)

type backtestVariablesStruct struct {
	buyPrice                float64
	sellPrice               float64
	tpSell                  float64
	slSell                  float64
	tpBuy                   float64
	slBuy                   float64
	buyResult               float64
	sellResult              float64
	capital                 float64
	initialCapital          float64
	upperChannel            float64
	lowerChannel            float64
	buyFlag                 bool
	sellFlag                bool
	universalFlag           bool
	updateBuy               bool
	updateSell              bool
	isJpy                   bool
	leverage                bool
	totalTradesInOneDay     int32
	flagLowerChannelCounter uint32
	flagUpperChannelCounter uint32
	iBuy                    uint32
	iSell                   uint32
	iTotalTrades            uint32
	totalTrades             []float64
	buyTrades               []float64
	sellTrades              []float64
}

func backtestHeart(balance float64, sizeBuy uint32, sizeSell uint32, optimationType string, jpy bool, leverage bool) *backtestVariablesStruct {
	if optimationType == "fimathe" || optimationType == "fimathe-with-strategy" {
		sizeBuy = 100000
		sizeSell = 100000
	} else if optimationType != "fimathe" && optimationType != "fimathe-with-strategy" && sizeBuy == 0 || sizeSell == 0 {
		log.Panic("Error in backtestHeart - Size Buy:", sizeBuy, "Size Sell:", sizeBuy)
	}
	var backtestVariables = backtestVariablesStruct{
		buyPrice:            0.,
		sellPrice:           0.,
		tpSell:              0.,
		slSell:              0.,
		tpBuy:               0.,
		slBuy:               0.,
		buyResult:           0.,
		sellResult:          0.,
		capital:             balance,
		initialCapital:      balance,
		buyFlag:             true,
		sellFlag:            true,
		universalFlag:       true,
		isJpy:               jpy,
		leverage:            leverage,
		totalTradesInOneDay: 0,
		iBuy:                1,
		iSell:               1,
		iTotalTrades:        1,
		totalTrades:         make([]float64, sizeBuy+sizeSell),
		buyTrades:           make([]float64, sizeBuy),
		sellTrades:          make([]float64, sizeSell),
	}

	backtestVariables.totalTrades[0] = balance
	backtestVariables.buyTrades[0] = balance
	backtestVariables.sellTrades[0] = balance

	return &backtestVariables
}

func updateBacktestHeart(backtestVariables *backtestVariablesStruct, action string) {
	if action == "buy" {
		(*backtestVariables).totalTrades[(*backtestVariables).iTotalTrades] = (*backtestVariables).capital
		(*backtestVariables).buyTrades[(*backtestVariables).iBuy] = conversion.Round((*backtestVariables).buyTrades[(*backtestVariables).iBuy-1]+(*backtestVariables).buyResult, 2)
		(*backtestVariables).iTotalTrades++
		(*backtestVariables).iBuy++
		(*backtestVariables).buyFlag = true
		(*backtestVariables).universalFlag = true
		(*backtestVariables).updateBuy = false
		(*backtestVariables).tpBuy = 0.
		(*backtestVariables).slBuy = 0.
		(*backtestVariables).buyPrice = 0.
	} else if action == "sell" {
		(*backtestVariables).totalTrades[(*backtestVariables).iTotalTrades] = (*backtestVariables).capital
		(*backtestVariables).sellTrades[(*backtestVariables).iSell] = conversion.Round((*backtestVariables).sellTrades[(*backtestVariables).iSell-1]+(*backtestVariables).sellResult, 2)
		(*backtestVariables).iTotalTrades++
		(*backtestVariables).iSell++
		(*backtestVariables).sellFlag = true
		(*backtestVariables).universalFlag = true
		(*backtestVariables).updateSell = false
		(*backtestVariables).tpSell = 0.
		(*backtestVariables).slSell = 0.
		(*backtestVariables).sellPrice = 0.
	}
}

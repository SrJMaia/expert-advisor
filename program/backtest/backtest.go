package backtest

import (
	"log"

	"github.com/SrJMaia/expert-advisor/program/check"
	"github.com/SrJMaia/expert-advisor/program/conversion"
	"github.com/SrJMaia/expert-advisor/program/data"
)

func updateFimatheChannels(backtestVariables *backtestVariablesStruct, price float64, channelInPips float64) {
	// flag zona neutra?
	if price > (*backtestVariables).upperChannel {
		for price >= (*backtestVariables).upperChannel {
			(*backtestVariables).lowerChannel = conversion.RoundIsJpy((*backtestVariables).upperChannel, (*backtestVariables).isJpy)
			(*backtestVariables).upperChannel = conversion.RoundIsJpy((*backtestVariables).upperChannel+channelInPips, (*backtestVariables).isJpy)
			(*backtestVariables).flagLowerChannelCounter = 0
			(*backtestVariables).flagUpperChannelCounter++
		}
		if !backtestVariables.universalFlag && !backtestVariables.buyFlag {
			backtestVariables.slBuy = conversion.RoundIsJpy((*backtestVariables).lowerChannel-channelInPips, (*backtestVariables).isJpy)
		}
	} else if price < (*backtestVariables).lowerChannel {
		for price <= (*backtestVariables).lowerChannel {
			(*backtestVariables).upperChannel = conversion.RoundIsJpy((*backtestVariables).lowerChannel, (*backtestVariables).isJpy)
			(*backtestVariables).lowerChannel = conversion.RoundIsJpy((*backtestVariables).lowerChannel-channelInPips, (*backtestVariables).isJpy)
			(*backtestVariables).flagUpperChannelCounter = 0
			(*backtestVariables).flagLowerChannelCounter++
		}
		if !backtestVariables.universalFlag && !backtestVariables.sellFlag {
			backtestVariables.slSell = conversion.RoundIsJpy((*backtestVariables).upperChannel+channelInPips, (*backtestVariables).isJpy)
		}
	}
}

func financeCalculation(backtestVariables *backtestVariablesStruct, calculationMode string) {
	var comission float64
	var lot float64
	if (*backtestVariables).leverage && (*backtestVariables).capital > 3*(*backtestVariables).initialCapital {
		lot = conversion.Round(((*backtestVariables).capital - (*backtestVariables).initialCapital), -3)
	} else if (*backtestVariables).totalTradesInOneDay >= 1 {
		lot = (*backtestVariables).initialCapital * 1
	} else {
		lot = (*backtestVariables).initialCapital
	}
	comission = lot / 1000 * .07

	if calculationMode == "buy-tp" {
		(*backtestVariables).buyResult = conversion.RoundIsJpy((lot*((*backtestVariables).tpBuy-(*backtestVariables).buyPrice))/(*backtestVariables).tpBuy-comission, (*backtestVariables).isJpy)
		(*backtestVariables).capital = conversion.Round((*backtestVariables).buyResult+(*backtestVariables).capital, 2)
		check.MyCheckingFinanceCalculation((*backtestVariables).capital, (*backtestVariables).tpBuy, (*backtestVariables).buyPrice, (*backtestVariables).buyResult)
		updateBacktestHeart(backtestVariables, "buy")
		(*backtestVariables).totalTradesInOneDay++
	} else if calculationMode == "buy-sl" {
		(*backtestVariables).buyResult = conversion.RoundIsJpy((lot*((*backtestVariables).slBuy-(*backtestVariables).buyPrice))/(*backtestVariables).slBuy-comission, (*backtestVariables).isJpy)
		(*backtestVariables).capital = conversion.Round((*backtestVariables).buyResult+(*backtestVariables).capital, 2)
		check.MyCheckingFinanceCalculation((*backtestVariables).capital, (*backtestVariables).slBuy, (*backtestVariables).buyPrice, (*backtestVariables).buyResult)
		updateBacktestHeart(backtestVariables, "buy")
		(*backtestVariables).totalTradesInOneDay--
	} else if calculationMode == "sell-tp" {
		(*backtestVariables).sellResult = conversion.RoundIsJpy((lot*((*backtestVariables).sellPrice-(*backtestVariables).tpSell))/(*backtestVariables).tpSell-comission, (*backtestVariables).isJpy)
		(*backtestVariables).capital = conversion.Round((*backtestVariables).sellResult+(*backtestVariables).capital, 2)
		check.MyCheckingFinanceCalculation((*backtestVariables).capital, (*backtestVariables).tpSell, (*backtestVariables).sellPrice, (*backtestVariables).sellResult)
		updateBacktestHeart(backtestVariables, "sell")
		(*backtestVariables).totalTradesInOneDay++
	} else if calculationMode == "sell-sl" {
		(*backtestVariables).sellResult = conversion.RoundIsJpy((lot*((*backtestVariables).sellPrice-(*backtestVariables).slSell))/(*backtestVariables).slSell-comission, (*backtestVariables).isJpy)
		(*backtestVariables).capital = conversion.Round((*backtestVariables).sellResult+(*backtestVariables).capital, 2)
		check.MyCheckingFinanceCalculation((*backtestVariables).capital, (*backtestVariables).slSell, (*backtestVariables).sellPrice, (*backtestVariables).sellResult)
		updateBacktestHeart(backtestVariables, "sell")
		(*backtestVariables).totalTradesInOneDay--
	} else {
		log.Panic("Finance Calculation Wrong Mode")
	}

}

func tpslCalculation(backtestVariables *backtestVariablesStruct, price float64, multiplyTp float64, multiplySl float64, tpSlValue float64, buy bool) {
	if buy {
		(*backtestVariables).buyPrice = price
		(*backtestVariables).tpBuy = conversion.RoundIsJpy(price+(multiplyTp*tpSlValue), (*backtestVariables).isJpy)
		(*backtestVariables).slBuy = conversion.RoundIsJpy(price-(multiplySl*tpSlValue), (*backtestVariables).isJpy)
		(*backtestVariables).totalTradesInOneDay++
	} else {
		(*backtestVariables).sellPrice = price
		(*backtestVariables).tpSell = conversion.RoundIsJpy(price-(multiplyTp*tpSlValue), (*backtestVariables).isJpy)
		(*backtestVariables).slSell = conversion.RoundIsJpy(price+(multiplySl*tpSlValue), (*backtestVariables).isJpy)
		(*backtestVariables).totalTradesInOneDay++
	}
}

func checkTpslExit(open float64, high float64, low float64, close float64, tpsl float64, higher bool) bool {
	if higher {
		return open >= tpsl || high >= tpsl || low >= tpsl || close >= tpsl
	}
	return open <= tpsl || high <= tpsl || low <= tpsl || close <= tpsl
}

func takeProfitFimathe(backtestVariables *backtestVariablesStruct, multiply float64, channelInPips float64, buy bool) {
	if buy {
		(*backtestVariables).tpBuy = conversion.RoundIsJpy((*backtestVariables).upperChannel+(channelInPips*2), (*backtestVariables).isJpy)
	} else {
		(*backtestVariables).tpSell = conversion.RoundIsJpy((*backtestVariables).lowerChannel-(channelInPips*2), (*backtestVariables).isJpy)
	}
}

func BacktestMain(dt *data.LayoutData, balance float64, backtestType string, jpy bool, leverage bool) ([]float64, []float64, []float64) {
	// Canal =~ 10 a 200 pips ou 0.0010 a 0.0200
	// StartPoint range(+halfChannel até -halfChannel)
	// multiplyTpChannel = range(1,6)
	// 2 Tipos, tp sl fixos ou trailing stop

	var backtestVariables = backtestHeart(balance, dt.SizeBuy, dt.SizeSell, backtestType, jpy, leverage)

	if (*dt).StartPointInPips == 0 {
		backtestVariables.upperChannel = conversion.RoundIsJpy((*dt).Open[0]+((*dt).ChannelInPips/2.), jpy)
		backtestVariables.lowerChannel = conversion.RoundIsJpy((*dt).Open[0]-((*dt).ChannelInPips/2.), jpy)
	} else {
		backtestVariables.upperChannel = conversion.RoundIsJpy((((*dt).StartPointInPips) + (*dt).Open[0]), jpy)
		backtestVariables.lowerChannel = conversion.RoundIsJpy(backtestVariables.upperChannel-(*dt).ChannelInPips, jpy)
	}

	// var dayTeste = false

	for i := 1; i < len((*dt).Open); i++ {

		// if (*dt).TimeWeekDays[i] != (*dt).TimeWeekDays[i-1] && !dayTeste {
		// 	dayTeste = true
		// }
		// if (*dt).TimeWeekDays[i] != (*dt).TimeWeekDays[i-1] {
		// 	(*backtestVariables).totalTradesInOneDay = 0
		// }

		// if (*dt).TimeHour[i] > 10 && (*dt).TimeHour[i] < 18 { //  && (*backtestVariables).totalTradesInOneDay < 2
		if (*dt).PriceTf[i] != (*dt).PriceTf[i-1] {
			if backtestType == "netting" {
				backtestNetting(dt, backtestVariables, i, "open")
			} else if backtestType == "hedging" {
				backtestHedging(dt, backtestVariables, i, "open")
			} else if backtestType == "no-tpsl" {
				backtestNoTpsl(dt, backtestVariables, i)
			} else if backtestType == "fimathe" {
				backtestFimathe(dt, backtestVariables, i, "open")
			} else if backtestType == "fimathe-with-strategy" {
				backtestFimatheWithStrategy(dt, backtestVariables, i, "open")
			}
		} else if (*dt).PriceTf[i] == (*dt).PriceTf[i-1] {
			if backtestType == "netting" {
				backtestNetting(dt, backtestVariables, i, "close")
			} else if backtestType == "hedging" {
				backtestHedging(dt, backtestVariables, i, "close")
			} else if backtestType == "fimathe" {
				backtestFimathe(dt, backtestVariables, i, "close")
			} else if backtestType == "fimathe-with-strategy" {
				backtestFimatheWithStrategy(dt, backtestVariables, i, "close")
			}
		}
		// }
		// if (*dt).TimeHour[i] >= 20 && (*dt).TimeHour[i] <= 23 && (*dt).TimeWeekDays[i] == 5 {
		// 	if !backtestVariables.buyFlag {
		// 		backtestVariables.tpBuy = (*dt).Open[i]
		// 		financeCalculation(backtestVariables, "buy-tp")
		// 	}
		// 	if !backtestVariables.sellFlag {
		// 		backtestVariables.tpSell = (*dt).Open[i]
		// 		financeCalculation(backtestVariables, "sell-tp")
		// 	}
		// 	dayTeste = false
		// }
	}

	var totalTrades = conversion.RemoveExcessZeros(backtestVariables.totalTrades)
	var buyTrades = conversion.RemoveExcessZeros(backtestVariables.buyTrades)
	var sellTrades = conversion.RemoveExcessZeros(backtestVariables.sellTrades)

	return totalTrades, buyTrades, sellTrades

}

func backtestFimatheWithStrategy(dt *data.LayoutData, backtestVariables *backtestVariablesStruct, index int, mode string) {
	if mode == "open" {
		updateFimatheChannels(backtestVariables, (*dt).PriceTf[index], (*dt).ChannelInPips)
		if (*dt).PriceTf[index] > (*backtestVariables).upperChannel && backtestVariables.buyFlag && backtestVariables.universalFlag && (*dt).BuyFlag[index] {
			backtestVariables.buyPrice = (*dt).PriceTf[index]
			takeProfitFimathe(backtestVariables, (*dt).MultiplyTpChannel, (*dt).ChannelInPips, true)
			backtestVariables.slBuy = conversion.RoundIsJpy((*backtestVariables).lowerChannel-((*dt).ChannelInPips*1.5), (*backtestVariables).isJpy)
			backtestVariables.buyFlag = false
			backtestVariables.universalFlag = false
		} else if (*dt).PriceTf[index] < (*backtestVariables).lowerChannel && backtestVariables.sellFlag && backtestVariables.universalFlag && (*dt).SellFlag[index] {
			backtestVariables.sellPrice = (*dt).PriceTf[index]
			takeProfitFimathe(backtestVariables, (*dt).MultiplyTpChannel, (*dt).ChannelInPips, false)
			backtestVariables.slSell = conversion.RoundIsJpy((*backtestVariables).upperChannel+((*dt).ChannelInPips*1.5), (*backtestVariables).isJpy)
			backtestVariables.sellFlag = false
			backtestVariables.universalFlag = false
		}
	} else if mode == "close" {
		if !backtestVariables.universalFlag && !backtestVariables.buyFlag {
			if checkTpslExit((*dt).Open[index], (*dt).High[index], (*dt).Low[index], (*dt).Close[index], backtestVariables.slBuy, false) {
				financeCalculation(backtestVariables, "buy-sl")
			} else if checkTpslExit((*dt).Open[index], (*dt).High[index], (*dt).Low[index], (*dt).Close[index], backtestVariables.tpBuy, true) {
				financeCalculation(backtestVariables, "buy-tp")
			}
		}
		if !backtestVariables.universalFlag && !backtestVariables.sellFlag {
			if checkTpslExit((*dt).Open[index], (*dt).High[index], (*dt).Low[index], (*dt).Close[index], backtestVariables.slSell, true) {
				financeCalculation(backtestVariables, "sell-sl")
			} else if checkTpslExit((*dt).Open[index], (*dt).High[index], (*dt).Low[index], (*dt).Close[index], backtestVariables.tpSell, false) {
				financeCalculation(backtestVariables, "sell-tp")
			}
		}
	} else {
		log.Panic("Fimathe With Strategy Mode Wrong.")
	}
}

func backtestFimathe(dt *data.LayoutData, backtestVariables *backtestVariablesStruct, index int, mode string) {
	if mode == "open" {
		updateFimatheChannels(backtestVariables, (*dt).PriceTf[index], (*dt).ChannelInPips)
		if (*dt).PriceTf[index] > (*backtestVariables).lowerChannel && (*dt).PriceTf[index-1] < (*backtestVariables).lowerChannel && backtestVariables.buyFlag && backtestVariables.universalFlag && (*backtestVariables).flagUpperChannelCounter == 2 {
			backtestVariables.buyPrice = (*dt).PriceTf[index]
			takeProfitFimathe(backtestVariables, (*dt).MultiplyTpChannel, (*dt).ChannelInPips, true)
			backtestVariables.slBuy = conversion.RoundIsJpy((*backtestVariables).lowerChannel-((*dt).ChannelInPips*1.5), (*backtestVariables).isJpy)
			backtestVariables.buyFlag = false
			backtestVariables.universalFlag = false
		} else if (*dt).PriceTf[index] < (*backtestVariables).upperChannel && (*dt).PriceTf[index-1] > (*backtestVariables).upperChannel && backtestVariables.sellFlag && backtestVariables.universalFlag && (*backtestVariables).flagLowerChannelCounter == 2 {
			backtestVariables.sellPrice = (*dt).PriceTf[index]
			takeProfitFimathe(backtestVariables, (*dt).MultiplyTpChannel, (*dt).ChannelInPips, false)
			backtestVariables.slSell = conversion.RoundIsJpy((*backtestVariables).upperChannel+((*dt).ChannelInPips*1.5), (*backtestVariables).isJpy)
			backtestVariables.sellFlag = false
			backtestVariables.universalFlag = false
		}
	} else if mode == "close" {
		if !backtestVariables.universalFlag && !backtestVariables.buyFlag {
			if checkTpslExit((*dt).Open[index], (*dt).High[index], (*dt).Low[index], (*dt).Close[index], backtestVariables.slBuy, false) {
				financeCalculation(backtestVariables, "buy-sl")
			} else if checkTpslExit((*dt).Open[index], (*dt).High[index], (*dt).Low[index], (*dt).Close[index], backtestVariables.tpBuy, true) {
				financeCalculation(backtestVariables, "buy-tp")
			}
		}
		if !backtestVariables.universalFlag && !backtestVariables.sellFlag {
			if checkTpslExit((*dt).Open[index], (*dt).High[index], (*dt).Low[index], (*dt).Close[index], backtestVariables.slSell, true) {
				financeCalculation(backtestVariables, "sell-sl")
			} else if checkTpslExit((*dt).Open[index], (*dt).High[index], (*dt).Low[index], (*dt).Close[index], backtestVariables.tpSell, false) {
				financeCalculation(backtestVariables, "sell-tp")
			}
		}
	} else {
		log.Panic("Fimathe Mode Wrong.")
	}
}

func backtestNoTpsl(dt *data.LayoutData, backtestVariables *backtestVariablesStruct, index int) {
	if (*dt).BuyFlag[index] {
		if !(*backtestVariables).sellFlag {
			(*backtestVariables).buyPrice = (*dt).PriceTf[index]
			(*backtestVariables).buyFlag = false
			(*backtestVariables).sellFlag = true
			(*backtestVariables).tpSell = (*dt).PriceTf[index]
			financeCalculation(backtestVariables, "sell-tp")
			(*backtestVariables).totalTradesInOneDay++
			return
		} else if (*backtestVariables).buyFlag {
			(*backtestVariables).buyPrice = (*dt).PriceTf[index]
			(*backtestVariables).buyFlag = false
			(*backtestVariables).totalTradesInOneDay++
			return
		}
	}
	if (*dt).SellFlag[index] {
		if !(*backtestVariables).buyFlag {
			(*backtestVariables).sellPrice = (*dt).PriceTf[index]
			(*backtestVariables).sellFlag = false
			(*backtestVariables).buyFlag = true
			(*backtestVariables).tpBuy = (*dt).PriceTf[index]
			financeCalculation(backtestVariables, "buy-tp")
			(*backtestVariables).totalTradesInOneDay++
		} else if (*backtestVariables).sellFlag {
			(*backtestVariables).sellPrice = (*dt).PriceTf[index]
			(*backtestVariables).sellFlag = false
			(*backtestVariables).totalTradesInOneDay++
		}
	}
}

func backtestHedging(dt *data.LayoutData, backtestVariables *backtestVariablesStruct, index int, mode string) {
	if mode == "open" {
		if (*dt).BuyFlag[index] && (*backtestVariables).buyFlag {
			if (*dt).IsFixTpsl {
				tpslCalculation(backtestVariables, (*dt).PriceTf[index], (*dt).MultiplyTp, (*dt).MultiplySl, (*dt).TpslFix, true)
			} else {
				tpslCalculation(backtestVariables, (*dt).PriceTf[index], (*dt).MultiplyTp, (*dt).MultiplySl, (*dt).TpslNonFix[index], true)
			}
			(*backtestVariables).buyFlag = false
		}
		if (*dt).SellFlag[index] && (*backtestVariables).sellFlag {
			if (*dt).IsFixTpsl {
				tpslCalculation(backtestVariables, (*dt).PriceTf[index], (*dt).MultiplyTp, (*dt).MultiplySl, (*dt).TpslFix, false)
			} else {
				tpslCalculation(backtestVariables, (*dt).PriceTf[index], (*dt).MultiplyTp, (*dt).MultiplySl, (*dt).TpslNonFix[index], false)
			}
			(*backtestVariables).sellFlag = false
		}
	} else if mode == "close" {
		if !(*backtestVariables).buyFlag {
			if checkTpslExit((*dt).Open[index], (*dt).High[index], (*dt).Low[index], (*dt).Close[index], (*backtestVariables).slBuy, false) {
				financeCalculation(backtestVariables, "buy-sl")
			} else if checkTpslExit((*dt).Open[index], (*dt).High[index], (*dt).Low[index], (*dt).Close[index], (*backtestVariables).tpBuy, true) {
				financeCalculation(backtestVariables, "buy-tp")
			}
		}
		if !(*backtestVariables).sellFlag {
			if checkTpslExit((*dt).Open[index], (*dt).High[index], (*dt).Low[index], (*dt).Close[index], (*backtestVariables).slSell, true) {
				financeCalculation(backtestVariables, "sell-sl")
			} else if checkTpslExit((*dt).Open[index], (*dt).High[index], (*dt).Low[index], (*dt).Close[index], (*backtestVariables).tpSell, false) {
				financeCalculation(backtestVariables, "sell-tp")
			}
		}
	} else {
		log.Panic("Backtest Hedging Mode Wrong")
	}
}

func backtestNetting(dt *data.LayoutData, backtestVariables *backtestVariablesStruct, index int, mode string) {
	if mode == "open" {
		if (*dt).BuyFlag[index] && (*backtestVariables).universalFlag && (*backtestVariables).buyFlag {
			if (*dt).IsFixTpsl {
				tpslCalculation(backtestVariables, (*dt).PriceTf[index], (*dt).MultiplyTp, (*dt).MultiplySl, (*dt).TpslFix, true)
			} else {
				tpslCalculation(backtestVariables, (*dt).PriceTf[index], (*dt).MultiplyTp, (*dt).MultiplySl, (*dt).TpslNonFix[index], true)
			}
			(*backtestVariables).buyFlag = false
			(*backtestVariables).universalFlag = false
		} else if (*dt).SellFlag[index] && (*backtestVariables).universalFlag && (*backtestVariables).sellFlag {
			if (*dt).IsFixTpsl {
				tpslCalculation(backtestVariables, (*dt).PriceTf[index], (*dt).MultiplyTp, (*dt).MultiplySl, (*dt).TpslFix, false)
			} else {
				tpslCalculation(backtestVariables, (*dt).PriceTf[index], (*dt).MultiplyTp, (*dt).MultiplySl, (*dt).TpslNonFix[index], false)
			}
			(*backtestVariables).sellFlag = false
			(*backtestVariables).universalFlag = false
		}
	} else if mode == "close" {
		if !(*backtestVariables).universalFlag && !(*backtestVariables).buyFlag {
			if checkTpslExit((*dt).Open[index], (*dt).High[index], (*dt).Low[index], (*dt).Close[index], (*backtestVariables).slBuy, false) {
				financeCalculation(backtestVariables, "buy-sl")
			} else if checkTpslExit((*dt).Open[index], (*dt).High[index], (*dt).Low[index], (*dt).Close[index], (*backtestVariables).tpBuy, true) {
				financeCalculation(backtestVariables, "buy-tp")
			}
		}
		if !(*backtestVariables).universalFlag && !(*backtestVariables).sellFlag {
			if checkTpslExit((*dt).Open[index], (*dt).High[index], (*dt).Low[index], (*dt).Close[index], (*backtestVariables).slSell, true) {
				financeCalculation(backtestVariables, "sell-sl")
			} else if checkTpslExit((*dt).Open[index], (*dt).High[index], (*dt).Low[index], (*dt).Close[index], (*backtestVariables).tpSell, false) {
				financeCalculation(backtestVariables, "sell-tp")
			}
		}
	} else {
		log.Panic("Backtest Netting Mode Wrong")
	}

}

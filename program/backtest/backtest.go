package backtest

import (
	"log"

	"github.com/SrJMaia/expert-advisor/program/check"
	"github.com/SrJMaia/expert-advisor/program/conversion"
	"github.com/SrJMaia/expert-advisor/program/data"
)

func financeCalculation(backtestMain *backtestVariables, leverage bool, isJpy bool, calculationMode string) {
	var comission float64
	var lot float64
	if leverage && (*backtestMain).capital > 3*(*backtestMain).initialCapital {
		lot = conversion.Round(((*backtestMain).capital - (*backtestMain).initialCapital), -3)
	} else {
		lot = (*backtestMain).initialCapital
	}
	comission = lot / 1000 * .07

	if calculationMode == "buy-tp" {
		(*backtestMain).buyResult = conversion.RoundIsJpy((lot*((*backtestMain).tpBuy-(*backtestMain).buyPrice))/(*backtestMain).tpBuy-comission, isJpy)
		(*backtestMain).capital = conversion.Round((*backtestMain).buyResult+(*backtestMain).capital, 2)
		check.MyCheckingFinanceCalculation((*backtestMain).capital, (*backtestMain).tpBuy, (*backtestMain).buyPrice, (*backtestMain).buyResult)
		updateBacktestHeart(backtestMain, "buy")
	} else if calculationMode == "buy-sl" {
		(*backtestMain).buyResult = conversion.RoundIsJpy((lot*((*backtestMain).slBuy-(*backtestMain).buyPrice))/(*backtestMain).slBuy-comission, isJpy)
		(*backtestMain).capital = conversion.Round((*backtestMain).buyResult+(*backtestMain).capital, 2)
		check.MyCheckingFinanceCalculation((*backtestMain).capital, (*backtestMain).slBuy, (*backtestMain).buyPrice, (*backtestMain).buyResult)
		updateBacktestHeart(backtestMain, "buy")
	} else if calculationMode == "sell-tp" {
		(*backtestMain).sellResult = conversion.RoundIsJpy((lot*((*backtestMain).sellPrice-(*backtestMain).tpSell))/(*backtestMain).tpSell-comission, isJpy)
		(*backtestMain).capital = conversion.Round((*backtestMain).sellResult+(*backtestMain).capital, 2)
		check.MyCheckingFinanceCalculation((*backtestMain).capital, (*backtestMain).tpSell, (*backtestMain).sellPrice, (*backtestMain).sellResult)
		updateBacktestHeart(backtestMain, "sell")
	} else if calculationMode == "sell-sl" {
		(*backtestMain).sellResult = conversion.RoundIsJpy((lot*((*backtestMain).sellPrice-(*backtestMain).slSell))/(*backtestMain).slSell-comission, isJpy)
		(*backtestMain).capital = conversion.Round((*backtestMain).sellResult+(*backtestMain).capital, 2)
		check.MyCheckingFinanceCalculation((*backtestMain).capital, (*backtestMain).slSell, (*backtestMain).sellPrice, (*backtestMain).sellResult)
		updateBacktestHeart(backtestMain, "sell")
	} else {
		log.Panic("Finance Calculation Wrong Mode")
	}

}

func tpslCalculation(backtestMain *backtestVariables, price float64, multiplyTp float64, multiplySl float64, tpSlValue float64, isJpy bool, buy bool) {
	if buy {
		(*backtestMain).tpBuy = conversion.RoundIsJpy(price+(multiplyTp*tpSlValue), isJpy)
		(*backtestMain).slBuy = conversion.RoundIsJpy(price-(multiplySl*tpSlValue), isJpy)
	} else {
		(*backtestMain).tpSell = conversion.RoundIsJpy(price-(multiplyTp*tpSlValue), isJpy)
		(*backtestMain).slSell = conversion.RoundIsJpy(price+(multiplySl*tpSlValue), isJpy)
	}
}

func checkTpslExit(open float64, high float64, low float64, close float64, tpsl float64, higher bool) bool {
	if higher {
		return open >= tpsl || high >= tpsl || low >= tpsl || close >= tpsl
	}
	return open <= tpsl || high <= tpsl || low <= tpsl || close <= tpsl
}

func takeProfitFimathe(channel float64, multiply float64, channelInPips float64, isJpy bool, buy bool) float64 {
	if buy {
		return conversion.RoundIsJpy(channel+(multiply*channelInPips), isJpy)
	}
	return conversion.RoundIsJpy(channel-(multiply*channelInPips), isJpy)
}

func NoTpslBacktest(dt *data.LayoutData, balance float64, jpy bool, leverage bool) ([]float64, []float64, []float64) {

	var backtestMain = backtestHeart(balance, dt.SizeBuy, dt.SizeSell)

	for i := 1; i < len((*dt).Open); i++ {

		if (*dt).PriceTf[i] != (*dt).PriceTf[i-1] {
			if (*dt).BuyFlag[i] {
				if !backtestMain.sellFlag {
					backtestMain.buyPrice = (*dt).PriceTf[i]
					backtestMain.buyFlag = false
					backtestMain.sellFlag = true
					backtestMain.tpSell = (*dt).PriceTf[i]
					financeCalculation(&backtestMain, leverage, jpy, "sell-tp")
					backtestMain.totalTrades[backtestMain.iTotalTrades] = backtestMain.capital
					backtestMain.sellTrades[backtestMain.iSell] = conversion.Round(backtestMain.sellTrades[backtestMain.iSell-1]+backtestMain.sellResult, 2)
					backtestMain.iTotalTrades++
					backtestMain.iSell++
				} else {
					backtestMain.buyPrice = (*dt).PriceTf[i]
					backtestMain.buyFlag = false
				}
			}
			if (*dt).SellFlag[i] {
				if !backtestMain.buyFlag {
					backtestMain.sellPrice = (*dt).PriceTf[i]
					backtestMain.sellFlag = false
					backtestMain.buyFlag = true
					backtestMain.tpBuy = (*dt).PriceTf[i]
					financeCalculation(&backtestMain, leverage, jpy, "buy-tp")
					backtestMain.totalTrades[backtestMain.iTotalTrades] = backtestMain.capital
					backtestMain.buyTrades[backtestMain.iBuy] = conversion.Round(backtestMain.buyTrades[backtestMain.iBuy-1]+backtestMain.buyResult, 2)
					backtestMain.iTotalTrades++
					backtestMain.iBuy++
				} else {
					backtestMain.sellPrice = (*dt).PriceTf[i]
					backtestMain.sellFlag = false
				}
			}
		}
	}

	var totalTrades = conversion.RemoveExcessZeros(backtestMain.totalTrades)
	var buyTrades = conversion.RemoveExcessZeros(backtestMain.buyTrades)
	var sellTrades = conversion.RemoveExcessZeros(backtestMain.sellTrades)

	return totalTrades, buyTrades, sellTrades

}

func HedgingBacktest(dt *data.LayoutData, multiplyTp float64, multiplySl float64, balance float64, jpy bool, leverage bool) ([]float64, []float64, []float64) {

	var backtestMain backtestVariables = backtestHeart(balance, dt.SizeBuy, dt.SizeSell)
	// var dayTeste = true

	for i := 1; i < len((*dt).Open); i++ {

		// if (*dt).TimeWeekDays[i] != (*dt).TimeWeekDays[i-1] && !dayTeste {
		// 	dayTeste = true
		// }
		//if (*dt).TimeHour[i] > 6 && (*dt).TimeHour[i] < 20 {

		if (*dt).PriceTf[i] != (*dt).PriceTf[i-1] {
			if (*dt).BuyFlag[i] && backtestMain.buyFlag {
				backtestMain.buyPrice = (*dt).PriceTf[i]
				if (*dt).IsFixTpsl {
					tpslCalculation(&backtestMain, (*dt).PriceTf[i], multiplyTp, multiplySl, (*dt).TpslFix, jpy, true)
				} else {
					tpslCalculation(&backtestMain, (*dt).PriceTf[i], multiplyTp, multiplySl, (*dt).TpslNonFix[i], jpy, true)
				}
				backtestMain.buyFlag = false
			}
			if (*dt).SellFlag[i] && backtestMain.sellFlag {
				backtestMain.sellPrice = (*dt).PriceTf[i]
				if (*dt).IsFixTpsl {
					tpslCalculation(&backtestMain, (*dt).PriceTf[i], multiplyTp, multiplySl, (*dt).TpslFix, jpy, false)
				} else {
					tpslCalculation(&backtestMain, (*dt).PriceTf[i], multiplyTp, multiplySl, (*dt).TpslNonFix[i], jpy, false)
				}
				backtestMain.sellFlag = false
			}
		} else if (*dt).PriceTf[i] == (*dt).PriceTf[i-1] {
			if !backtestMain.buyFlag {
				if checkTpslExit((*dt).Open[i], (*dt).High[i], (*dt).Low[i], (*dt).Close[i], backtestMain.slBuy, false) {
					financeCalculation(&backtestMain, leverage, jpy, "buy-sl")
				} else if checkTpslExit((*dt).Open[i], (*dt).High[i], (*dt).Low[i], (*dt).Close[i], backtestMain.tpBuy, true) {
					financeCalculation(&backtestMain, leverage, jpy, "buy-tp")
				}
			}
			if !backtestMain.sellFlag {
				if checkTpslExit((*dt).Open[i], (*dt).High[i], (*dt).Low[i], (*dt).Close[i], backtestMain.slSell, true) {
					financeCalculation(&backtestMain, leverage, jpy, "sell-sl")
				} else if checkTpslExit((*dt).Open[i], (*dt).High[i], (*dt).Low[i], (*dt).Close[i], backtestMain.tpSell, false) {
					financeCalculation(&backtestMain, leverage, jpy, "sell-tp")
				}
			}
		}
		//}
		// if (*dt).TimeHour[i] >= 20 && (*dt).TimeHour[i] <= 23 && (*dt).TimeWeekDays[i] == 5 {
		// 	if !backtestMain.buyFlag {
		//		ANTES DE ENVIAR PARA FINANCE CALCULATION
		// 		TROCAR O TP DE TUDO E FAZER COMO NO NOTPSLBACKTEST
		// 		backtestMain.capital, backtestMain.buyResult = financeCalculation(backtestMain.capital, (*dt).Open[i], backtestMain.buyPrice, backtestMain.buyPrice, balance, leverage, jpy)
		// 		updateBacktestHeart(&backtestMain, "buy")
		// 	}
		// 	if !backtestMain.sellFlag {
		// 		backtestMain.capital, backtestMain.sellResult = financeCalculation(backtestMain.capital, backtestMain.sellPrice, (*dt).Open[i], backtestMain.sellPrice, balance, leverage, jpy)
		// 		updateBacktestHeart(&backtestMain, "sell")
		// 	}
		// 	dayTeste = false
		// }
	}

	var totalTrades []float64 = conversion.RemoveExcessZeros(backtestMain.totalTrades)
	var buyTrades []float64 = conversion.RemoveExcessZeros(backtestMain.buyTrades)
	var sellTrades []float64 = conversion.RemoveExcessZeros(backtestMain.sellTrades)

	return totalTrades, buyTrades, sellTrades

}

func NettingBacktest(dt *data.LayoutData, multiplyTp float64, multiplySl float64, balance float64, jpy bool, leverage bool) ([]float64, []float64, []float64) {

	var backtestMain backtestVariables = backtestHeart(balance, dt.SizeBuy, dt.SizeSell)

	for i := 1; i < len((*dt).Open); i++ {

		if (*dt).PriceTf[i] != (*dt).PriceTf[i-1] {
			if (*dt).BuyFlag[i] && backtestMain.universalFlag && backtestMain.buyFlag {
				backtestMain.buyPrice = (*dt).PriceTf[i]
				if (*dt).IsFixTpsl {
					tpslCalculation(&backtestMain, (*dt).PriceTf[i], multiplyTp, multiplySl, (*dt).TpslFix, jpy, true)
				} else {
					tpslCalculation(&backtestMain, (*dt).PriceTf[i], multiplyTp, multiplySl, (*dt).TpslNonFix[i], jpy, true)
				}
				backtestMain.buyFlag = false
				backtestMain.universalFlag = false
			} else if (*dt).SellFlag[i] && backtestMain.universalFlag && backtestMain.sellFlag {
				backtestMain.sellPrice = (*dt).PriceTf[i]
				if (*dt).IsFixTpsl {
					tpslCalculation(&backtestMain, (*dt).PriceTf[i], multiplyTp, multiplySl, (*dt).TpslFix, jpy, false)
				} else {
					tpslCalculation(&backtestMain, (*dt).PriceTf[i], multiplyTp, multiplySl, (*dt).TpslNonFix[i], jpy, false)
				}
				backtestMain.sellFlag = false
				backtestMain.universalFlag = false
			}
		} else if (*dt).PriceTf[i] == (*dt).PriceTf[i-1] {
			if !backtestMain.universalFlag && !backtestMain.buyFlag {
				if checkTpslExit((*dt).Open[i], (*dt).High[i], (*dt).Low[i], (*dt).Close[i], backtestMain.slBuy, false) {
					financeCalculation(&backtestMain, leverage, jpy, "buy-sl")
				} else if checkTpslExit((*dt).Open[i], (*dt).High[i], (*dt).Low[i], (*dt).Close[i], backtestMain.tpBuy, true) {
					financeCalculation(&backtestMain, leverage, jpy, "buy-tp")
				}
			}
			if !backtestMain.universalFlag && !backtestMain.sellFlag {
				if checkTpslExit((*dt).Open[i], (*dt).High[i], (*dt).Low[i], (*dt).Close[i], backtestMain.slSell, true) {
					financeCalculation(&backtestMain, leverage, jpy, "sell-sl")
				} else if checkTpslExit((*dt).Open[i], (*dt).High[i], (*dt).Low[i], (*dt).Close[i], backtestMain.tpSell, false) {
					financeCalculation(&backtestMain, leverage, jpy, "sell-tp")
				}
			}
		}
	}

	var totalTrades []float64 = conversion.RemoveExcessZeros(backtestMain.totalTrades)
	var buyTrades []float64 = conversion.RemoveExcessZeros(backtestMain.buyTrades)
	var sellTrades []float64 = conversion.RemoveExcessZeros(backtestMain.sellTrades)

	return totalTrades, buyTrades, sellTrades

}

func Fimathe(dt *data.LayoutData, channelInPips float64, startPointInPips float64, multiplyTpChannel float64, balance float64, jpy bool, leverage bool) ([]float64, []float64, []float64) {
	// Canal =~ 10 a 200 pips ou 0.0010 a 0.0200
	// StartPoint range(+halfChannel até -halfChannel)
	// multiplyTpChannel = range(1,6)
	// 2 Tipos, tp sl fixos ou trailing stop

	var halfChanneInlPips = channelInPips / 2.
	startPointInPips += (*dt).Open[0]

	var upperChannel = conversion.RoundIsJpy(startPointInPips+halfChanneInlPips, jpy)
	var lowerChannel = conversion.RoundIsJpy(startPointInPips-halfChanneInlPips, jpy)

	var backtestMain backtestVariables = backtestHeart(balance, 100000, 100000)

	for i := 1; i < len((*dt).Open); i++ {

		if (*dt).PriceTf[i] != (*dt).PriceTf[i-1] {
			if (*dt).PriceTf[i] > upperChannel && backtestMain.buyFlag && backtestMain.universalFlag && (*dt).BuyFlag[i] {
				backtestMain.buyPrice = (*dt).PriceTf[i]
				backtestMain.tpBuy = takeProfitFimathe(upperChannel, multiplyTpChannel, channelInPips, jpy, true)
				backtestMain.slBuy = conversion.RoundIsJpy(lowerChannel-channelInPips, jpy)
				backtestMain.buyFlag = false
				backtestMain.universalFlag = false
			} else if (*dt).PriceTf[i] < lowerChannel && backtestMain.sellFlag && backtestMain.universalFlag && (*dt).SellFlag[i] {
				backtestMain.sellPrice = (*dt).PriceTf[i]
				backtestMain.tpSell = takeProfitFimathe(lowerChannel, multiplyTpChannel, channelInPips, jpy, false)
				backtestMain.slSell = conversion.RoundIsJpy(upperChannel+channelInPips, jpy)
				backtestMain.sellFlag = false
				backtestMain.universalFlag = false
			}
			if (*dt).PriceTf[i] > upperChannel {
				lowerChannel = conversion.RoundIsJpy(upperChannel, jpy)
				for (*dt).PriceTf[i] > upperChannel {
					upperChannel = conversion.RoundIsJpy(upperChannel+channelInPips, jpy)
				}
			} else if (*dt).PriceTf[i] < lowerChannel {
				upperChannel = conversion.RoundIsJpy(lowerChannel, jpy)
				for (*dt).PriceTf[i] < lowerChannel {
					lowerChannel = conversion.RoundIsJpy(lowerChannel-channelInPips, jpy)
				}
			}
		} else if (*dt).PriceTf[i] == (*dt).PriceTf[i-1] {
			if !backtestMain.universalFlag && !backtestMain.buyFlag {
				if checkTpslExit((*dt).Open[i], (*dt).High[i], (*dt).Low[i], (*dt).Close[i], backtestMain.slBuy, false) {
					financeCalculation(&backtestMain, leverage, jpy, "buy-sl")
				} else if checkTpslExit((*dt).Open[i], (*dt).High[i], (*dt).Low[i], (*dt).Close[i], backtestMain.tpBuy, true) {
					financeCalculation(&backtestMain, leverage, jpy, "buy-tp")
				}
			}
			if !backtestMain.universalFlag && !backtestMain.sellFlag {
				if checkTpslExit((*dt).Open[i], (*dt).High[i], (*dt).Low[i], (*dt).Close[i], backtestMain.slSell, true) {
					financeCalculation(&backtestMain, leverage, jpy, "sell-sl")
				} else if checkTpslExit((*dt).Open[i], (*dt).High[i], (*dt).Low[i], (*dt).Close[i], backtestMain.tpSell, false) {
					financeCalculation(&backtestMain, leverage, jpy, "sell-tp")
				}
			}
		}

	}

	var totalTrades []float64 = conversion.RemoveExcessZeros(backtestMain.totalTrades)
	var buyTrades []float64 = conversion.RemoveExcessZeros(backtestMain.buyTrades)
	var sellTrades []float64 = conversion.RemoveExcessZeros(backtestMain.sellTrades)

	return totalTrades, buyTrades, sellTrades
}

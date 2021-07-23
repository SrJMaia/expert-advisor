package backtest

import (
	"log"
	"math"

	"github.com/SrJMaia/expert-advisor/program/conversion"
	"github.com/SrJMaia/expert-advisor/program/data"
)

type backtestVariables struct {
	buyPrice      float64
	sellPrice     float64
	tpSell        float64
	slSell        float64
	tpBuy         float64
	slBuy         float64
	buyResult     float64
	sellResult    float64
	capital       float64
	buyFlag       bool
	sellFlag      bool
	universalFlag bool
	updateBuy     bool
	updateSell    bool
	iBuy          uint32
	iSell         uint32
	iTotalTrades  uint32
	totalTrades   []float64
	buyTrades     []float64
	sellTrades    []float64
}

func backtestHeart(balance float64, sizeBuy uint32, sizeSell uint32) backtestVariables {
	var backtestMain = backtestVariables{
		buyPrice:      0.,
		sellPrice:     0.,
		tpSell:        0.,
		slSell:        0.,
		tpBuy:         0.,
		slBuy:         0.,
		buyResult:     0.,
		sellResult:    0.,
		capital:       balance,
		buyFlag:       true,
		sellFlag:      true,
		universalFlag: true,
		iBuy:          1,
		iSell:         1,
		iTotalTrades:  1,
		totalTrades:   make([]float64, sizeBuy+sizeSell),
		buyTrades:     make([]float64, sizeBuy),
		sellTrades:    make([]float64, sizeSell),
	}

	backtestMain.totalTrades[0] = balance
	backtestMain.buyTrades[0] = balance
	backtestMain.sellTrades[0] = balance

	return backtestMain
}

func financeCalculation(balance float64, initialPrice float64, finalPrice float64, eurPrice float64, initialBalance float64, leverage bool) (float64, float64) {
	var lot float64
	var comission float64
	var result float64
	var total float64
	if leverage && balance > 3*initialBalance {
		lot = conversion.Round((balance - initialBalance), -3)
	} else {
		lot = initialBalance
	}
	comission = lot / 1000 * .07

	result = conversion.Round((lot*(initialPrice-finalPrice))/eurPrice-comission, 5)
	total = conversion.Round(result+balance, 2)

	if math.IsNaN(total) || initialPrice == 0 || finalPrice == 0 || eurPrice == 0 {
		log.Panic("Finance Calculation Error. NaN")
	}
	return total, result

}

func tpslCalculationFix(price float64, multiplyTpF float64, multiplySlF float64, jpyRound bool, buy bool) (float64, float64) {
	var tp float64
	var sl float64
	var roundTPSL int
	if jpyRound {
		roundTPSL = 3
	} else {
		roundTPSL = 5
	}
	if buy {
		tp = conversion.Round(price+multiplyTpF, roundTPSL)
		sl = conversion.Round(price-multiplySlF, roundTPSL)
	} else {
		tp = conversion.Round(price-multiplyTpF, roundTPSL)
		sl = conversion.Round(price+multiplySlF, roundTPSL)
	}
	return tp, sl
}

func tpslCalculationNonFix(price float64, multiplyTpF float64, multiplySlF float64, tpSlValue float64, jpyRound bool, buy bool) (float64, float64) {
	// Example: ATR
	var tp float64
	var sl float64
	var roundTPSL int
	if jpyRound {
		roundTPSL = 3
	} else {
		roundTPSL = 5
	}
	if buy {
		tp = conversion.Round(price+(multiplyTpF*tpSlValue), roundTPSL)
		sl = conversion.Round(price-(multiplySlF*tpSlValue), roundTPSL)
	} else {
		tp = conversion.Round(price-(multiplyTpF*tpSlValue), roundTPSL)
		sl = conversion.Round(price+(multiplySlF*tpSlValue), roundTPSL)
	}
	return tp, sl
}

func checkTpslExit(open float64, high float64, low float64, close float64, tpsl float64, higher bool) bool {
	if higher {
		return open >= tpsl || high >= tpsl || low >= tpsl || close >= tpsl
	}
	return open <= tpsl || high <= tpsl || low <= tpsl || close <= tpsl
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

func NoTpslBacktest(dt *data.LayoutData, balance float64, jpy bool, leverage bool) ([]float64, []float64, []float64) {

	var backtestMain = backtestHeart(balance, dt.SizeBuy, dt.SizeSell)

	for i := 1; i < len((*dt).Open); i++ {

		if (*dt).PriceTf[i] != (*dt).PriceTf[i-1] {
			if (*dt).BuyFlag[i] {
				if !backtestMain.sellFlag {
					backtestMain.buyPrice = (*dt).PriceTf[i]
					backtestMain.buyFlag = false
					backtestMain.sellFlag = true
					backtestMain.capital, backtestMain.sellResult = financeCalculation(backtestMain.capital, backtestMain.sellPrice, (*dt).PriceTf[i], (*dt).PriceTf[i], balance, leverage)
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
					backtestMain.capital, backtestMain.buyResult = financeCalculation(backtestMain.capital, (*dt).PriceTf[i], backtestMain.buyPrice, (*dt).PriceTf[i], balance, leverage)
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

func HedgingBacktest(dt *data.LayoutData, multiplyTpB float64, multiplySlB float64, balance float64, jpy bool, leverage bool) ([]float64, []float64, []float64) {

	var backtestMain backtestVariables = backtestHeart(balance, dt.SizeBuy, dt.SizeSell)
	var dayTeste = true

	for i := 1; i < len((*dt).Open); i++ {

		if (*dt).TimeWeekDays[i] != (*dt).TimeWeekDays[i-1] && !dayTeste {
			dayTeste = true
		}
		//if (*dt).TimeHour[i] > 6 && (*dt).TimeHour[i] < 20 {

		if (*dt).PriceTf[i] != (*dt).PriceTf[i-1] {
			if (*dt).BuyFlag[i] && backtestMain.buyFlag {
				backtestMain.buyPrice = (*dt).PriceTf[i]
				backtestMain.tpBuy, backtestMain.slBuy = tpslCalculationNonFix((*dt).PriceTf[i], multiplyTpB, multiplySlB, (*dt).Tpsl[i], jpy, true)
				backtestMain.buyFlag = false
			}
			if (*dt).SellFlag[i] && backtestMain.sellFlag {
				backtestMain.sellPrice = (*dt).PriceTf[i]
				backtestMain.tpSell, backtestMain.slSell = tpslCalculationNonFix((*dt).PriceTf[i], multiplyTpB, multiplySlB, (*dt).Tpsl[i], jpy, false)
				backtestMain.sellFlag = false
			}
		} else if (*dt).PriceTf[i] == (*dt).PriceTf[i-1] {
			if !backtestMain.buyFlag {
				if checkTpslExit((*dt).Open[i], (*dt).High[i], (*dt).Low[i], (*dt).Close[i], backtestMain.slBuy, false) {
					backtestMain.capital, backtestMain.buyResult = financeCalculation(backtestMain.capital, backtestMain.slBuy, backtestMain.buyPrice, backtestMain.slBuy, balance, leverage)
					updateBacktestHeart(&backtestMain, "buy")
				} else if checkTpslExit((*dt).Open[i], (*dt).High[i], (*dt).Low[i], (*dt).Close[i], backtestMain.tpBuy, true) {
					backtestMain.capital, backtestMain.buyResult = financeCalculation(backtestMain.capital, backtestMain.tpBuy, backtestMain.buyPrice, backtestMain.tpBuy, balance, leverage)
					updateBacktestHeart(&backtestMain, "buy")
				}
			}
			if !backtestMain.sellFlag {
				if checkTpslExit((*dt).Open[i], (*dt).High[i], (*dt).Low[i], (*dt).Close[i], backtestMain.slSell, true) {
					// Istead of returning, i could update there
					backtestMain.capital, backtestMain.sellResult = financeCalculation(backtestMain.capital, backtestMain.sellPrice, backtestMain.slSell, backtestMain.sellPrice, balance, leverage)
					updateBacktestHeart(&backtestMain, "sell")
				} else if checkTpslExit((*dt).Open[i], (*dt).High[i], (*dt).Low[i], (*dt).Close[i], backtestMain.tpSell, false) {
					backtestMain.capital, backtestMain.sellResult = financeCalculation(backtestMain.capital, backtestMain.sellPrice, backtestMain.tpSell, backtestMain.sellPrice, balance, leverage)
					updateBacktestHeart(&backtestMain, "sell")
				}
			}
		}
		//}
		if (*dt).TimeHour[i] >= 20 && (*dt).TimeHour[i] <= 23 && (*dt).TimeWeekDays[i] == 5 {
			if !backtestMain.buyFlag {
				backtestMain.capital, backtestMain.buyResult = financeCalculation(backtestMain.capital, (*dt).Open[i], backtestMain.buyPrice, backtestMain.buyPrice, balance, leverage)
				updateBacktestHeart(&backtestMain, "buy")
			}
			if !backtestMain.sellFlag {
				backtestMain.capital, backtestMain.sellResult = financeCalculation(backtestMain.capital, backtestMain.sellPrice, (*dt).Open[i], backtestMain.sellPrice, balance, leverage)
				updateBacktestHeart(&backtestMain, "sell")
			}
			dayTeste = false
		}
	}

	var totalTrades []float64 = conversion.RemoveExcessZeros(backtestMain.totalTrades)
	var buyTrades []float64 = conversion.RemoveExcessZeros(backtestMain.buyTrades)
	var sellTrades []float64 = conversion.RemoveExcessZeros(backtestMain.sellTrades)

	return totalTrades, buyTrades, sellTrades

}

func NettingBacktest(dt *data.LayoutData, multiplyTpB float64, multiplySlB float64, balance float64, jpy bool, leverage bool) ([]float64, []float64, []float64) {

	var backtestMain backtestVariables = backtestHeart(balance, dt.SizeBuy, dt.SizeSell)

	for i := 1; i < len((*dt).Open); i++ {

		if (*dt).PriceTf[i] != (*dt).PriceTf[i-1] {
			if (*dt).BuyFlag[i] && backtestMain.universalFlag && backtestMain.buyFlag {
				backtestMain.buyPrice = (*dt).PriceTf[i]
				backtestMain.tpBuy, backtestMain.slBuy = tpslCalculationNonFix((*dt).PriceTf[i], multiplyTpB, multiplySlB, (*dt).Tpsl[i], jpy, true)
				backtestMain.buyFlag = false
				backtestMain.universalFlag = false
			} else if (*dt).SellFlag[i] && backtestMain.universalFlag && backtestMain.sellFlag {
				backtestMain.sellPrice = (*dt).PriceTf[i]
				backtestMain.tpSell, backtestMain.slSell = tpslCalculationNonFix((*dt).PriceTf[i], multiplyTpB, multiplySlB, (*dt).Tpsl[i], jpy, false)
				backtestMain.sellFlag = false
				backtestMain.universalFlag = false
			}
		} else if (*dt).PriceTf[i] == (*dt).PriceTf[i-1] {
			if !backtestMain.universalFlag && !backtestMain.buyFlag {
				if checkTpslExit((*dt).Open[i], (*dt).High[i], (*dt).Low[i], (*dt).Close[i], backtestMain.slBuy, false) {
					backtestMain.capital, backtestMain.buyResult = financeCalculation(backtestMain.capital, backtestMain.slBuy, backtestMain.buyPrice, backtestMain.slBuy, balance, leverage)
					updateBacktestHeart(&backtestMain, "buy")
				} else if checkTpslExit((*dt).Open[i], (*dt).High[i], (*dt).Low[i], (*dt).Close[i], backtestMain.tpBuy, true) {
					backtestMain.capital, backtestMain.buyResult = financeCalculation(backtestMain.capital, backtestMain.tpBuy, backtestMain.buyPrice, backtestMain.tpBuy, balance, leverage)
					updateBacktestHeart(&backtestMain, "buy")
				}
			}
			if !backtestMain.universalFlag && !backtestMain.sellFlag {
				if checkTpslExit((*dt).Open[i], (*dt).High[i], (*dt).Low[i], (*dt).Close[i], backtestMain.slSell, true) {
					backtestMain.capital, backtestMain.sellResult = financeCalculation(backtestMain.capital, backtestMain.sellPrice, backtestMain.slSell, backtestMain.sellPrice, balance, leverage)
					updateBacktestHeart(&backtestMain, "sell")
				} else if checkTpslExit((*dt).Open[i], (*dt).High[i], (*dt).Low[i], (*dt).Close[i], backtestMain.tpSell, false) {
					backtestMain.capital, backtestMain.sellResult = financeCalculation(backtestMain.capital, backtestMain.sellPrice, backtestMain.tpSell, backtestMain.sellPrice, balance, leverage)
					updateBacktestHeart(&backtestMain, "sell")
				}
			}
		}
	}

	var totalTrades []float64 = conversion.RemoveExcessZeros(backtestMain.totalTrades)
	var buyTrades []float64 = conversion.RemoveExcessZeros(backtestMain.buyTrades)
	var sellTrades []float64 = conversion.RemoveExcessZeros(backtestMain.sellTrades)

	return totalTrades, buyTrades, sellTrades

}

func takeProfitFimathe(channel float64, multiply float64, channelInPips float64, jpy bool, buy bool) float64 {
	var roundValue int
	if jpy {
		roundValue = 3
	} else {
		roundValue = 5
	}
	if buy {
		return conversion.Round(channel+(multiply*channelInPips), roundValue)
	}
	return conversion.Round(channel-(multiply*channelInPips), roundValue)
}

func Fimathe(dt *data.LayoutData, channelInPips float64, startPointInPips float64, multiplyTpChannel float64, balance float64, jpy bool, leverage bool) ([]float64, []float64, []float64) {
	// Canal =~ 10 a 200 pips ou 0.0010 a 0.0200
	// StartPoint range(+halfChannel até -halfChannel)
	// multiplyTpChannel = range(1,6)
	// 2 Tipos, tp sl fixos ou trailing stop
	var dayTeste = true
	var daysCounter int
	var halfChanneInlPips = channelInPips / 2.
	startPointInPips += (*dt).Open[0]
	var roundPlace int
	if jpy {
		roundPlace = 3
	} else {
		roundPlace = 5
	}

	var upperChannel = conversion.Round(startPointInPips+halfChanneInlPips, roundPlace)
	var lowerChannel = conversion.Round(startPointInPips-halfChanneInlPips, roundPlace)

	var backtestMain backtestVariables = backtestHeart(balance, 100000, 100000)

	for i := 1; i < len((*dt).Open); i++ {

		if (*dt).TimeWeekDays[i] != (*dt).TimeWeekDays[i-1] && !dayTeste {
			dayTeste = true
		}

		if (*dt).TimeWeekDays[i] != (*dt).TimeWeekDays[i-1] {
			daysCounter = 0
		}
		if (*dt).TimeHour[i] > 8 && (*dt).TimeHour[i] < 18 {

			if (*dt).PriceTf[i] != (*dt).PriceTf[i-1] && daysCounter < 1 {
				if (*dt).PriceTf[i] > upperChannel && backtestMain.buyFlag && backtestMain.universalFlag && (*dt).BuyFlag[i] {
					backtestMain.buyPrice = (*dt).PriceTf[i]
					backtestMain.tpBuy = takeProfitFimathe(upperChannel, multiplyTpChannel, channelInPips, jpy, true)
					backtestMain.slBuy = conversion.Round(lowerChannel-channelInPips, roundPlace)
					backtestMain.buyFlag = false
					backtestMain.universalFlag = false
					daysCounter++
				} else if (*dt).PriceTf[i] < lowerChannel && backtestMain.sellFlag && backtestMain.universalFlag && (*dt).SellFlag[i] {
					backtestMain.sellPrice = (*dt).PriceTf[i]
					backtestMain.tpSell = takeProfitFimathe(lowerChannel, multiplyTpChannel, channelInPips, jpy, false)
					backtestMain.slSell = conversion.Round(upperChannel+channelInPips, roundPlace)
					backtestMain.sellFlag = false
					backtestMain.universalFlag = false
					daysCounter++
				}
				if (*dt).PriceTf[i] > upperChannel {
					lowerChannel = conversion.Round(upperChannel, roundPlace)
					for (*dt).PriceTf[i] > upperChannel {
						upperChannel = conversion.Round(upperChannel+channelInPips, roundPlace)
					}
				} else if (*dt).PriceTf[i] < lowerChannel {
					upperChannel = conversion.Round(lowerChannel, roundPlace)
					for (*dt).PriceTf[i] < lowerChannel {
						lowerChannel = conversion.Round(lowerChannel-channelInPips, roundPlace)
					}
				}

			} else if (*dt).PriceTf[i] == (*dt).PriceTf[i-1] {
				if !backtestMain.universalFlag && !backtestMain.buyFlag {
					if checkTpslExit((*dt).Open[i], (*dt).High[i], (*dt).Low[i], (*dt).Close[i], backtestMain.slBuy, false) {
						backtestMain.capital, backtestMain.buyResult = financeCalculation(backtestMain.capital, backtestMain.slBuy, backtestMain.buyPrice, backtestMain.slBuy, balance, leverage)
						updateBacktestHeart(&backtestMain, "buy")
					} else if checkTpslExit((*dt).Open[i], (*dt).High[i], (*dt).Low[i], (*dt).Close[i], backtestMain.tpBuy, true) {
						backtestMain.capital, backtestMain.buyResult = financeCalculation(backtestMain.capital, backtestMain.tpBuy, backtestMain.buyPrice, backtestMain.tpBuy, balance, leverage)
						updateBacktestHeart(&backtestMain, "buy")
					}
				}
				if !backtestMain.universalFlag && !backtestMain.sellFlag {
					if checkTpslExit((*dt).Open[i], (*dt).High[i], (*dt).Low[i], (*dt).Close[i], backtestMain.slSell, true) {
						backtestMain.capital, backtestMain.sellResult = financeCalculation(backtestMain.capital, backtestMain.sellPrice, backtestMain.slSell, backtestMain.sellPrice, balance, leverage)
						updateBacktestHeart(&backtestMain, "sell")
					} else if checkTpslExit((*dt).Open[i], (*dt).High[i], (*dt).Low[i], (*dt).Close[i], backtestMain.tpSell, false) {
						backtestMain.capital, backtestMain.sellResult = financeCalculation(backtestMain.capital, backtestMain.sellPrice, backtestMain.tpSell, backtestMain.sellPrice, balance, leverage)
						updateBacktestHeart(&backtestMain, "sell")
					}
				}
			}
		}
		if (*dt).TimeHour[i] >= 20 && (*dt).TimeHour[i] <= 23 && (*dt).TimeWeekDays[i] == 5 {
			if !backtestMain.buyFlag {
				backtestMain.capital, backtestMain.buyResult = financeCalculation(backtestMain.capital, (*dt).Open[i], backtestMain.buyPrice, backtestMain.buyPrice, balance, leverage)
				updateBacktestHeart(&backtestMain, "buy")
			}
			if !backtestMain.sellFlag {
				backtestMain.capital, backtestMain.sellResult = financeCalculation(backtestMain.capital, backtestMain.sellPrice, (*dt).Open[i], backtestMain.sellPrice, balance, leverage)
				updateBacktestHeart(&backtestMain, "sell")
			}
			dayTeste = false
		}
	}

	var totalTrades []float64 = conversion.RemoveExcessZeros(backtestMain.totalTrades)
	var buyTrades []float64 = conversion.RemoveExcessZeros(backtestMain.buyTrades)
	var sellTrades []float64 = conversion.RemoveExcessZeros(backtestMain.sellTrades)

	return totalTrades, buyTrades, sellTrades
}

func FimatheCosturada(dt *data.LayoutData, channelInPips float64, startPointInPips float64, multiplyTpChannel float64, balance float64, jpy bool, leverage bool) ([]float64, []float64, []float64) {
	// Canal =~ 10 a 200 pips ou 0.0010 a 0.0200
	// StartPoint range(+halfChannel até -halfChannel)
	// multiplyTpChannel = range(1,6)
	// 2 Tipos, tp sl fixos ou trailing stop
	var dayTeste = true
	var halfChanneInlPips = channelInPips / 2.
	startPointInPips += (*dt).Open[0]
	var roundPlace int
	if jpy {
		roundPlace = 3
	} else {
		roundPlace = 5
	}

	var upperChannel = conversion.Round(startPointInPips+halfChanneInlPips, roundPlace)
	var lowerChannel = conversion.Round(startPointInPips-halfChanneInlPips, roundPlace)

	var backtestMain backtestVariables = backtestHeart(balance, 100000, 100000)

	for i := 1; i < len((*dt).Open); i++ {

		if (*dt).TimeWeekDays[i] != (*dt).TimeWeekDays[i-1] && !dayTeste {
			dayTeste = true
		}
		if (*dt).TimeHour[i] > 8 && (*dt).TimeHour[i] < 18 {

			if (*dt).PriceTf[i] != (*dt).PriceTf[i-1] {
				if (*dt).PriceTf[i] > upperChannel && backtestMain.buyFlag && backtestMain.universalFlag {
					backtestMain.buyPrice = (*dt).PriceTf[i]
					backtestMain.tpBuy = takeProfitFimathe(upperChannel, multiplyTpChannel, channelInPips, jpy, true)
					backtestMain.slBuy = conversion.Round(lowerChannel-channelInPips, roundPlace)
					backtestMain.buyFlag = false
					backtestMain.universalFlag = false
				} else if (*dt).PriceTf[i] < lowerChannel && backtestMain.sellFlag && backtestMain.universalFlag {
					backtestMain.sellPrice = (*dt).PriceTf[i]
					backtestMain.tpSell = takeProfitFimathe(lowerChannel, multiplyTpChannel, channelInPips, jpy, false)
					backtestMain.slSell = conversion.Round(upperChannel+channelInPips, roundPlace)
					backtestMain.sellFlag = false
					backtestMain.universalFlag = false
				}
				if (*dt).PriceTf[i] > upperChannel {
					lowerChannel = conversion.Round(upperChannel, roundPlace)
					for (*dt).PriceTf[i] > upperChannel {
						upperChannel = conversion.Round(upperChannel+channelInPips, roundPlace)
					}
				} else if (*dt).PriceTf[i] < lowerChannel {
					upperChannel = conversion.Round(lowerChannel, roundPlace)
					for (*dt).PriceTf[i] < lowerChannel {
						lowerChannel = conversion.Round(lowerChannel-channelInPips, roundPlace)
					}
				}

			} else if (*dt).PriceTf[i] == (*dt).PriceTf[i-1] {
				if !backtestMain.universalFlag && !backtestMain.buyFlag {
					if checkTpslExit((*dt).Open[i], (*dt).High[i], (*dt).Low[i], (*dt).Close[i], backtestMain.slBuy, false) { // Stop Loss
						// Ccopy the tp and sl and put in the next trade
						backtestMain.capital, backtestMain.buyResult = financeCalculation(backtestMain.capital, backtestMain.slBuy, backtestMain.buyPrice, backtestMain.slBuy, balance, leverage)
						updateBacktestHeart(&backtestMain, "buy")
					} else if checkTpslExit((*dt).Open[i], (*dt).High[i], (*dt).Low[i], (*dt).Close[i], backtestMain.tpBuy, true) { // Take Profit
						backtestMain.capital, backtestMain.buyResult = financeCalculation(backtestMain.capital, backtestMain.tpBuy, backtestMain.buyPrice, backtestMain.tpBuy, balance, leverage)
						updateBacktestHeart(&backtestMain, "buy")
					}
				}
				if !backtestMain.universalFlag && !backtestMain.sellFlag {
					if checkTpslExit((*dt).Open[i], (*dt).High[i], (*dt).Low[i], (*dt).Close[i], backtestMain.slSell, true) { // Stop Loss
						backtestMain.capital, backtestMain.sellResult = financeCalculation(backtestMain.capital, backtestMain.sellPrice, backtestMain.slSell, backtestMain.sellPrice, balance, leverage)
						updateBacktestHeart(&backtestMain, "sell")
					} else if checkTpslExit((*dt).Open[i], (*dt).High[i], (*dt).Low[i], (*dt).Close[i], backtestMain.tpSell, false) { // Take Profit
						backtestMain.capital, backtestMain.sellResult = financeCalculation(backtestMain.capital, backtestMain.sellPrice, backtestMain.tpSell, backtestMain.sellPrice, balance, leverage)
						updateBacktestHeart(&backtestMain, "sell")
					}
				}
			}
		}
		if (*dt).TimeHour[i] >= 20 && (*dt).TimeHour[i] <= 23 && (*dt).TimeWeekDays[i] == 5 {
			if !backtestMain.buyFlag {
				backtestMain.capital, backtestMain.buyResult = financeCalculation(backtestMain.capital, (*dt).Open[i], backtestMain.buyPrice, backtestMain.buyPrice, balance, leverage)
				updateBacktestHeart(&backtestMain, "buy")
			}
			if !backtestMain.sellFlag {
				backtestMain.capital, backtestMain.sellResult = financeCalculation(backtestMain.capital, backtestMain.sellPrice, (*dt).Open[i], backtestMain.sellPrice, balance, leverage)
				updateBacktestHeart(&backtestMain, "sell")
			}
			dayTeste = false
		}
	}

	var totalTrades []float64 = conversion.RemoveExcessZeros(backtestMain.totalTrades)
	var buyTrades []float64 = conversion.RemoveExcessZeros(backtestMain.buyTrades)
	var sellTrades []float64 = conversion.RemoveExcessZeros(backtestMain.sellTrades)

	return totalTrades, buyTrades, sellTrades
}

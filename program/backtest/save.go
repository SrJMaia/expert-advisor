package backtest

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/SrJMaia/expert-advisor/program/check"
)

func SaveBacktest(tot *[]float64, buy *[]float64, sell *[]float64, columns string) {

	var totValue = make([]string, len(*tot)+1)
	var buyValue = make([]string, len(*buy)+1)
	var sellValue = make([]string, len(*sell)+1)
	var indexes = make([]string, len(*tot)*2)

	var _, errCheck = os.Stat("C:/Users/johnk/Google Drive/Programming/Dados/BacktestResults.csv")
	if os.IsNotExist(errCheck) {
		var file, err = os.OpenFile("C:/Users/johnk/Google Drive/Programming/Dados/BacktestResults.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		var writer = csv.NewWriter(file)
		for i := 0; i < len(indexes); i++ {
			indexes[i] = fmt.Sprintf("%d", i)
		}

		err = writer.Write(indexes)
		check.MyCheckingError(err)
		writer.Flush()
		file.Close()
	}

	var file, err = os.OpenFile("C:/Users/johnk/Google Drive/Programming/Dados/BacktestResults.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	check.MyCheckingError(err)
	defer file.Close()

	var writer = csv.NewWriter(file)
	defer writer.Flush()

	totValue[0] = columns + "Tot"
	buyValue[0] = columns + "Buy"
	sellValue[0] = columns + "Sell"

	for i := 1; i < len(*tot)+1; i++ {
		totValue[i] = strconv.FormatFloat((*tot)[i-1], 'E', -1, 64)
	}
	for i := 1; i < len(*buy)+1; i++ {
		buyValue[i] = strconv.FormatFloat((*buy)[i-1], 'E', -1, 64)
	}
	for i := 1; i < len(*sell)+1; i++ {
		sellValue[i] = strconv.FormatFloat((*sell)[i-1], 'E', -1, 64)
	}

	err = writer.Write(totValue)
	check.MyCheckingError(err)

	err = writer.Write(buyValue)
	check.MyCheckingError(err)

	err = writer.Write(sellValue)
	check.MyCheckingError(err)

}

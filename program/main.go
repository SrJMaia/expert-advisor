package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// const (
// 	jpy             = false
// 	leverage        = false
// 	saveCsvBacktest = false
// 	saveCsvAnalysis = false
// 	tf              = "H1"
// 	capital         = 1000.
// )

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

const (
	host     = "localhost"
	database = "prices"
	user     = "root"
	password = "root"
)

type MyValues struct {
	anything int
	Time     sql.NullTime `db:"date_hour"`
	Open     float64      `db:"open_prices"`
	High     float64      `db:"high_prices"`
	Low      float64      `db:"low_prices"`
	Close    float64      `db:"close_prices"`
}

func main() {

	/*
		Faço a transformação no codigo
		Dividir data em duas struct
		Retirar atrtpsl?
			- Antes de calcular o tpsl, checar se é fix tpsl
		Mudar backtest?
			- Após checar se o candle é do tf ou nao, ou seja, o primeiro if else
			- Dividir o codigo posterior em pequenas funções, e chamar cada uma no qual queira

	*/

	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3307)/%s?parseTime=true", user, password, host, database)

	// var db, err = sql.Open("mysql", connectionString)
	// checkError(err)

	// defer db.Close()
	// fmt.Println("Connected")

	var start = time.Now()
	var db, err = sqlx.Open("mysql", connectionString)
	checkError(err)

	var a []MyValues
	err = db.Select(&a, "SELECT eurusd.date_hour,eurusd.open_prices, eurusd.low_prices, eurusd.high_prices, eurusd.close_prices FROM prices.eurusd")
	checkError(err)
	var elapsed = time.Since(start)
	fmt.Println(elapsed)
	fmt.Println(a[0].Time.Time)
	fmt.Println(a[0].Time.Time.Minute())
	fmt.Printf("%T\n", a[0].Time)
	fmt.Println(a[0])
	fmt.Println(a[1])

	// var filePath = "/home/maia/Documents/Project/expert-advisor-database/EURUSD_Clean_20y.csv"

	// table.StartHead()
	// start := time.Now()

	// csvFile, err := os.Open(filePath)
	// check.MyCheckingError(err)
	// defer csvFile.Close()

	// csvLines, err := csv.NewReader(csvFile).ReadAll()
	// check.MyCheckingError(err)

	// fmt.Println(len(csvLines))

	// var stringSize int
	// var stringToSlicce string

	// var openValue float64
	// var highValue float64
	// var lowValue float64
	// var closeValue float64

	// var allValues []string

	// for i, v := range csvLines {
	// 	if i < 1 {
	// 		continue
	// 	}

	// 	openValue, err = strconv.ParseFloat(v[1], 64)
	// 	check.MyCheckingError(err)

	// 	highValue, err = strconv.ParseFloat(v[2], 64)
	// 	check.MyCheckingError(err)

	// 	lowValue, err = strconv.ParseFloat(v[3], 64)
	// 	check.MyCheckingError(err)

	// 	closeValue, err = strconv.ParseFloat(v[4], 64)
	// 	check.MyCheckingError(err)

	// 	stringToSlicce = fmt.Sprintf("(STR_TO_DATE('%s', '%%Y-%%m-%%d %%H:%%i:%%s' ) ,%.5f, %.5f, %.5f, %.5f)", v[0], openValue, highValue, lowValue, closeValue)

	// 	if stringSize+len(stringToSlicce) >= 4000000 {
	// 		x := fmt.Sprintf("INSERT INTO prices.eurusd(eurusd.date_hour, eurusd.open_prices, eurusd.high_prices, eurusd.low_prices, eurusd.close_prices) VALUES %s", strings.Join(allValues, ","))
	// 		insert, err := db.Query(x)
	// 		checkError(err)
	// 		insert.Close()
	// 		allValues = nil
	// 		stringSize = 0
	// 	}

	// 	allValues = append(allValues, stringToSlicce)
	// 	stringSize += len(stringToSlicce)

	// 	if i == len(csvLines)-1 {
	// 		x := fmt.Sprintf("INSERT INTO prices.eurusd(eurusd.date_hour, eurusd.open_prices, eurusd.high_prices, eurusd.low_prices, eurusd.close_prices) VALUES %s", strings.Join(allValues, ","))
	// 		insert, err := db.Query(x)
	// 		checkError(err)
	// 		insert.Close()
	// 		allValues = nil
	// 		stringSize = 0
	// 	}
	// }

	// elapsed := time.Since(start)

	// table.StartBody("Successfully Read CSV file.", elapsed)

	// -------------------------------------------------------------------
	// var myAnalysisStruct = myanalysis.LayoutAnalysis{}
	// var myData = data.LayoutData{}

	// data.ReadData("/home/maia/Documents/Project/expert-advisor-database/EURUSD_Clean_20y.csv", &myData)
	// var start = time.Now()
	// data.NormalizePricesTFAndDateTime(&myData, tf)
	// var elapsed = time.Since(start)
	// fmt.Println("Time to normalize:", elapsed)

	// fmt.Println(myData.SizeTimeFrame)
	// indicator.AtrTpsl(&myData, 14, tf, jpy)
	// //indicator.FixTpsl(&myData, 10, tf, jpy)

	// start = time.Now()
	// // for slow := 60; slow < 390; slow += 30 {
	// // 	//for fast := 30; fast < slow; fast += 30 {
	// // 	//for i := 90; i < 180; i += 30 {
	// // 	table.LoopStepsPrint(float64(slow))
	// // 	indicator.StrategyOneDecyclerPriceCross(&myData, float64(slow), tf, jpy)
	// // 	if myData.SizeBuy < 2 || myData.SizeSell < 2 {
	// // 		continue
	// // 	}
	// // 	backtest.Optimize(&myData, &myAnalysisStruct, "hedging", capital, jpy, saveCsvBacktest, leverage, tf, float64(slow))
	// // 	//}
	// // 	//}
	// // }

	// for slow := 60; slow < 300; slow += 60 {
	// 	indicator.StrategyOneDecyclerPriceCross(&myData, float64(slow), tf, jpy)
	// 	if myData.SizeBuy < 2 || myData.SizeSell < 2 {
	// 		continue
	// 	}
	// 	backtest.Optimize(&myData, &myAnalysisStruct, "fimathe", capital, jpy, saveCsvBacktest, leverage, tf, float64(slow))
	// }

	// // indicator.StrategyOneDecyclerPriceCross(&myData, 30., tf, jpy)
	// // backtest.Optimize(&myData, &myAnalysisStruct, "fimathe", capital, jpy, saveCsvBacktest, leverage, tf)

	// elapsed = time.Since(start)
	// table.StartBody("Successfully Backtested.", elapsed)
	// myanalysis.FilterOtimization(&myAnalysisStruct, capital, float64(myData.SizeTimeFrame), saveCsvAnalysis)
}

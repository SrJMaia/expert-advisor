package main

import (
	"log"
	"os"
	"os/exec"

	_ "github.com/go-sql-driver/mysql"
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

// const (
// 	host     = "localhost"
// 	database = "test_db"
// 	user     = "root"
// 	password = "root"
// )

func main() {

	cmd := exec.Command("cd ..")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	// var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3307)/", user, password, host)

	// var db, err = sql.Open("mysql", connectionString)
	// checkError(err)

	// defer db.Close()
	// fmt.Println("Connected")

	// // Create DataBase
	// _, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + database)
	// checkError(err)

	// db.Close()

	// connectionString = fmt.Sprintf("%s:%s@tcp(%s:3307)/%s", user, password, host, database)
	// db, err = sql.Open("mysql", connectionString)
	// checkError(err)
	// defer db.Close()

	// // Drop previous table of same name if one exists.
	// _, err = db.Exec("DROP TABLE IF EXISTS test;")
	// checkError(err)
	// fmt.Println("Finished dropping table (if existed).")

	// // Create table.
	// _, err = db.Exec("CREATE TABLE test (id serial PRIMARY KEY, name VARCHAR(50), quantity INTEGER);")
	// checkError(err)
	// fmt.Println("Finished creating table.")

	// // Insert some data into table.
	// sqlStatement, err := db.Prepare("INSERT INTO test (name, quantity) VALUES (?, ?);")
	// checkError(err)
	// res, err := sqlStatement.Exec("banana", 150)
	// checkError(err)
	// rowCount, err := res.RowsAffected()
	// checkError(err)
	// fmt.Printf("Inserted %d row(s) of data.\n", rowCount)

	// res, err = sqlStatement.Exec("orange", 154)
	// checkError(err)
	// rowCount, err = res.RowsAffected()
	// checkError(err)
	// fmt.Printf("Inserted %d row(s) of data.\n", rowCount)

	// res, err = sqlStatement.Exec("apple", 100)
	// checkError(err)
	// rowCount, err = res.RowsAffected()
	// checkError(err)
	// fmt.Printf("Inserted %d row(s) of data.\n", rowCount)
	// fmt.Println("Done.")

	// var myAnalysisStruct = myanalysis.LayoutAnalysis{}
	// var myData = data.LayoutData{}

	// data.ReadData("/home/maia/Documents/Git/expert-advisor/data/EURUSD_Clean_20y.csv", &myData)

	// data.NormalizePricesTFAndDateTime(&myData, tf)
	// indicator.AtrTpsl(&myData, 14, tf, jpy)
	// //indicator.FixTpsl(&myData, 10, tf, jpy)

	// var start = time.Now()
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

	// var elapsed = time.Since(start)
	// table.StartBody("Successfully Backtested.", elapsed)
	// myanalysis.FilterOtimization(&myAnalysisStruct, capital, float64(myData.SizeTimeFrame), saveCsvAnalysis)
}

package check

import (
	"log"
	"math"
)

func MyCheckingError(err error) {
	if err != nil {
		log.Panic("Error:", err)
	}
}

func MyCheckingNan(value float64, index int) {
	if math.IsNaN(value) {
		log.Println("Nan found, index:", index)
	}
}

func MyCheckingFinanceCalculation(total float64, initialPrice float64, finalPrice float64, result float64) {
	if math.IsNaN(total) || initialPrice == 0 || finalPrice == 0 || math.IsNaN(result) {
		log.Panic("Finance Calculation Error. NaN")
	}
}

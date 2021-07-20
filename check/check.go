package check

import (
	"log"
	"math"
)

func MyCheckingError(err error) {
	if err != nil {
		log.Println("Error:", err)
	}
}

func MyCheckingNan(value float64, index int) {
	if math.IsNaN(value) {
		log.Println("Nan found, index:", index)
	}
}

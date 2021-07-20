package check

import (
	"math"
	"testing"
)

func TestMyCheckingNan(t *testing.T) {
	var arrayTest = [5]float64{0, 1, math.NaN(), 3, 4}
	for i, v := range arrayTest {
		MyCheckingNan(v, i)
	}
}

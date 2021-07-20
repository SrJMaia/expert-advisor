package myanalysis

import (
	"testing"
)

func TestAnnualizedReturn(t *testing.T) {
	var backtest = []float64{1000, 3000}
	var leng = 132744
	var tf = "H1"
	var result = AnnualizedReturn(backtest, leng, tf)
	if result != 7.52 {
		t.Error("Expected:", 7.52, "Got:", result)
	}

}

func TestMaximumDrawdown(t *testing.T) {
	x, y := MaximumDrawdown([]float64{-1, 1, 2, 3, 4, 5, 6, 4, 3, 5, 6, 7, 2, 8, -10, 1, 1, 1, 2, 2, 3})
	if x != -225 {
		t.Error("X - Expected:", -87.5, "Got:", x)
	} else if y != 18 {
		t.Error("Y - Expected:", 18, "Got:", y)
	}

}

func TestNetProfit(t *testing.T) {
	type test struct {
		data   []float64
		answer float64
	}

	var testsStruct = []test{
		{
			data:   []float64{100, 120},
			answer: 20,
		},
		{
			data:   []float64{100, 0, 10, 15, -20, -150, 120},
			answer: 20,
		},
		{
			data:   []float64{100, 0},
			answer: -100,
		},
		{
			data:   []float64{100, -100},
			answer: -200,
		},
	}

	for _, v := range testsStruct {
		var x = NetProfit(v.data)
		if x != v.answer {
			t.Error("Expected:", v.answer, "Got:", x)
		}
	}
}

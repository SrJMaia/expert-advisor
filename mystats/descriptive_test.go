package mystats

import "testing"

func TestQuartile(t *testing.T) {
	var oddSlice = []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var oddQ1, oddQ3 = Quartile(oddSlice)
	if oddQ1 != 2 {
		t.Error("Odd - Expected:", 2, "Got:", oddQ1)
	} else if oddQ3 != 8 {
		t.Error("Odd - Expected:", 8, "Got:", oddQ3)
	}
	var evenSlice = []float64{0, 1, 10, 3, 4, 11, 6, 7, 8, 9, 2, 5}
	var evenQ1, evenQ3 = Quartile(evenSlice)
	if evenQ1 != 2.5 {
		t.Error("Even - Expected:", 2.5, "Got:", evenQ1)
	} else if evenQ3 != 8.5 {
		t.Error("Even - Expected:", 8.5, "Got:", evenQ3)
	}
}

func TestMaxBelowZero(t *testing.T) {
	x := MaxBelowZero([]float64{0, 1, 2, -1, -2, -10, 10, 20}...)
	if x != -1 {
		t.Error("Expected:", -1, "Got:", x)
	}
}

func TestMedian(t *testing.T) {
	var oddSlice = []float64{5, 4, 0, 10, 7}
	var x = Median(oddSlice)
	if x != 5 {
		t.Error("Odd - Expected:", 5, "Got:", x)
	}
	var evenSlice = []float64{10, 4, 1, 2}
	var y = Median(evenSlice)
	if y != 3 {
		t.Error("Even - Expected:", 3, "Got:", y)
	}
}

func TestStandardDeviationPopulation(t *testing.T) {
	var slice = []float64{600, 470, 170, 430, 300}
	var x = StandardDeviationPopulation(slice)
	if x != 147.32277 {
		t.Error("Expected:", 147.32277, "Got:", x)
	}
}

func TestStandardDeviationSample(t *testing.T) {
	var slice = []float64{600, 470, 170, 430, 300}
	var x = StandardDeviationSample(slice)
	if x != 164.71187 {
		t.Error("Expected:", 164.71187, "Got:", x)
	}
}

func TestVarianceSample(t *testing.T) {
	var slice = []float64{600, 470, 170, 430, 300}
	var x = VarianceSample(slice)
	if x != 27130 {
		t.Error("Expected:", 27130, "Got:", x)
	}
}

func TestVariancePopulation(t *testing.T) {
	var slice = []float64{600, 470, 170, 430, 300}
	var x = VariancePopulation(slice)
	if x != 21704 {
		t.Error("Expected:", 21704, "Got:", x)
	}
}

func TestMax(t *testing.T) {
	var slice = []float64{-1, -2, 1, 10, 2, 5, 20, -20}
	var x = Max(slice...)
	if x != 20 {
		t.Error("Expected:", 20, "Got:", x)
	}
}

func TestMin(t *testing.T) {
	var slice = []float64{-1, -2, 1, 10, 2, 5, 20, -20}
	var x = Min(slice...)
	if x != -20 {
		t.Error("Expected:", -20, "Got:", x)
	}
}

func TestMinAboveZero(t *testing.T) {
	var slice = []float64{-1, -2, 1, 10, 2, 5, 20, -20}
	var x = MinAboveZero(slice...)
	if x != 1 {
		t.Error("Expected:", 1, "Got:", x)
	}
}

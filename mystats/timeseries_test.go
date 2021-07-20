package mystats

import (
	"math"
	"testing"
)

func TestMaxValueOverPeriod(t *testing.T) {
	var array = []float64{10, 10, 11, 11, 11, 10, 9, 5, 10, 8, 1, 20, 2}
	x := MaxValueOverPeriod(&array, 2)
	var answer = []float64{10, 10, 10, 11, 11, 11, 11, 10, 9, 10, 10, 8, 20}
	for i := range answer {
		if x[i] != answer[i] {
			t.Error("Expected:", answer, "Got:", x)
		}
	}
}

func TestMinValueOverPeriod(t *testing.T) {
	var array = []float64{10, 10, 11, 11, 11, 10, 9, 5, 10, 8, 1, 20, 2}
	x := MinValueOverPeriod(&array, 2)
	var answer = []float64{10, 10, 10, 10, 11, 11, 10, 9, 5, 5, 8, 1, 1}
	for i := range answer {
		if x[i] != answer[i] {
			t.Error("Expected:", answer, "Got:", x)
		}
	}
}

func TestGetAmountValuesOverX(t *testing.T) {
	var array = []float64{10, 10, 11, 11, 11, 10, 9}
	x := GetAmountValuesOverX(&array, 10)
	if x != 3 {
		t.Error("Expected:", 3, "Got:", x)
	}
}

func TestGetValuesByIndexString(t *testing.T) {
	var array = []string{"10", "20", "30", "40", "50", "60", "70", "80", "90"}
	var indexes = []int{2, 5, 1}
	x := GetValuesByIndexString(&array, indexes)
	answer := []string{"30", "60", "20"}
	for i := range x {
		if x[i] != answer[i] {
			t.Error("Expected:", answer, "Got:", x)
		}
	}
}

func TestGetValuesByIndexFloat(t *testing.T) {
	var array = []float64{10, 20, 30, 40, 50, 60, 70, 80, 90}
	var indexes = []int{2, 5, 1}
	x := GetValuesByIndexFloat(&array, indexes)
	answer := []float64{30, 60, 20}
	for i := range x {
		if x[i] != answer[i] {
			t.Error("Expected:", answer, "Got:", x)
		}
	}
}

func TestPctChange(t *testing.T) {
	x := PctChange([]float64{10, 14, 20, 25, 12.5, 13, 0, 50})
	answer := []float64{0, 0.4, 0.42857142857142855, 0.25, -0.5, 0.04, -1, math.Inf(1)}
	for i := range x {
		if x[i] != answer[i] {
			t.Error("Expected:", answer, "Got:", x)
		}
	}
}

func TestDiff(t *testing.T) {
	x := Diff([]float64{10, 14, 20, 25, 12.5, 13, 0, 50})
	answer := []float64{0, 4, 6, 5, -12.5, 0.5, -13, 50}
	for i := range x {
		if x[i] != answer[i] {
			t.Error("Expected:", answer, "Got:", x)
		}
	}
}

func TestMaximumAccumulate(t *testing.T) {
	x := MaximumAccumulate([]float64{-1, 1, 2, 3, 4, 5, 6, 4, 3, 5, 6, 7, 2, 8, 1, 1, 1, 1, 2, 2, 3})
	answer := []float64{-1, 1, 2, 3, 4, 5, 6, 6, 6, 6, 6, 7, 7, 8, 8, 8, 8, 8, 8, 8, 8}
	for i := range x {
		if x[i] != answer[i] {
			t.Error("Expected:", answer, "Got:", x)
		}
	}
}

func TestFindFirstIndexValue(t *testing.T) {
	x := FindFirstIndexValue(&[]float64{-1, 0, 5, 10, -50}, 10)
	if x != 3 {
		t.Error("Expected:", 3, "Got:", x)
	}
}

func TestFindIndexesValue(t *testing.T) {
	x := FindIndexesValue(&[]float64{-1, 0, 5, 10, 10, -50}, 10)
	var answer = []int{3, 4}
	for i := range x {
		if x[i] != answer[i] {
			t.Error("Expected:", answer, "Got:", x)
		}
	}
}

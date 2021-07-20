package conversion

import "testing"

func TestRound(t *testing.T) {
	var result = Round(1.12345, 3)
	if result != 1.123 {
		t.Error("Expected:", 1.123, "Got:", result)
	}
}

func TestRemoveExcessZeros(t *testing.T) {
	var result = RemoveExcessZeros([]float64{1, 2, 3, 4, 5, 0, 0, 0, 0, 0})
	var answer = []float64{1, 2, 3, 4, 5}
	for i := range answer {
		if result[i] != answer[i] {
			t.Error("Expected:", answer, "Got:", result)
		}
	}
}

func TestFillWithValue(t *testing.T) {
	var result = FillWithValue([]float64{1, 2, 3, 4, 5, 0, 0, 0, 0, 0}, 0, 10)
	var answer = []float64{1, 2, 3, 4, 5, 10, 10, 10, 10, 10}
	for i := range answer {
		if result[i] != answer[i] {
			t.Error("Expected:", answer, "Got:", result)
		}
	}
}

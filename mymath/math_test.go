package mymath

import "testing"

func TestQuartile(t *testing.T) {
	var result = DegreesToRadians(60)
	if result != 1.0471975511965976 {
		t.Error("Expected:", 1.0471975511965976, "Got:", result)
	}
}

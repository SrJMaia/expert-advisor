package mystats

import (
	"math"
	"sort"

	"github.com/SrJMaia/expert-advisor/program/conversion"
)

func Quartile(s []float64) (float64, float64) {
	sort.Float64s(s)
	var size = float64(len(s))
	var q1 float64
	var q3 float64
	if len(s)%2 == 0 {
		q1 = (s[int((size+1)*.25)-1] + s[int((size+1)*.25)]) / 2
		q3 = (s[int((size+1)*.75)-1] + s[int((size+1)*.75)]) / 2
		return q1, q3
	}
	q1 = s[int((size)*.25)]
	q3 = s[int((size)*.75)]
	return q1, q3
}

func Median(s []float64) float64 {
	sort.Float64s(s)
	if len(s)%2 == 0 {
		var middle = len(s) / 2
		return (s[middle-1] + s[middle]) / 2
	}
	return s[int(len(s)/2)]
}

func Mean(s []float64) float64 {
	var avg float64
	for i := range s {
		avg += s[i]
	}
	avg = avg / float64(len(s))

	return avg
}

func VarianceSample(s []float64) float64 {
	var avg = Mean(s)
	var variance float64
	for i := range s {
		variance = variance + (math.Pow(s[i]-avg, 2))
	}
	variance = variance / float64(len(s)-1)
	return variance
}

func VariancePopulation(s []float64) float64 {
	var avg = Mean(s)
	var variance float64
	for i := range s {
		variance = variance + (math.Pow(s[i]-avg, 2))
	}
	variance = variance / float64(len(s))
	return variance
}

func StandardDeviationSample(s []float64) float64 {
	var std = VarianceSample(s)
	return conversion.Round(math.Sqrt(std), 5)
}

func StandardDeviationPopulation(s []float64) float64 {
	var std = VariancePopulation(s)
	return conversion.Round(math.Sqrt(std), 5)
}

func Max(n ...float64) float64 {
	sort.Float64s(n)
	if len(n) == 0 {
		return 0.
	}
	return n[len(n)-1]
}

func MaxBelowZero(n ...float64) float64 {
	var maxValue = -math.MaxFloat64
	if len(n) == 1 {
		return n[0]
	} else if len(n) == 0 {
		return 0.
	}
	for i := range n {
		if n[i] >= 0 {
			continue
		}
		if n[i] > maxValue || !math.IsNaN(n[i]) || !math.IsInf(n[i], 0) {
			maxValue = n[i]
		}
	}
	return maxValue
}

func Min(n ...float64) float64 {
	sort.Float64s(n)
	if len(n) == 0 {
		return 0.
	}
	return n[0]
}

func MaxMin(n []float64) (float64, float64) {
	if len(n) == 0 {
		return 0., 0.
	}
	var max float64
	var min = math.MaxFloat64
	for _, value := range n {
		if value > max {
			max = value
		}
		if value < min {
			min = value
		}
	}
	return max, min
}

func MinAboveZero(n ...float64) float64 {
	var minValue float64
	if len(n) == 1 {
		return n[0]
	} else if len(n) == 0 {
		return 0.
	}
	for i := range n {
		if i == 0 {
			minValue = math.MaxFloat64
			continue
		}
		if 0 < n[i] && n[i] < minValue || !math.IsNaN(n[i]) || !math.IsInf(n[i], 0) {
			minValue = n[i]
		}
	}
	if minValue == math.MaxFloat64 {
		minValue = 0
	}
	return minValue
}

package mystats

func MaxValueOverPeriod(array *[]float64, period int) []float64 {
	var maxValue = make([]float64, len(*array))
	for i := 0; i < len(*array); i++ {
		if i < period {
			maxValue[i] = (*array)[i]
			continue
		}
		maxValue[i] = Max((*array)[i-period : i]...)
	}
	return maxValue
}

func MinValueOverPeriod(array *[]float64, period int) []float64 {
	var maxValue = make([]float64, len(*array))
	for i := 0; i < len(*array); i++ {
		if i < period {
			maxValue[i] = (*array)[i]
			continue
		}
		maxValue[i] = Min((*array)[i-period : i]...)
	}
	return maxValue
}

func GetAmountValuesOverX(array *[]float64, value float64) float64 {
	var totValues float64
	for i := range *array {
		if (*array)[i] > value {
			totValues++
		}
	}
	return totValues
}

func GetValuesByIndexFloat(array *[]float64, indexes []int) []float64 {
	var results = make([]float64, len(indexes))
	for i, v := range indexes {
		results[i] = (*array)[v]
	}
	return results
}

func GetValuesByIndexString(array *[]string, indexes []int) []string {
	var results = make([]string, len(indexes))
	for i, v := range indexes {
		results[i] = (*array)[v]
	}
	return results
}

func Diff(array []float64) []float64 {
	var arrayDiff = make([]float64, len(array))
	for i := range array {
		if i < 1 {
			arrayDiff[i] = 0.
			continue
		}
		arrayDiff[i] = (array)[i] - (array)[i-1]
	}
	return arrayDiff
}

func PctChange(array []float64) []float64 {
	var arrayDiff = make([]float64, len(array))
	for i := range array {
		if i < 1 {
			arrayDiff[i] = 0.
			continue
		}
		arrayDiff[i] = (array[i] - array[i-1]) / array[i-1]
	}
	return arrayDiff
}

func MaximumAccumulate(array []float64) []float64 {
	// Same as np.maximum.accumulate from python
	var maxValue = array[0]
	var newArray = make([]float64, len(array))
	for i, v := range array {
		if v > maxValue {
			maxValue = v
			newArray[i] = v
		} else {
			newArray[i] = maxValue
		}
	}
	return newArray
}

func FindFirstIndexValue(array *[]float64, value float64) int {
	// Find the first index where a value appears
	var index int
	for i := range *array {
		if (*array)[i] == value {
			index = i
			break
		}
	}
	return index
}

func FindIndexesValue(array *[]float64, value float64) []int {
	// Find all the indexes where a value appears
	var index []int
	for i := range *array {
		if (*array)[i] == value {
			index = append(index, i)
		}
	}
	return index
}

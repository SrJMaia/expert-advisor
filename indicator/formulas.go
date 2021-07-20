package indicator

import (
	"math"

	"github.com/SrJMaia/expert-advisor/conversion"
	"github.com/SrJMaia/expert-advisor/mymath"
)

func Decycler(price *[]float64, cutoff float64, jpy bool) []float64 {
	var alpha1 = (math.Cos(mymath.DegreesToRadians(360./cutoff)) + math.Sin(mymath.DegreesToRadians(360./cutoff)) - 1) / math.Cos(mymath.DegreesToRadians(360./cutoff))
	var decyclerSlice = make([]float64, len(*price))
	var roundPLace int
	if jpy {
		roundPLace = 3
	} else {
		roundPLace = 5
	}
	for i := range *price {
		if i < 1 {
			decyclerSlice[i] = (*price)[i]
			continue
		}
		decyclerSlice[i] = conversion.Round((alpha1/2)*((*price)[i]+(*price)[i-1])+(1-alpha1)*decyclerSlice[i-1], roundPLace)
	}
	return decyclerSlice
}

func DecyclerOscilator(price *[]float64, hpperiod1 float64, hpperiod2 float64, jpy bool) []float64 {
	var alpha1 = (math.Cos(mymath.DegreesToRadians(.707*360./hpperiod1)) + math.Sin(mymath.DegreesToRadians(.707*360./hpperiod1)) - 1.) / math.Cos(mymath.DegreesToRadians(.707*360/hpperiod1))
	var alpha2 = (math.Cos(mymath.DegreesToRadians(.707*360./hpperiod2)) + math.Sin(mymath.DegreesToRadians(.707*360./hpperiod2)) - 1.) / math.Cos(mymath.DegreesToRadians(.707*360/hpperiod2))
	var hp1 = make([]float64, len(*price))
	var hp2 = make([]float64, len(*price))
	var decyclerOscSlice = make([]float64, len(*price))
	var roundPLace int
	if jpy {
		roundPLace = 3
	} else {
		roundPLace = 5
	}
	for i := range *price {
		if i < 2 {
			hp1[i] = (*price)[i]
			hp2[i] = (*price)[i]
			continue
		}
		hp1[i] = (1.-alpha1/2.)*(1.-alpha1/2.)*((*price)[i]-2*(*price)[i-1]+(*price)[i-2]) + 2*(1.-alpha1)*hp1[i-1] - (1.-alpha1)*(1.-alpha1)*hp1[i-2]
		hp2[i] = (1.-alpha2/2.)*(1.-alpha2/2.)*((*price)[i]-2*(*price)[i-1]+(*price)[i-2]) + 2*(1.-alpha2)*hp2[i-1] - (1.-alpha2)*(1.-alpha2)*hp2[i-2]

		decyclerOscSlice[i] = conversion.Round(hp2[i]-hp1[i], roundPLace)
	}
	return decyclerOscSlice
}

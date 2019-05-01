package maths

import (
	"errors"
	"fmt"
	"math"
)

// RMSError calculates the RMS error from two float64 slices of equal length
func RMSError(series1, series2 []float64) (float64, error) {
	if len(series1) != len(series2) {
		return -1., errors.New(fmt.Sprintf("Un-equal length of input arrays being %d and %d", len(series1), len(series2)))
	}
	nElms := len(series1)
	diff := 0.
	for i, elm := range series1 {
		diff = diff + math.Pow((elm-series2[i]), 2.)
	}
	return math.Sqrt(math.Abs(diff) / float64(nElms)), nil
}

func maxFloatInSlice(nums []float64) float64 {
	largest := nums[0]
	for _, num := range nums {
		if num > largest {
			largest = num
		}
	}
	return largest
}

func minFloatInSlice(nums []float64) float64 {
	smallest := nums[0]
	for _, num := range nums {
		if num < smallest {
			smallest = num
		}
	}
	return smallest
}

// EnclosedWithin returns the values of an array which are within the range of values in another array
func EnclosedWithin(serie1 []float64, serie2 []float64) []float64 {
	max1 := maxFloatInSlice(serie1)
	min1 := minFloatInSlice(serie1)

	max2 := maxFloatInSlice(serie2)
	min2 := minFloatInSlice(serie2)

	max := math.Min(max1, max2)
	min := math.Max(min1, min2)

	results := []float64{}

	for _, val := range serie1 {
		if val >= min && val <= max {
			results = append(results, val)
		}
	}
	return results

}

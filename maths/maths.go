package maths

import (
	"log"
	"math"
)

func ArgMin(vals []float64) (int, float64) {
	min := vals[0]
	minInd := 0
	for index, val := range vals {
		if val < min {
			min = val
			minInd = index
		}
	}
	return minInd, min
}

func SortBy(vals []float64, indices []int) []float64 {
	if !(len(vals)==len(indices)){
		log.Panic("Values and indices are not of equal length")
	}
	organized := make([]float64, len(vals))
	for i, index := range indices {
		organized[i] = vals[index]
	}
	return organized
}

func ContainsElementWithinTol(values[]float64,val float64,tol float64) bool{
	for _,value:=range values{
		if math.Abs(value-val)<tol{
			return true
		}
	}
	return false
}

func Arange(min float64, max float64, step float64) []float64 {
	results := []float64{}
	val := min
	for val <= max {
		results = append(results, val)
		val += step


	}
	return results
}

func Linspace(min float64, max float64, Nsteps int) []float64 {
	results := []float64{}
	step := (max-min)/float64(Nsteps-1)
	val := min
	for val <= max {
		results = append(results, val)
		val += step


	}
	return results
}
	


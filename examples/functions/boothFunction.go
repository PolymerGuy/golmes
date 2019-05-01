package functions

import (
	"github.com/PolymerGuy/tinyOpt/examples"
	"math"
)

func booth(args...float64) []float64{
	value := math.Pow(args[0]+2*args[1]-7,2)+math.Pow(2*args[0]+args[1]-5,2)
	return []float64{value}
}

func NewBoothFunction() main.Expression {
	return main.Expression{
		expr:booth,
		nArgs:2,
	}
}


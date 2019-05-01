package functions

import "github.com/PolymerGuy/tinyOpt/examples"

func unity(args...float64) []float64{
	return args
}

func NewUnityFunction() main.Expression {
	return main.Expression{
		expr:  unity,
		nArgs: 2,
	}
}


package functions

import (
	"github.com/PolymerGuy/tinyOpt/examples"
	"gonum.org/v1/gonum/floats"
	"testing"
)

func TestFunction(t *testing.T) {
	const tol = 1e-9
	// Test-Expression returns its arguments
	expression := main.NewUnityFunction()
	f := main.NewFunction(expression)

	f.Args([]float64{2,5})
	f.Evaluate()
	if !floats.EqualApprox(f.Value(),[]float64{2,5},tol){
		t.Errorf("Returned wrong results %v",f.Value())
	}

	// Mutate with new arguments
	f.Args([]float64{3,7})
	f.Evaluate()
	if !floats.EqualApprox(f.Value(),[]float64{3,7},tol){
		t.Errorf("Returned wrong results %v",f.Value())
	}


}




package maths

import (
	"fmt"
	"gonum.org/v1/gonum/floats"
	"testing"
)

const tol  = 1e-6
// Weak test comparing it to numpy
func TestLinspace(t *testing.T) {
	correct := []float64{0.1, 2.4, 4.7, 7. , 9.3}
	values := Linspace(0.1,9.3,5)

	diff := make([]float64,len(correct))
	floats.SubTo(diff,correct,values)
	fmt.Println(diff)

	if floats.Max(diff)>tol || floats.Min(diff)< -tol{
		t.Fail()
	}

}

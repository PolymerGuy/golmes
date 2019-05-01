package functions

import (
	"fmt"
	"github.com/PolymerGuy/tinyOpt/examples"
	"gonum.org/v1/gonum/floats"
	"testing"
)

func ExampleBoothMinima() {
	fmt.Println(main.booth(1,3))
	// Output: [0]
}


func TestBooth(t *testing.T){
	if !floats.EqualApprox(main.booth(1,3),[]float64{0},1e-9){
		t.Errorf("Did not return 0")
	}

	if floats.EqualApprox(main.booth(1,4),[]float64{0},1e-9){
		t.Errorf("Did return 0")
	}
}



package maths

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/floats"
)

func TestRMSError(t *testing.T) {
	tol := 1e-9
	serie1 := Arange(0., 10., 1.)
	serie2 := Arange(1., 11., 1.)
	serie3 := Arange(1., 10., 1.)

	// Test for equal series. Should return 0
	if val, _ := RMSError(serie1, serie1); val > tol {
		t.Errorf("Equal serie does not return a zero error")
	}

	// An offset of 1 should give a RMS error of 1
	if val, _ := RMSError(serie1, serie2); math.Abs(val-1.0) > tol {
		t.Errorf("Unity offset should yield unity error")
	}

	// Should return UnEqualLengthError when un-equal length slices are given
	_, err := RMSError(serie1, serie3)
	if err == nil {
		t.Errorf("Panic occured for unequal length slices")
	}

}

func Test_maxFloatInSlice(t *testing.T) {
	serie1 := []float64{-5., -2.0, 9.0, -12.0}
	if maxFloatInSlice(serie1) != 9.0 {
		t.Errorf("Did not find largest element in slice")
	}

}

func Test_minFloatInSlice(t *testing.T) {
	serie1 := []float64{-5., -2.0, -90.0, 12.0}
	if minFloatInSlice(serie1) != -90.0 {
		t.Errorf("Did not find smallest element in slice")
	}
}

func TestEnclosedWithin(t *testing.T) {
	serie1 := Arange(-5, 5, 1)
	serie2 := Arange(-2, 3, 1)
	serie3 := Arange(6, 12, 1)
	serie4 := Arange(-2.5, 3.5, 1)

	// Check for no union, should return empty slice
	if !floats.Equal(EnclosedWithin(serie1, serie3), []float64{}) {
		t.Errorf("Did not handle no union")
	}

	if !floats.Equal(EnclosedWithin(serie1, serie2), serie2) {
		t.Errorf("Did not handle union")
	}

	if !floats.Equal(EnclosedWithin(serie1, serie4), serie2) {
		t.Errorf("Did not handle union")
	}
}

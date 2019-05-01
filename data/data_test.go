package data

import (
	"github.com/PolymerGuy/golmes/maths"
	"gonum.org/v1/gonum/floats"
	"testing"
)

// Trivial test of interpolator where it is checked that it returns
// the values it was initialized with and the midpoints between them.
// TODO: Implement test for interpolating a harmonic
func TestSerieWithArgs_ReadAt(t *testing.T) {

	args := NewSeries([]float64{1, 2, 3, 4, 5})
	vals := NewSeries([]float64{1, 2, 3, 4, 5})

	series := SerieWithArgs{args, vals}

	res := series.ReadAt(args)
	if !floats.Equal(res, vals.Read()) {
		t.Error("Interpolator does not return the values it was initialized with")
	}

	res = series.ReadAt([]float64{1.5, 2.5, 3.5, 4.5})
	if !floats.Equal(res, []float64{1.5, 2.5, 3.5, 4.5}) {
		t.Error("Interpolator does not return the values between interpolants")
	}

}

func TestPairSyncArgs_Compare(t *testing.T) {
	//TODO: Define behavior when there are no overlapping args
	serieArgs1 := NewSeries(maths.Arange(0., 2., 0.1))
	serieArgs2 := NewSeries(maths.Arange(0., 2., 0.2))
	serieArgs3 := NewSeries(maths.Arange(-2., 0., 0.1))
	serieArgs4 := NewSeries(maths.Arange(0.5, 1.5, 0.1))
	emptyFloatSlice := NewSeries([]float64{})

	serie1 := SerieWithArgs{serieArgs1, serieArgs1}
	serie2 := SerieWithArgs{serieArgs2, serieArgs2}
	serie3 := SerieWithArgs{serieArgs3, serieArgs3}
	serie4 := SerieWithArgs{serieArgs4, serieArgs4}

	NewPairWithArgs(serie1, serie1, emptyFloatSlice).Compare()

	// Compare equal series should return 0
	if val, _ := NewPairWithArgs(serie1, serie1, emptyFloatSlice).Compare(); val != 0.0 {
		t.Errorf("Compare does not return zero error for equal series")
	}

	// Compare should evaluate all points in serieArgs1 and return 0
	if val, _ := NewPairWithArgs(serie1, serie2, emptyFloatSlice).Compare(); val != 0.0 {
		t.Errorf("Compare does not return zero error for unequal args and common vals on linear data")
	}

	// Not overlapping series should return error
	if _, err := NewPairWithArgs(serie1, serie3, emptyFloatSlice).Compare(); err == nil {
		t.Errorf("Compare does not return zero error for unequal args and common vals on linear data")
	}

	// Compare should evaluate the point at 1 in serieArgs1 and return 0
	if val, _ := NewPairWithArgs(serie1, serie4, emptyFloatSlice).Compare(); val != 0.0 {
		t.Errorf("Compare does not return zero error for unequal args and common vals on linear data")
	}

	// Compare should evaluate the point in serieArgs4 and return 0
	if val, _ := NewPairWithArgs(serie1, serie4, serieArgs4).Compare(); val != 0.0 {
		t.Errorf("Compare does not return zero error for unequal args and common vals on linear data")
	}

	// Not overlapping args should return error
	if _, err := NewPairWithArgs(serie1, serie1, serieArgs3).Compare(); err == nil {
		t.Errorf("Compare does not return zero error for unequal args and common vals on linear data")
	}

}

func TestPair_Compare(t *testing.T) {
	//TODO: Define behavior when there are no overlapping args
	serie1 := NewSeries(maths.Arange(0., 2., 0.1))
	serie2 := NewSeries(maths.Arange(1., 3., 0.1))
	serie3 := NewSeries(maths.Arange(1., 3., 0.2))

	// Compare equal series should return 0
	if val, _ := NewPair(serie1, serie1).Compare(); val != 0.0 {
		t.Errorf("Compare does not return zero error for equal series")
	}

	// Compare equal series should return 0
	if val, _ := NewPair(serie1, serie2).Compare(); !floats.EqualWithinAbs(val, 1.0, 1e-6) {
		t.Errorf("Compare does not return 1 error for unity offset, returned: %v", val)
	}

	// Un equal length series should return error
	if _, err := NewPair(serie1, serie3).Compare(); err == nil {
		t.Errorf("Compare does not return error when uneqal length series are recieved")
	}

}

func TestWeightedPairs_Compare(t *testing.T) {
	//TODO: Define behavior when there are no overlapping args
	serie1 := NewSeries(maths.Arange(0., 2., 0.1))
	serie2 := NewSeries(maths.Arange(1., 3., 0.1))
	serie3 := NewSeries(maths.Arange(1., 3., 0.2))

	pair1 := NewPair(serie1, serie1)
	pair2 := NewPair(serie1, serie2)

	weights := []float64{1.0, 2.0}

	weightedPairs := WeightedPairs{[]Pair{pair1, pair2}, weights}

	// Compare equal series should return 0
	if val, _ := pair1.Compare(); val != 0.0 {
		t.Errorf("Compare does not return zero error for equal series")
	}

	// Compare equal series should return 0
	if val, _ := pair2.Compare(); !floats.EqualWithinAbs(val, 1.0, 1e-6) {
		t.Errorf("Compare does not return 1 error for unity offset, returned: %v", val)
	}

	// Compare equal series should return 0
	if val, _ := weightedPairs.Compare(); !floats.EqualWithinAbs(val, 2.0, 1e-6) {
		t.Errorf("Compare does not return 1 error for unity offset, returned: %v", val)
	}

	// Un equal length series should return error
	if _, err := NewPair(serie1, serie3).Compare(); err == nil {
		t.Errorf("Compare does not return error when uneqal length series are recieved")
	}

}

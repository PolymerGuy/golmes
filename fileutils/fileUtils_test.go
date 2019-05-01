package fileutils

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"

	"gonum.org/v1/gonum/floats"
)

func Test_substitute(t *testing.T) {
	data := "This one\n And these two\n Should be three"
	targets := []string{"one", "two", "three"}
	substitutes := []string{"1", "2", "3"}

	subData := string(replaceStrings([]byte(data), targets, substitutes))

	if !strings.EqualFold(subData, "This 1\n And these 2\n Should be 3") {
		t.Errorf("Substitution did not function as intended")
	}

}

func Test_floatsColumn(t *testing.T) {
	data := `fld1,fld2,fld3
1,2,3
4,5,6
7,8,9`

	mixedData := `fld1,fld2,fld3
1,2,3
4,5,6
q,w,e
7,8,9`

	columnInd := 1

	reader := csv.NewReader(bytes.NewBufferString(data))

	// Check if correct column is returned
	if column, _ := floatsColumn(reader, columnInd); !floats.Equal(column, []float64{2, 5, 8}) {
		t.Fatalf("Should return %v returned %v", []float64{2, 5, 8}, column)
	}
	mixedReader := csv.NewReader(bytes.NewBufferString(mixedData))

	// Check if correct column is returned
	if column, _ := floatsColumn(mixedReader, columnInd); !floats.Equal(column, []float64{2, 5, 8}) {
		t.Fatalf("Should return %v returned %v", []float64{2, 5, 8}, column)
	}

	reader = csv.NewReader(bytes.NewBufferString(data))

	// Check that an error is returned if columnInd is out of bounds
	if column, err := floatsColumn(reader, 5); err == nil {
		t.Fatalf("Should return an error, returned %v", column)
	}

}

func Test_findWordInd(t *testing.T) {
	data := []string{"This", "is", "a", "test"}

	if ind, err := findWordInd(data, "is"); ind != 1 || err != nil {
		t.Fatalf("Found wrong index")
	}

	if _, err := findWordInd(data, "wrong"); err == nil {
		t.Fatalf("Did not return error when key is not present")
	}
}

// Note that this test is done without tolerances
func Test_parseFloat(t *testing.T) {
	tol := 1e-9

	if val,err := parseFloat("1e9");!floats.EqualWithinAbs(val,float64(1e9),tol)||err!=nil{
		t.Fatalf("Did not correct the correct float. Should be %v got %v",float64(1e9),val)
	}

	if val,err := parseFloat("-1e9");!floats.EqualWithinAbs(val,float64(-1e9),tol)||err!=nil{
		t.Fatalf("Did not correct the correct float. Should be %v got %v",float64(-1e9),val)
	}

	if val,err := parseFloat("-1.23456");!floats.EqualWithinAbs(val,float64(-1.23456),tol)||err!=nil{
		t.Fatalf("Did not correct the correct float. Should be %v got %v",float64(-1.23456),val)
	}
}

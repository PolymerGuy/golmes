package data

import (
	"errors"
	"fmt"
	"github.com/PolymerGuy/golmes/fileutils"
	"github.com/PolymerGuy/golmes/maths"
	"github.com/pkelchte/spline"
	"gonum.org/v1/gonum/floats"
	"log"
	"math"
	"sort"
)




type DataReader interface {
	Read() []float64
}

type DataReaderWithArgs interface {
	Read() []float64
	ReadArgs() []float64
	ReadAt([]float64) []float64
}

type Comparator interface {
	Compare(dataset DataReaderWithArgs) (float64, error)
}

type Serie []float64
type SerieWithArgs struct {
	arguments      DataReader
	functionValues DataReader
}

type Pair struct {
	data1 DataReader
	data2 DataReader
}
type PairWithArgs struct {
	data1      DataReaderWithArgs
	data2      DataReaderWithArgs
	commonArgs DataReader
}

type DataFromFile struct {
	fileName string
	key      string
}

type WeightedPairs struct {
	objectiveFunctions []Pair
	weights            []float64
}

func NewSeries(data []float64) Serie {
	return data
}

func NewSeriesWithArgs(arguments DataReader, functionValues DataReader) SerieWithArgs {
	return SerieWithArgs{functionValues: functionValues, arguments: arguments}
}

func NewPair(data1 DataReader, data2 DataReader) Pair {
	return Pair{data1, data2}
}

func NewPairWithArgs(data1 DataReaderWithArgs, data2 DataReaderWithArgs, commonArgs DataReader) PairWithArgs {
	return PairWithArgs{data1, data2, commonArgs}
}

func NewSeriesFromFile(fileName string, key string) DataFromFile {
	return DataFromFile{fileName, key}
}

func (pair Pair) GetFields() []DataReader {
	return []DataReader{pair.data1, pair.data2}
}
func (pair PairWithArgs) GetFields() []DataReaderWithArgs {
	return []DataReaderWithArgs{pair.data1, pair.data2}
}

func (series DataFromFile) Read() []float64 {
	return fileutils.GetKeyFromCSVFile(series.fileName, series.key)
}

func (series Serie) Read() []float64 {
	return series
}

func (series SerieWithArgs) Read() []float64 {
	return series.functionValues.Read()
}
func (series SerieWithArgs) ReadAt(xs []float64) []float64 {

	if !sort.Float64sAreSorted(series.arguments.Read()) {
		fmt.Println("Args arguments: ", series.arguments.Read())
		log.Fatal("arguments are not sorted")
	}

	//interpolator := spline.NewCubic(series.arguments.Read(), series.functionValues.Read())
	// Check if it works properly...
	// Check if xs is witin the args of series
	if !floats.Equal(maths.EnclosedWithin(xs, series.arguments.Read()), xs) {
		log.Fatalln("Tried to evaluate interpolator outside its domain")
	}

	s := spline.Spline{}
	if len(series.arguments.Read()) != len(series.functionValues.Read()) {
		log.Fatalln("Cannot initialize interpolator, arguments and values are of unequal length")
	}

	s.Set_points(series.arguments.Read(), series.functionValues.Read(), true)

	results := []float64{}
	for _, x := range xs {
		results = append(results, s.Operate(x))
	}

	return results
}
func (series SerieWithArgs) ReadArgs() []float64 {
	return series.arguments.Read()
}

func (data Pair) Compare() (float64, error) {
	data1 := data.data1.Read()
	data2 := data.data2.Read()

	rmsError, err := maths.RMSError(data1, data2)
	if err != nil {
		return math.Inf(1), errors.New(fmt.Sprintf("Comparison of data failed: %s", err))
	}
	return rmsError, nil
}

// Compare return the RMS error of two series, synced to the largest common overlapping range.
// If an commonArgs serie is provided, this is used to sync the series
// TODO: Extend Compare to return an error as well...
func (data PairWithArgs) Compare(datas2 DataReaderWithArgs) (float64, error) {
	args := []float64{}
	switch len(data.commonArgs.Read()) {
	case 0:
		args = maths.EnclosedWithin(data.data1.ReadArgs(), datas2.ReadArgs())
	default:
		// The common arguments should be enclosed within the args of serie1 and serie2
		commonArguments := maths.EnclosedWithin(data.data1.ReadArgs(), datas2.ReadArgs())
		args = maths.EnclosedWithin(commonArguments, data.commonArgs.Read())

	}
	if floats.Equal(args, []float64{}) {
		return math.Inf(1), errors.New(fmt.Sprintln("No overlapping arguments for compare"))
	}

	data1 := data.data1.ReadAt(args)

	data2 := datas2.ReadAt(args)

	rmsError, err := maths.RMSError(data1, data2)
	if err != nil {
		return math.Inf(1), errors.New(fmt.Sprintf("Comparison of data failed: %s", err))
	}
	return rmsError, nil
}

func (objFuncs WeightedPairs) Compare() (float64, error) {
	var weightedFuncVals float64
	for i, objFunc := range objFuncs.objectiveFunctions {
		comp, err := objFunc.Compare()
		if err != nil {
			return math.Inf(1), errors.New(fmt.Sprintf("Comparison of data failed: %s", err))

		}
		weightedFuncVals = weightedFuncVals + objFuncs.weights[i]*comp

	}

	return weightedFuncVals, nil
}

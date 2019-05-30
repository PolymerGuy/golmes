package yamlparser

import (
	"fmt"
	"github.com/PolymerGuy/golmes/appwrapper"
	"log"
	"strconv"

	"github.com/PolymerGuy/golmes/data"
	"gonum.org/v1/gonum/optimize"
	"gopkg.in/yaml.v2"
	"strings"
)

// Note: struct fields must be public in order for unmarshal to
// correctly populate the data.
type DataComparators struct {
	Type           string
	Referencefile  string //
	Currentfile    string //
	Commonargsfile string
	Keywords       []string
}

type AbqSetting struct {
	Path          string
	Inputfile     string
	Keywords      []string
	InitialParams []float64 `initial_parameters`
}

type SolverSet struct {
	Method       string
	Threshold    float64
	Evaluations  int
	CoarseSearch CoarseSearch `coarsesearch`
}

type CoarseSearch struct {
	Seed       []string
	Limits     []string
	Scales     []int
	Refinement []float64
}

type YamlData struct {
	Comparators    []DataComparators `DataComparators`
	AbqSettings    AbqSetting        `Abaqus_settings`
	SolverSettings SolverSet         `Solver_settings`
	WorkDir        string            `Work_directory`
}

func Parse(yamlData []byte) YamlData {
	var template YamlData

	err := yaml.Unmarshal(yamlData, &template)

	if err != nil {
		log.Println("Could not marshall the file")
		log.Fatalf("Unable to unmarshal data: %v", err)
	}

	return template
}

func (yamlData YamlData) NewComparator() []data.PairWithArgs {

	comparators := []data.PairWithArgs{}

	for _, templ := range yamlData.Comparators {
		comparators = append(comparators, MakeComparator(templ))
	}
	return comparators
}
func (yamlData YamlData) NewCostFunction() appwrapper.CostFunction {

	costFunction := appwrapper.CostFunction{
		ExecPath:              yamlData.AbqSettings.Path,
		InputFileTemplateName: yamlData.AbqSettings.Inputfile,
		WorkDirectory:         yamlData.WorkDir,
		ArgKeywords:           yamlData.AbqSettings.Keywords,
		InitialParameters:     yamlData.AbqSettings.InitialParams,
		ResKeywords:           yamlData.Comparators[0].Keywords,
	}

	return costFunction
}
func (data YamlData) NewOptimizerSettings() optimize.Settings {
	var settings = optimize.Settings{
		FuncEvaluations: data.SolverSettings.Evaluations,
		Converger:       &optimize.FunctionConverge{Absolute: data.SolverSettings.Threshold, Iterations: 10},
	}

	return settings
}

func MakeComparator(settings DataComparators) data.PairWithArgs {

	refStrain := data.NewSeriesFromFile(settings.Referencefile, settings.Keywords[0])

	refStress := data.NewSeriesFromFile(settings.Referencefile, settings.Keywords[1])

	refData := data.NewSeriesWithArgs(refStrain, refStress)

	curStrain := data.NewSeriesFromFile(settings.Currentfile, settings.Keywords[0])

	curStress := data.NewSeriesFromFile(settings.Currentfile, settings.Keywords[1])

	curData := data.NewSeriesWithArgs(curStrain, curStress)

	// TODO: Remove refData as default args!
	return data.NewPairWithArgs(refData, curData, refStrain)

}

func (data YamlData) NewOptimizerMethod() optimize.Method {

	method := lowerCaseWithoutSeps(data.SolverSettings.Method)
	fmt.Println("Method:",method)
	fmt.Println("Method:",data.SolverSettings.Method)
	switch method {
	case "neldermead":
		return &optimize.NelderMead{}
	case "gradient":
		return &optimize.GradientDescent{}
	default:
		log.Println("Using default fine search method: Nelder Mead")
		return &optimize.NelderMead{}

	}
}

type CoarseSearchSettings struct {
	Seeds      []int
	Bounds     []float64
	NPts       int
	Refinement []float64
}

func (yamlData YamlData) NewCoarseSearch() CoarseSearchSettings {
	nPts := 1
	for _, seed := range stringSliceToIntSlice(yamlData.SolverSettings.CoarseSearch.Seed) {
		nPts *= seed
	}

	seeds := stringSliceToIntSlice(yamlData.SolverSettings.CoarseSearch.Seed)

	bounds := stringSliceToFloatSlice(yamlData.SolverSettings.CoarseSearch.Limits)


	if len(seeds) * 2 != len(bounds){
		log.Fatal("The number of seeds should match the number of bounds.")
	}


	refinement := yamlData.SolverSettings.CoarseSearch.Refinement

	coarseSearch := CoarseSearchSettings{
		Seeds:      seeds,
		Bounds:     bounds,
		NPts:       nPts,
		Refinement: refinement,
	}
	return coarseSearch

}

// lowerCaseWithoutSeps removes whitespaces and turns all chars to lowercase
func lowerCaseWithoutSeps(s string) string {
	s = strings.Replace(s, "-", "", -1)
	s = strings.Replace(s, " ", "", -1)
	return strings.ToLower(s)

}

func stringSliceToIntSlice(t []string) []int {
	var t2 = []int{}

	for _, i := range t {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		t2 = append(t2, j)
	}
	return t2
}

func stringSliceToFloatSlice(t []string) []float64 {
	var t2 = []float64{}

	for _, i := range t {
		for _, k := range strings.Fields(i) {
			j, err := strconv.ParseFloat(k, 64)
			if err != nil {
				panic(err)
			}
			t2 = append(t2, j)
		}
	}
	return t2
}

package cmd

import (
	"github.com/PolymerGuy/golmes/minimize"
	"github.com/PolymerGuy/golmes/output"
	"github.com/PolymerGuy/golmes/yamlparser"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
)

var Optimize = cli.Command{
	Name:        "run",
	Usage:       "Run optimization",
	Description: `Bla bla bla`,
	Action:      optimize,
}

func optimize(c *cli.Context) {
	args := c.Args()
	if len(args) == 0 {
		log.Panic("No input was received")
		os.Exit(1)
	} else {
		yamlFileName := args.Get(0)

		file, err := ioutil.ReadFile(yamlFileName)
		if err != nil {
			log.Println(err)
		}
		coarseSeach, optJob := optJobFromYAML(file)
		initialParameters, _ := minimize.CoarseSearch(optJob, coarseSeach)
		log.Println("Coarse search gave:")
		output.PrettyPrint(initialParameters)

		optJob.InitialParameters = initialParameters.X

		res, err := minimize.FindFunctionMinima(optJob)
		if err != nil {
			log.Panic(err)
		}
		log.Println("Fine search gave:")
		output.PrettyPrint(res)

	}
}

func optJobFromYAML(yamlFile []byte) (yamlparser.CoarseSearchSettings, minimize.OptimizationJob) {
	parser := yamlparser.Parse(yamlFile)

	//TODO: Remove list of comparators. This should be impossible...
	comparator := parser.NewComparator()[0]
	costFunction := parser.NewCostFunction()
	method := parser.NewOptimizerMethod()
	costFunction.Comparator = comparator
	settings := parser.NewOptimizerSettings()
	coarseSearch := parser.NewCoarseSearch()

	optJob := minimize.OptimizationJob{
		InitialParameters: costFunction.InitialParameters,
		Method:            method,
		CostFunc:          costFunction,
		Settings:          settings,
	}
	return coarseSearch, optJob
}


func optJobsFromYAML(yamlFile []byte) (yamlparser.CoarseSearchSettings, minimize.OptimizationJob) {
	parser := yamlparser.Parse(yamlFile)


	jobs := []minimize.OptimizationJob
	for ,_ := range parser.Comparators{

		comparator := parser.NewComparator()[0]
		costFunction := parser.NewCostFunction()
		method := parser.NewOptimizerMethod()
		costFunction.Comparator = comparator
		settings := parser.NewOptimizerSettings()
		coarseSearch := parser.NewCoarseSearch()

		optJob := minimize.OptimizationJob{
		InitialParameters: costFunction.InitialParameters,
		Method:            method,
		CostFunc:          costFunction,
		Settings:          settings,
	}
		jobs = append(jobs,optJob)
	}
	return jobs
}

package cmd

import (
	"fmt"
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
			log.Fatal(err)
		}
		coarseSeach, optJob, err := optJobFromYAML(file)


		if err != nil {
			log.Fatal(err)
		}

		if coarseSeach != nil {
			for _, task := range coarseSeach {
				// Do a response surface naive search
				coarseResults, err := minimize.CoarseSearchSurf(optJob, task)
				if err != nil {
					log.Println(err)
				}

				// Use the results from the coarse search as initial parameters
				optJob.InitialParameters = coarseResults
				log.Println("Snip here...")
			}
		}

		// Do a fine search
		res, err := minimize.FindFunctionMinima(optJob)
		if err != nil {
			log.Panic(err)
		}
		log.Println("Fine search gave:")
		output.PrettyPrint(res)

	}
}

func optJobFromYAML(yamlFile []byte) ([]yamlparser.CoarseSearchSettings, minimize.OptimizationJob, error) {
	parser := yamlparser.Parse(yamlFile)

	fmt.Println(parser.AbqSettings)
	fmt.Println(parser.SolverSettings)

	coarseSearchs := []yamlparser.CoarseSearchSettings{}

	//TODO: Remove list of comparators. This should be impossible...
	comparator := parser.NewComparator()[0]
	costFunction := parser.NewCostFunction()
	method := parser.NewOptimizerMethod()
	costFunction.Comparator = comparator
	settings := parser.NewOptimizerSettings()
	coarseSearch := parser.NewCoarseSearch()

	coarseSearchs = append(coarseSearchs, coarseSearch)

	optJob := minimize.OptimizationJob{
		InitialParameters: costFunction.InitialParameters,
		Method:            method,
		CostFunc:          costFunction,
		Settings:          settings,
	}
	return coarseSearchs, optJob, nil
}

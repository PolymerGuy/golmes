package appwrapper

import (
	"github.com/PolymerGuy/golmes/data"
	"github.com/PolymerGuy/golmes/fileutils"
	"github.com/lithammer/shortuuid"
	"log"
	"os"
	"os/exec"
	"strings"
)

type CostFunction struct {
	ExecPath              string
	InputFileTemplateName string
	WorkDirectory         string
	ArgKeywords           []string
	ResKeywords           []string
	NParameters           int
	InitialParameters     []float64
	Comparator            data.Comparator
}

func (app CostFunction) Eval(args []float64) float64 {
	// Implements the costFunction interface
	// Generate input file with args

	if _, err := os.Stat(app.WorkDirectory); os.IsNotExist(err) {
		err = os.MkdirAll(app.WorkDirectory, 0755)
		if err != nil {
			panic(err)
		}
	}

	id := shortuuid.New()
	curWorkDir := app.WorkDirectory + id

	err := os.MkdirAll(curWorkDir, 0755)
	if err != nil {
		log.Panic(err)
	}

	inpurFileName, err := fileutils.MakeInputFile(app.InputFileTemplateName, curWorkDir, args, app.ArgKeywords)
	// Run Abaqus

	err = runApp(app.ExecPath, inpurFileName)
	if err != nil {
		log.Fatal("Abaqus failed!")
	}

	resFileName := strings.TrimSuffix(inpurFileName, ".txt") + "_res.txt"

	curStrain := data.NewSeriesFromFile(resFileName, app.ResKeywords[0])

	curStress := data.NewSeriesFromFile(resFileName, app.ResKeywords[1])

	curData := data.NewSeriesWithArgs(curStrain, curStress)

	res, err := app.Comparator.Compare(curData)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(id, args, res)

	return res
}

func runApp(appPath string, jobName string) error {
	run := exec.Command(appPath, jobName)
	err := run.Run()
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

package appwrapper

import (
	"errors"
	"github.com/PolymerGuy/golmes/data"
	"github.com/PolymerGuy/golmes/fileutils"
	"github.com/lithammer/shortuuid"
	"log"
	"os"
	"os/exec"
	"strings"
)

type CostFunction interface {
	Eval(args []float64) float64
}


type AppWrapper struct {
	ExecPath              string
	InputFileTemplateName string
	WorkDirectory         string
	ArgKeywords           []string
	ResKeywords           []string
	NParameters           int
	InitialParameters     []float64
	Comparator            data.Comparator
}

func (app AppWrapper) Eval(args []float64) float64 {
	// Implements the costFunction interface
	// Generate input file with args

	// Check if the main work directory is there, if not, make it
	if _, err := os.Stat(app.WorkDirectory); os.IsNotExist(err) {
		err = os.MkdirAll(app.WorkDirectory, 0755)
		if err != nil {
			log.Fatal("Could not make main work directory:",err)

		}
	}

	// Make unique work sub directory
	uniqueWorkDir,err := makeUniqueDir(app.WorkDirectory)
	if err != nil{
		log.Fatal("Could not make main unique directory:",err)

	}

	// Make an input file with all keywords set
	inpurFileName, err := fileutils.MakeInputFile(app.InputFileTemplateName, uniqueWorkDir, args, app.ArgKeywords)
	if err != nil {
		log.Fatal("Could not make input file from template ",err)
	}

	// Run application
	err = runApp(app.ExecPath, inpurFileName)
	if err != nil {
		log.Fatal("Running the application failed: ", err)
	}

	resFileName := strings.TrimSuffix(inpurFileName, ".txt") + "_res.txt"

	curStrain := data.NewSeriesFromFile(resFileName, app.ResKeywords[0])

	curStress := data.NewSeriesFromFile(resFileName, app.ResKeywords[1])

	curData := data.NewSeriesWithArgs(curStrain, curStress)

	res, err := app.Comparator.Compare(curData)
	if err != nil {
		log.Fatal("Could not compare data ", err)
	}
	log.Println(uniqueWorkDir, args, res)

	return res
}

func makeUniqueDir(path string) (string, error) {
	// Make a unique work sub directory
	id := shortuuid.New()
	curWorkDir := path + id
	err := os.MkdirAll(curWorkDir, 0755)
	if err != nil {
		return "", errors.New("Could not make the local work directory: "+ path+ err.Error())
	}
	return curWorkDir,nil

}

func runApp(appPath string, jobName string) error {
	run := exec.Command(appPath, jobName)
	err := run.Run()
	if err != nil {
		return err
	}
	return nil
}

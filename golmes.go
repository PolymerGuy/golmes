package main

import (
	"github.com/PolymerGuy/golmes/cmd"
	"github.com/urfave/cli"
	"log"
	"os"
)

// This is the main Golmes application
// It sets up logging and runs the commands
func main() {

	logFile := makeFile("runtime.log")
	log.SetOutput(logFile)
	defer logFile.Close()

	app := cli.NewApp()
	app.Name = "Golmes"
	app.Usage = "Multivariate optimization"

	app.Commands = []cli.Command{
		cmd.Optimize,
		cmd.Web,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func makeFile(logfile string) *os.File {
	f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return f

}

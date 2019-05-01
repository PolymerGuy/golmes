package cmd

import (
	"github.com/urfave/cli"
	"log"
)

var Web = cli.Command{
	Name:        "web",
	Usage:       "Start web interface",
	Description: `This is a web interface where you can monitor the optimization process`,
	Action:      webUI,
}

func webUI(c *cli.Context) {
	log.Println("The web ui is not implemented yet")
}

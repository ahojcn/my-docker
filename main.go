package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"mydocker/command"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "mydocker"
	app.Usage = "implementation of mydocker"

	app.Commands = []cli.Command{
		command.RunCommand,
		command.InitCommand,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln("Run() error", err)
	}
}

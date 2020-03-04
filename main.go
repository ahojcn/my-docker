package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"mydocker/command"
	_ "mydocker/nsenter"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "mydocker"
	app.Usage = "implementation of mydocker"

	app.Commands = []cli.Command{
		command.RunCommand,
		command.InitCommand,
		command.CommitCommand,
		command.ListCommand,
		command.LogCommand,
		command.ExecCommand,
		command.StopCommand,
		// command.RemoveCommand,
		command.CopyCommand,
		command.NetworkCommand,
	}

	app.Before = func(context *cli.Context) error {
		// Log as JSON instead of the default ASCII formatter.
		log.SetFormatter(&log.TextFormatter{})
		log.SetOutput(os.Stdout)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln("Run() error", err)
	}
}
package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"time"
)

// Command
func main() {
	log.Infoln("main() start.")

	app := cli.NewApp()
	app.Name = "appName..."
	app.Usage = "appUsage..."

	runCmd := cli.Command{
		Name: "run",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "it",
				Usage: "enable tty",
			},
		},
		Action: func(c *cli.Context) error {
			log.Infoln("run cmd args:", c.Args())
			log.Infoln("run cmd tty:", c.Bool("it"))
			for i := 0; i < 5; i++ {
				fmt.Println("run cmd sleep", i)
				time.Sleep(time.Second * 1)
			}
			return nil
		},
	}

	app.Commands = []*cli.Command{
		&runCmd,
	}

	app.Action = func(c *cli.Context) error {
		log.Infoln("main() args", c.Args())
		for i := 0; i < 5; i++ {
			log.Infoln("main() sleep", i)
			time.Sleep(time.Second * 1)
		}
		return nil
	}

	log.Infoln("Run()")
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln("app.Run() error", err)
	}
	log.Infoln("main() end.")
}

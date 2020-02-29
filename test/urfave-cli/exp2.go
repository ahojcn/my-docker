package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"time"
)

// 加入 flags
func main() {
	log.Infoln("main() start.")

	app := cli.NewApp()
	app.Name = "appName"
	app.Usage = "appUsage"

	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "flag",
			Usage: "enable tty",
		},
		&cli.StringFlag{
			Name:  "lang",
			Value: "english",
		},
	}

	app.Action = func(c *cli.Context) error {
		log.Infoln("args:", c.Args())
		log.Infoln("flag:", c.Bool("flag"))
		log.Infoln("lang:", c.String("lang"))
		for i := 0; i < 5; i++ {
			log.Infoln("sleep", i)
			time.Sleep(time.Second * 1)
		}

		return nil
	}

	log.Infoln("Run()")
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal("Run() error", err)
	}
	log.Infoln("main() end.")
}

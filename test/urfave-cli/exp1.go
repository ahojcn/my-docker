package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"time"
)

// 基础用法
func main() {
	log.Infoln("main() start.")

	app := cli.NewApp()
	app.Name = "app名字"
	app.Usage = "使用方式 ..."

	app.Action = func(c *cli.Context) error {
		log.Infoln("args:", c.Args())
		for i := 0; i < 5; i++ {
			log.Infoln("sleeping:", i)
			time.Sleep(time.Second * 1)
		}

		return nil
	}

	log.Infoln("Run(), os.Args:", os.Args)
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln("Run() error,", err)
	}

	log.Infoln("main() end.")
}

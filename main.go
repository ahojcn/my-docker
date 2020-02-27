package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

const usage = `mydocker is a simple container runtime implementation.
			   The purpose of this project is to learn how docker works and how to write a docker by ourselves
			   Enjoy it, just for fun.`

func main() {
	app := cli.NewApp()
	app.Name = "mydocker"
	app.Usage = usage

	app.Commands = []*cli.Command{
		&initCommand,
		&runCommand,
	}

	// 初始化 logrus 日志配置
	app.Before = func(context *cli.Context) error {
		// Log as JSON instead of the default ASCII formatter.
		//log.SetFormatter(&log.JSONFormatter{})  // 使用 json 格式的 log 信息
		log.Infoln("初始化 logrus 配置。")
		log.SetFormatter(&log.TextFormatter{}) // 使用 text 格式的 log 信息
		log.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

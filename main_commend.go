package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"mydocker/cgroups/subsystems"
	"mydocker/container"
)

/*
定义了 runCommand 的 Flags
作用类似于运行命令时候使用 -- 指定参数
*/
var runCommand = cli.Command{
	Name: "run",
	Usage: `Create a container with namespace and cgroups limit
			mydocker run -ti [command]`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
		&cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
		&cli.StringFlag{
			Name:  "cpushare",
			Usage: "cpushare limit",
		},
		&cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit",
		},
	},
	/*
		这里是 Run() 命令执行的真正函数
		1. 判断参数是否包含 command
		2. 获取用户指定的 command
		3. 调用 Run() 去准备启动容器
	*/
	Action: func(context *cli.Context) error {
		log.Infoln("解析参数。")
		if context.Args().Len() < 1 {
			return fmt.Errorf("Missing container command")
		}
		var cmdArray []string

		for i := 0; i < context.Args().Len(); i++ {
			cmdArray = append(cmdArray, context.Args().Get(i))
		}

		tty := context.Bool("ti")
		resConf := &subsystems.ResourceConfig{
			MemoryLimit: context.String("m"),
			CpuSet:      context.String("cpuset"),
			CpuShare:    context.String("cpushare"),
		}

		Run(tty, cmdArray, resConf)
		return nil
	},
}

/*
定义了 initCommand 的具体操作
此操作为内部方法，禁止外部调用
*/
//var initCommand = cli.Command{
//	Name:  "init",
//	Usage: "Init container process run user's process in container. Do not call it outside",
//	/*
//		1. 获取传递过来的 command 参数
//		2. 执行容器初始化操作
//	*/
//	Action: func(context *cli.Context) error {
//		log.Infof("init come on")
//		cmd := context.Args().Get(0)
//		log.Infof("command %s", cmd)
//		err := container.RunContainerInitProcess()
//		return err
//	},
//}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",
	Action: func(context *cli.Context) error {
		log.Infof("容器初始化操作。")
		err := container.RunContainerInitProcess()
		return err
	},
}

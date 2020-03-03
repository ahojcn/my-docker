package command

import (
	"github.com/urfave/cli"
	"mydocker/cgroups"
	"mydocker/cgroups/subsystems"
)

var RunCommand = cli.Command{
	Name: "run",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "it",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name:  "m",
			Usage: "limit memory usage",
		},
		cli.StringFlag{
			Name:  "r",
			Usage: "set root path",
		},
		cli.StringSliceFlag{
			Name:  "v",
			Usage: "enable volume",
		},
		cli.BoolFlag{
			Name:  "d",
			Usage: "enable detach",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "container name",
		},
	},

	Action: func(c *cli.Context) error {
		tty := c.Bool("it")
		memory := c.String("m")
		rootPath := c.String("r")
		volumes := c.StringSlice("v")
		detach := c.Bool("d")
		containerName := c.String("name")
		command := c.Args().Get(0)

		/**
		如果用户设置了该参数，此时需要把该限制对应的 subsystem 加入到 CgroupManager 中的属性 SubsystemIns 数组中对其限制
		并把 memory 值放到 ResourceConfig 中的 MemoryLimit 属性
		*/
		res := subsystems.ResourceConfig{MemoryLimit: memory}
		cg := cgroups.CgroupManger{
			Resource:      &res,
			SubsystemsIns: make([]subsystems.Subsystem, 0),
		}
		if memory != "" {
			cg.SubsystemsIns = append(cg.SubsystemsIns, &subsystems.MemorySubsystem{})
		}

		if detach { // 如果有参数 d，就把 tty 设置为 false 后台运行
			tty = false
		}

		Run(command, tty, &cg, rootPath, volumes, containerName)

		return nil
	},
}

var CommitCommand = cli.Command{
	Name: "commit",
	Action: func(c *cli.Context) error {
		imageName := c.Args().Get(0)
		Commit(imageName)
		return nil
	},
}

var ListCommand = cli.Command{
	Name: "ps",
	Action: func(c *cli.Context) error {
		List()
		return nil
	},
}

var LogCommand = cli.Command{
	Name: "logs",
	Action: func(c *cli.Context) error {
		containerName := c.Args().Get(0)
		Logs(containerName)
		return nil
	},
}

var InitCommand = cli.Command{
	Name: "init",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "it",
			Usage: "enable tty",
		},
	},
	Action: func(c *cli.Context) error {
		command := c.Args().Get(0)
		Init(command)
		return nil
	},
}

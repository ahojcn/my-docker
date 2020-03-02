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
	},

	Action: func(c *cli.Context) error {
		tty := c.Bool("it")
		memory := c.String("m")
		rootPath := c.String("r")
		volumes := c.StringSlice("v")
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

		Run(command, tty, &cg, rootPath, volumes)

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

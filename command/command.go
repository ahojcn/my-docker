package command

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"mydocker/cgroups"
	"mydocker/cgroups/subsystems"
	"mydocker/network"
	"os"
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
		cli.StringSliceFlag{
			Name:  "e",
			Usage: "set env",
		},
		cli.StringFlag{
			Name:  "net",
			Usage: "container network",
		},
		cli.StringSliceFlag{
			Name:  "p",
			Usage: "port mapping",
		},
	},

	Action: func(c *cli.Context) error {
		tty := c.Bool("it")
		memory := c.String("m")
		rootPath := c.String("r")
		volumes := c.StringSlice("v")
		detach := c.Bool("d")
		containerName := c.String("name")
		envSlice := c.StringSlice("e")
		network := c.String("net")
		portMapping := c.StringSlice("p")

		imageName := c.Args().Get(0)
		command := c.Args().Get(1)

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

		Run(command, tty, &cg, rootPath, volumes, containerName, imageName, envSlice, network, portMapping)

		return nil
	},
}

var CommitCommand = cli.Command{
	Name: "commit",
	Action: func(c *cli.Context) error {
		// imageName := c.Args().Get(0)
		containerName := c.Args().Get(0)
		imageName := c.Args().Get(1)
		Commit(containerName, imageName)
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

var ExecCommand = cli.Command{
	Name: "exec",
	Action: func(c *cli.Context) error {
		if os.Getenv("mydocker_pid") != "" {
			log.Infoln("pid callback pid:", os.Getgid())
			return nil
		}
		containerName := c.Args().Get(0)
		command := c.Args().Get(1)
		log.Printf("container name : %s, command: %s", containerName, command)
		Exec(containerName, command)
		return nil
	},
}

var StopCommand = cli.Command{
	Name:  "stop",
	Usage: "stop a container",
	Action: func(c *cli.Context) error {
		containerName := c.Args().Get(0)
		Stop(containerName)
		return nil
	},
}

var RemoveCommand = cli.Command{
	Name:  "rm",
	Usage: "remove a stopped container",
	Action: func(c *cli.Context) error {
		containerName := c.Args().Get(0)
		Remove(containerName)
		return nil
	},
}

var CopyCommand = cli.Command{
	Name:  "cp",
	Usage: "copy files",
	Action: func(c *cli.Context) error {
		src := c.Args().Get(0)
		dst := c.Args().Get(1)
		log.Infof("src: %s, dst: %s", src, dst)
		Copy(src, dst)
		return nil
	},
}

var NetworkCommand = cli.Command{
	Name:  "network",
	Usage: "container network commands",
	Subcommands: []cli.Command{
		{
			Name:  "create",
			Usage: "create a container network",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "driver",
					Usage: "network driver",
				},
				cli.StringFlag{
					Name:  "subnet",
					Usage: "subnet cidr",
				},
			},
			Action: func(context *cli.Context) error {
				if len(context.Args()) < 1 {
					return fmt.Errorf("Missing network name")
				}
				network.Init()
				err := network.CreateNetwork(context.String("driver"), context.String("subnet"), context.Args()[0])
				if err != nil {
					return fmt.Errorf("create network error: %+v", err)
				}
				return nil
			},
		},
		{
			Name:  "list",
			Usage: "list container network",
			Action: func(context *cli.Context) error {
				network.Init()
				network.ListNetwork()
				return nil
			},
		},
		{
			Name:  "remove",
			Usage: "remove container network",
			Action: func(context *cli.Context) error {
				if len(context.Args()) < 1 {
					return fmt.Errorf("Missing network name")
				}
				network.Init()
				err := network.DeleteNetwork(context.Args()[0])
				if err != nil {
					return fmt.Errorf("remove network error: %+v", err)
				}
				return nil
			},
		},
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

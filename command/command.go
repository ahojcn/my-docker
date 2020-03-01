package command

import "github.com/urfave/cli"

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
	},

	Action: func(c *cli.Context) error {
		tty := c.Bool("it")
		memory := c.String("m")
		command := c.Args().Get(0)

		Run(command, tty, memory)

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

package command

import "github.com/urfave/cli"

var RunCommand = cli.Command{
	Name: "run",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "it",
			Usage: "enable tty",
		},
	},

	Action: func(c *cli.Context) error {
		tty := c.Bool("it")
		command := c.Args().Get(0)
		Run(command, tty)
		return nil
	},
}

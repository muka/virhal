package main

import (
	"fmt"
	"os"

	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()

	app.Name = "virhal"
	app.Usage = "Virtualized Hardware Abstraction Layer CLI tools"
	app.Version = "1.0.0-alpha1"
	flags := []cli.Flag{
		cli.IntFlag{
			Name:   "debug, d",
			Usage:  "debug output level",
			EnvVar: "DEBUG",
		},
	}
	app.Flags = flags

	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{},
			Usage:   "run a container based on label selection",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "filter, f",
					Usage: "filter for selection",
					// EnvVar: "FILTER",
				},
			},
			Action: func(c *cli.Context) error {
				fmt.Println("added task: ", c.String("flag"))
				return nil
			},
		},
		// {
		// 	Name:    "template",
		// 	Aliases: []string{"t"},
		// 	Usage:   "options for task templates",
		// 	Subcommands: []cli.Command{
		// 		{
		// 			Name:  "add",
		// 			Usage: "add a new template",
		// 			Action: func(c *cli.Context) error {
		// 				fmt.Println("new task template: ", c.Args().First())
		// 				return nil
		// 			},
		// 		},
		// 		{
		// 			Name:  "remove",
		// 			Usage: "remove an existing template",
		// 			Action: func(c *cli.Context) error {
		// 				fmt.Println("removed task template: ", c.Args().First())
		// 				return nil
		// 			},
		// 		},
		// 	},
		// },
	}

	app.Run(os.Args)
}

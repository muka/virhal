package main

import (
	"errors"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/muka/virhal/project"
	"github.com/muka/virhal/project/options"
	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()

	app.Name = "virhal"
	app.Usage = "Virtualized Hardware Abstraction Layer CLI tools"
	app.Version = "1.0.0-alpha1"

	flags := []cli.Flag{
		cli.BoolFlag{
			Name:   "debug, d",
			Usage:  "debug output level",
			EnvVar: "DEBUG",
		},
		cli.StringFlag{
			Name:   "file, f",
			Usage:  "configuration file",
			EnvVar: "CONFIG_FILE",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "start a project",
			Flags: flags,
			Action: func(c *cli.Context) error {

				if c.Bool("debug") {
					log.SetLevel(log.DebugLevel)
				}

				filepath := c.String("file")
				if filepath == "" {
					return errors.New("No file specified")
				}

				p, err := project.NewProjectFromFile(filepath)
				if err != nil {
					return err
				}

				return p.Start(options.Start{})
			},
		},
		{
			Name:    "status",
			Aliases: []string{"ps"},
			Usage:   "show status of a project",
			Flags:   flags,
			Action: func(c *cli.Context) error {

				if c.Bool("debug") {
					log.SetLevel(log.DebugLevel)
				}

				filepath := c.String("file")
				if filepath == "" {
					return errors.New("No file specified")
				}

				p, err := project.NewProjectFromFile(filepath)
				if err != nil {
					return err
				}

				return p.Status(options.Status{})
			},
		},
	}

	app.Run(os.Args)
}

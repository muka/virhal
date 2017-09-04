package main

import (
	"errors"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/muka/virhal/agent"
	"github.com/muka/virhal/docker"
	"github.com/muka/virhal/project"
	"github.com/muka/virhal/project/options"
	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()

	app.Name = "virhal"
	app.Usage = "Virtualized Hardware Abstraction Layer"
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
			Usage:   "Show status of a project",
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
		{
			Name:  "agent",
			Usage: "start the agent service",
			Flags: append(flags,
				cli.StringFlag{
					Name:   "port, p",
					Usage:  "port to listen for ",
					EnvVar: "GRPC_PORT",
					Value:  ":50051",
				},
			),
			Action: func(c *cli.Context) error {

				debug := c.Bool("debug")
				if debug {
					log.SetLevel(log.DebugLevel)
				}

				log.Info("Agent start")

				agent.RunServer(c.String("port"))

				err := docker.WatchEvents()
				if err != nil {
					return err
				}

				log.Info("Agent exit")
				return nil
			},
		},
	}

	app.Run(os.Args)
}

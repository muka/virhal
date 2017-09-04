package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/muka/virhal/agent"
	"github.com/muka/virhal/docker"
	"github.com/muka/virhal/project"
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
				p, err := project.NewProjectFromFile(filepath)
				if err != nil {
					log.Errorf("Failed to load project: %s", err.Error())
					return err
				}

				return agent.Start(p)
			},
		},
		{
			Name:  "stop",
			Usage: "stop a project",
			Flags: flags,
			Action: func(c *cli.Context) error {

				if c.Bool("debug") {
					log.SetLevel(log.DebugLevel)
				}

				filepath := c.String("file")
				p, err := project.NewProjectFromFile(filepath)
				if err != nil {
					log.Errorf("Failed to load project: %s", err.Error())
					return err
				}

				return agent.Stop(p)
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
				p, err := project.NewProjectFromFile(filepath)
				if err != nil {
					log.Errorf("Failed to load project: %s", err.Error())
					return err
				}

				res, err := agent.Status(p)

				if res != "" {
					log.Println(res)
					return nil
				}

				return err
			},
		},
		{
			Name:  "agent",
			Usage: "start the agent service",
			Flags: flags,
			Action: func(c *cli.Context) error {

				debug := c.Bool("debug")
				if debug {
					log.SetLevel(log.DebugLevel)
				}

				go func() {
					err := agent.RunServer()
					if err != nil {
						log.Error("Failed to start grpc server")
						panic(err)
					}
				}()

				go func() {
					err := docker.WatchEvents()
					if err != nil {
						log.Error("Failed to listen for events")
						panic(err)
					}
				}()

				log.Info("Agent started, waiting for events..")
				select {}
			},
		},
	}

	app.Run(os.Args)
}

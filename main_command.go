package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"writeDocker/container"
)

var (
	runCommand = cli.Command{
		Name: "run",
		Usage: `Create a container with namespace and cgroups limit
				mydocker run -ti [command]`,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "ti",
				Usage: "enable tty",
			},
		},
		// real action of run command.
		Action: func(context *cli.Context) error {
			// check params
			if len(context.Args()) < 1 {
				return fmt.Errorf("missing container command")
			}
			// get target command
			cmd := context.Args().Get(0)
			tty := context.Bool("ti")
			// use Run function to start container
			Run(tty, cmd)
			return nil
		},
	}

	initCommand = cli.Command{
		Name: "init",
		Usage: "Init container process run user's process in container." +
			"You can not call it outside",
		// real action of init
		Action: func(context *cli.Context) error {
			logrus.Infof("init come on")
			cmd := context.Args().Get(0)
			logrus.Infof("command %s", cmd)
			err := container.RunContainerInitProcess(cmd, nil)
			return err
		},
	}
)

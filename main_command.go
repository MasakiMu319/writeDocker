package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"writeDocker/cgroups/subsystems"
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
			cli.StringFlag{
				Name:  "m",
				Usage: "memory limit",
			},
			cli.StringFlag{
				Name:  "cpushare",
				Usage: "cpushare limit",
			},
			cli.StringFlag{
				Name:  "cpuset",
				Usage: "cpuset limit",
			},
		},
		// real action of run command.
		Action: func(context *cli.Context) error {
			// check params
			if len(context.Args()) < 1 {
				return fmt.Errorf("missing container command")
			}
			// get target command
			var cmdArray []string
			for _, arg := range context.Args() {
				cmdArray = append(cmdArray, arg)
			}
			tty := context.Bool("ti")
			// use Run function to start container
			resConf := &subsystems.ResourceConfig{
				MemoryLimit: context.String("m"),
				CpuSet:      context.String("cpuset"),
				CpuShare:    context.String("cpushare"),
			}
			Run(tty, cmdArray, resConf)
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
			err := container.RunContainerInitProcess()
			return err
		},
	}
)

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
			cli.BoolFlag{
				Name:  "d",
				Usage: "detach container",
			},
			cli.StringFlag{
				Name:  "v",
				Usage: "volume",
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
			cli.StringFlag{
				Name:  "name",
				Usage: "container name",
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
			detach := context.Bool("d")
			volume := context.String("v")

			if tty && detach {
				return fmt.Errorf("ti and d paramer can not both provided")
			}
			// use Run function to start container
			resConf := &subsystems.ResourceConfig{
				MemoryLimit: context.String("m"),
				CpuSet:      context.String("cpuset"),
				CpuShare:    context.String("cpushare"),
			}
			logrus.Infof("Running in an interactive environment : %v", tty)
			containerName := context.String("name")
			Run(tty, cmdArray, resConf, volume, containerName)
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

	commitCommand = cli.Command{
		Name:  "commit",
		Usage: "commit a container into image",
		Action: func(context *cli.Context) error {
			if len(context.Args()) < 1 {
				return fmt.Errorf("Missing container name")
			}
			imageName := context.Args().Get(0)
			commitContainer(imageName)
			return nil
		},
	}

	listCommand = cli.Command{
		Name:  "ps",
		Usage: "list all the containers",
		Action: func(context *cli.Context) error {
			ListContainers()
			return nil
		},
	}
)

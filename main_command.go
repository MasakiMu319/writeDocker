package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"writeDocker/cgroups/subsystems"
	"writeDocker/container"
	"writeDocker/network"
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
			cli.StringSliceFlag{
				Name:  "e",
				Usage: "set environment",
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

			imageName := cmdArray[0]
			cmdArray = cmdArray[1:]

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
			envSlice := context.StringSlice("e")
			network := context.String("net")
			portmapping := context.StringSlice("p")
			Run(tty, cmdArray, resConf, volume, containerName, imageName, envSlice, network, portmapping)
			return nil
		},
	}

	initCommand = cli.Command{
		Name: "init",
		Usage: "Init container process run user's process in container." +
			"You can not use it directly",
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
			if len(context.Args()) < 2 {
				return fmt.Errorf("Missing container name and image name")
			}
			containerName := context.Args().Get(0)
			imageName := context.Args().Get(1)
			commitContainer(containerName, imageName)
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

	logCommand = cli.Command{
		Name:  "logs",
		Usage: "print logs of a container",
		Action: func(context *cli.Context) error {
			if len(context.Args()) < 1 {
				return fmt.Errorf("Please input your container name")
			}
			containerName := context.Args().Get(0)
			logContainer(containerName)
			return nil
		},
	}

	execCommand = cli.Command{
		Name:  "exec",
		Usage: "exec a command into container",
		Action: func(context *cli.Context) error {
			if os.Getenv(ENV_EXEC_PID) != "" {
				logrus.Infof("pid callback pid %s", os.Getgid())
				return nil
			}
			if len(context.Args()) < 2 {
				return fmt.Errorf("Missing container name or command")
			}
			containerName := context.Args().Get(0)
			var commandArray []string
			for _, arg := range context.Args().Tail() {
				commandArray = append(commandArray, arg)
			}
			ExecContainer(containerName, commandArray)
			return nil
		},
	}

	stopCommand = cli.Command{
		Name:  "stop",
		Usage: "stop a container",
		Action: func(context *cli.Context) error {
			if len(context.Args()) < 1 {
				return fmt.Errorf("Missing container name")
			}
			containerName := context.Args().Get(0)
			stopContainer(containerName)
			return nil
		},
	}

	removeCommand = cli.Command{
		Name:  "rm",
		Usage: "remove unused container",
		Action: func(context *cli.Context) error {
			if len(context.Args()) < 1 {
				return fmt.Errorf("Missing container name")
			}
			containerName := context.Args().Get(0)
			removeContainer(containerName)
			return nil
		},
	}

	networkCommand = cli.Command{
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
					logrus.Infof("create network")
					if len(context.Args()) < 1 {
						return fmt.Errorf("Missing network name")
					}
					if err := network.Init(); err != nil {
						logrus.Errorf("network init error")
					}
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
)

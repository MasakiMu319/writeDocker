package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

const usage = `mydocker is a simple container runtime implementation.
					Enjoy it, just for fun.`

// createCommandList returns a list of CLI commands used in the "mydocker" application.
// The commands include:
// - initCommand: Initializes the container process and runs the user's process in the container.
// - runCommand: Creates a container with namespace and cgroups limit.
// - commitCommand: Commits a container into an image.
// - listCommand: Lists all the containers.
// - logCommand: Prints logs of a container.
// - execCommand: Executes a command into a container.
// - stopCommand: Stops a container.
// - removeCommand: Removes an unused container.
// - networkCommand: Handles container network commands.
// Example usage:
//
//	app := cli.NewApp()
//	app.Name = "mydocker"
//	app.Usage = usage
//	app.Commands = createCommandList()
//	app.Before = func(context *cli.Context) error {
//	  logrus.SetFormatter(&logrus.JSONFormatter{})
//	  logrus.SetOutput(os.Stdout)
//	  return nil
//	}
//	if err := app.Run(os.Args); err != nil {
//	  logrus.Fatal(err)
//	}
func createCommandList() []cli.Command {
	return []cli.Command{
		initCommand,
		runCommand,
		commitCommand,
		listCommand,
		logCommand,
		execCommand,
		stopCommand,
		removeCommand,
		networkCommand,
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "mydocker"
	app.Usage = usage
	app.Commands = createCommandList()
	app.Before = func(context *cli.Context) error {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetOutput(os.Stdout)
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

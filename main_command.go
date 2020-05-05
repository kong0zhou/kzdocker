package main

import (
	"fmt"
	"kzdocker/container"
	"kzdocker/log"
	"os"

	"github.com/urfave/cli/v2"
)

var runCommand = &cli.Command{
	Name: "run",
	Usage: `Create a container with namespace and cgroups limit
			mydocker run -it [command]`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "it",
			Usage: "enable tty",
		},
	},
	Action: func(context *cli.Context) error {
		if context.Args().Len() < 1 {
			return fmt.Errorf("Missing container command")
		}
		cmd := context.Args().Get(0)
		log.Info(cmd)
		tty := context.Bool("it")
		run(tty, cmd)
		return nil
	},
}

func run(tty bool, command string) {
	parent := container.NewParentProcess(tty, command)
	if err := parent.Start(); err != nil {
		log.Error(err.Error())
	}
	parent.Wait()
	os.Exit(-1)
}

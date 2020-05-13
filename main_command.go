package main

import (
	"fmt"
	"kzdocker/cgroup"
	"kzdocker/container"
	"kzdocker/log"

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
		&cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
		&cli.StringFlag{
			Name:  "cpushare",
			Usage: "cpushare limit",
		},
		&cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit",
		},
	},
	Action: func(context *cli.Context) error {
		if context.Args().Len() < 1 {
			return fmt.Errorf("Missing container command")
		}
		// cmd := context.Args().Get(0)
		// log.Info(cmd)
		var cmdArray []string
		for _, arg := range context.Args().Slice() {
			cmdArray = append(cmdArray, arg)
		}
		tty := context.Bool("it")
		resConf := &cgroup.ResourceConfig{
			MemoryLimit: context.String("m"),
			CPUSet:      context.String("cpuset"),
			CPUShare:    context.String("cpushare"),
		}
		run(tty, cmdArray, resConf)
		return nil
	},
}

func run(tty bool, command []string, res *cgroup.ResourceConfig) {
	parent := container.NewParentProcess(tty, command)
	if err := parent.Start(); err != nil {
		log.Error(err.Error())
	}
	cgroupManager := cgroup.NewCGroupManager("kzdocker-cgroup", res)
	defer cgroupManager.Destroy()
	cgroupManager.Set()
	cgroupManager.Apply(parent.Process.Pid)
	parent.Wait()
	// cgroupManager.Destroy()
	// os.Exit(-1)
	return
}

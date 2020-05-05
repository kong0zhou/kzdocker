package main

import (
	"fmt"
	"kzdocker/base"
	"kzdocker/container"
	"kzdocker/log"
	"os"

	"github.com/urfave/cli/v2"
)

func init() {
	// fmt.Println(os.Args[0])
	if os.Args[0] == `/proc/self/exe` {
		fmt.Println(os.Args[1])
		err := container.RunContainerInitProcess(os.Args[1], nil)
		if err != nil {
			log.Panic(err.Error())
		}
	}
}

func main() {
	base.InitBase()
	log.InitLog()

	app := cli.NewApp()
	app.Name = "kzdocker"
	app.Usage = "just for fun."

	app.Commands = []*cli.Command{
		runCommand,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
	}
}

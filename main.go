package main

import (
	"fmt"
	"mtsaver/cmd"
	"mtsaver/mtsaver"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Version:              mtsaver.Global.Version,
		Name:                 mtsaver.Global.AppName,
		Copyright:            "MiTo Team, https://mito-team.com",
		Description:          "7-Zip based backup arhives retension",
		Usage:                "backup arhives retension",
		EnableBashCompletion: false,
		Before:               mtsaver.SetupBeforeCommand,

		DefaultCommand: "help",
		Commands: []*cli.Command{
			&cmd.CmdVersion,
			&cmd.CmdInfo,
		},

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "7zip",
				Value: "auto",
				Usage: "Command to run 7-Zip executable. \"auto\" = try to auto-detect",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("Failed to run app")
	}
}

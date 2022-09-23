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
		Version:     mtsaver.Version,
		Name:        mtsaver.AppName,
		Copyright:   "MiTo Team, https://mito-team.com",
		Description: "7-Zip based backup arhives retension",
		Usage:       "backup arhives retension",

		//default action
		Action: cmd.CmdVersion.Action,

		Commands: []*cli.Command{
			&cmd.CmdVersion,
			&cmd.CmdInfo,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("Failed to run app")
	}
}

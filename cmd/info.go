package cmd

import (
	"fmt"
	"mtsaver/mtsaver"

	"github.com/urfave/cli/v2"
)

var CmdInfo = cli.Command{
	Name:  "info",
	Usage: "Print information about system, environment and so on",
	Action: func(ctx *cli.Context) error {
		fmt.Println(mtsaver.AppName + " version: " + mtsaver.Version)
		fmt.Println("Built with: " + mtsaver.BuiltWith)
		return nil
	},
}

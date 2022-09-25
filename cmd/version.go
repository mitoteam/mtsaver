package cmd

import (
	"fmt"
	mtsaver "mtsaver/main"

	"github.com/urfave/cli/v2"
)

var CmdVersion = cli.Command{
	Name:    "version",
	Aliases: []string{"v"},
	Usage:   "Print version number",
	Action: func(ctx *cli.Context) error {
		fmt.Println(mtsaver.Global.Version)
		return nil
	},
}
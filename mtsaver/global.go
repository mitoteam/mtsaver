package mtsaver

import (
	"runtime"

	"github.com/urfave/cli/v2"
)

var Global struct {
	AppName     string
	Version     string
	BuiltWith   string
	SevenZipCmd string
}

func init() {
	Global.AppName = "mtsaver"
	Global.Version = "1.0.0-alpha"
	Global.BuiltWith = runtime.Version()
}

func SetupBeforeCommand(ctx *cli.Context) error {
	Global.SevenZipCmd = ctx.String("7zip")

	if len(Global.SevenZipCmd) == 0 || Global.SevenZipCmd == "auto" {
		//try autodetect
		Global.SevenZipCmd = "ha ha auto"
	}

	return nil
}

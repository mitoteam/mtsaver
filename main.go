package main

import (
	_ "embed"
	"log"
	"mtsaver/app"
	"mtsaver/cmd"
)

//go:embed LICENSE.md
var licenseString string

func main() {
	app.Global.License = licenseString

	//cli application - we just let cobra to do it job
	if err := cmd.Root().Execute(); err != nil {
		log.Fatalln(err)
	}
}

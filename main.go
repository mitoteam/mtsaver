package main

import (
	"log"
	"mtsaver/cmd"
)

func main() {
	//cli application - we just let cobra to do it job
	if err := cmd.Root().Execute(); err != nil {
		log.Fatalln(err)
	}
}

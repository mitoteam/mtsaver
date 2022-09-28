package cmd

import (
	"fmt"
	"mtsaver/app"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print information about system, environment and so on",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(" --- " + app.Global.AppName + " --- ")
		fmt.Println("Version: " + app.Global.Version)
		fmt.Println("Commit: " + app.Global.Commit)
		fmt.Println("Built with: " + app.Global.BuiltWith)
		fmt.Println("7-zip command: " + app.Global.SevenZipCmd)
		fmt.Println("7-zip info: " + app.Global.SevenZipInfo)
	},
}

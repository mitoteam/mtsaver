package cmd

import (
	"fmt"
	mtsaver "mtsaver/main"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print information about system, environment and so on",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(mtsaver.Global.AppName + " version: " + mtsaver.Global.Version)
		fmt.Println("Built with: " + mtsaver.Global.BuiltWith)
		fmt.Println("7-zip command: " + mtsaver.Global.SevenZipCmd)
		fmt.Println("7-zip info: " + mtsaver.Global.SevenZipInfo)
	},
}

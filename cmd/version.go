package cmd

import (
	"fmt"
	mtsaver "mtsaver/main"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the raw version number of " + mtsaver.Global.AppName,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(mtsaver.Global.Version)
	},
}

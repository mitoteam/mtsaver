package cmd

import (
	"fmt"
	"mtsaver/app"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the raw version number of " + app.Global.AppName,

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(app.Global.Version)
		},
	})
}

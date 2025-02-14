package cmd

import (
	"fmt"
	"mtsaver/app"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "license",
		Short: "Prints license information for " + app.Global.AppName + ".",

		Run: func(cmd *cobra.Command, args []string) {
			title := app.Global.AppName + " license information"

			fmt.Println(title + "\n" + strings.Repeat("-", len(title)) + "\n")
			fmt.Println(app.Global.License)
		},
	})
}

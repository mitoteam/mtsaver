package cmd

import (
	"fmt"
	"log"
	"mtsaver/app"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "info [/path/to/directory]",
		Short: "Print information about system, environment etc. If path is given settings for that folder are printed as well.",

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(" --- " + app.Global.AppName + " --- ")
			fmt.Println("Version: " + app.Global.Version)
			fmt.Println("Commit: " + app.Global.Commit)
			fmt.Println("Built with: " + app.Global.BuiltWith)
			fmt.Println("7-zip command: " + app.Global.SevenZipCmd)
			fmt.Println("7-zip info: " + app.Global.SevenZipInfo)
			fmt.Println()

			if len(args) > 0 {
				job, err := app.NewJob(args[0])
				if err != nil {
					log.Fatalln(err)
				}

				fmt.Println(" --- Directory ---")
				fmt.Println(" --- Path:", job.Path)
				fmt.Println(" --- Settings file:", job.SettingsFilename())
				job.Settings.Print()

				fmt.Println()
			}
		},
	})
}

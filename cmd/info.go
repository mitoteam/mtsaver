package cmd

import (
	"fmt"
	"mtsaver/app"
	"path/filepath"

	"github.com/mitoteam/mttools"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "info [/path/to/directory]",
		Short: "Print information about system, environment etc. If path is given settings for that folder are printed as well.",

		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(" --- " + app.Global.AppName + " info --- ")
			fmt.Println("Version: " + app.Global.Version)
			fmt.Println("Commit: " + app.Global.Commit)
			fmt.Println("Built with: " + app.Global.BuiltWith)
			fmt.Println("7-zip command: " + app.Global.SevenZipCmd)
			fmt.Println("7-zip info: " + app.Global.SevenZipInfo)
			fmt.Println()

			job, err := app.NewJobFromArgs(args)
			if err != nil {
				return err
			}

			settings_filename := job.SettingsFilename()

			fmt.Println(" --- Directory Info ---")
			fmt.Println("Path:", job.Path)

			if mttools.IsFileExists(settings_filename) {
				fmt.Println("Settings file:", settings_filename)
				fmt.Println("\n --- Directory Settings ---")
				job.Settings.Print()
			} else {
				fmt.Printf("no settings file found (%s)\n", filepath.Base(settings_filename))
			}

			fmt.Println()

			return nil
		},
	})
}

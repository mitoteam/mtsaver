package cmd

import (
	"errors"
	"fmt"
	"mtsaver/app"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "init [/path/to/directory]",
		Short: "Creates settings file with defaults. If no path is given current directory used. --settings option can be used to specify settings file name or location explicitly.",

		RunE: func(cmd *cobra.Command, args []string) error {
			//Check path
			var path string
			if len(args) > 0 {
				path = args[0]
			} else {
				path = "." //current directory
			}

			// run Job
			job, err := app.NewJob(path)
			if err != nil {
				return err
			}

			filename := job.SettingsFilename()

			if app.IsFileExists(filename) {
				return errors.New("can not initialize existing file: " + filename)
			}

			comment := `
File created automatically by 'mtsaver init' command. There are all available
options listed here with its default values. Recomendation is to edit options you
want to change and remove all others to keep this simple.
`

			if err := job.Settings.SaveToFile(filename, comment); err != nil {
				return err
			}

			fmt.Println("Default settings written to " + filename)

			return nil
		},
	})
}

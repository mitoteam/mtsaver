package cmd

import (
	"errors"
	"fmt"
	"mtsaver/app"

	"github.com/mitoteam/mttools"

	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "init [/path/to/directory]",
		Short: "Creates settings file with defaults. If no path is given current directory used. --settings option can be used to specify settings file name or location explicitly.",

		RunE: func(cmd *cobra.Command, args []string) error {
			job, err := app.NewJobFromArgs(args)
			if err != nil {
				return err
			}

			if app.JobRuntimeOptions.DefaultsFrom != "" {
				if err := job.Settings.LoadFromFile(app.JobRuntimeOptions.DefaultsFrom); err != nil {
					return err
				}
			}

			filename := job.SettingsFilename()

			if !app.JobRuntimeOptions.Print && mttools.IsFileExists(filename) {
				return errors.New("can not initialize existing file: " + filename)
			}

			comment := `
File created automatically by 'mtsaver init' command. There are all available
options listed here with its default values. Recommendation is to edit options you
want to change and remove all others to keep this as simple as possible.
`

			if app.JobRuntimeOptions.Print {
				job.Settings.Print()
			} else {
				if err := job.Settings.SaveToFile(filename, comment); err != nil {
					return err
				}

				fmt.Println("Default settings written to " + filename)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(
		&app.JobRuntimeOptions.DefaultsFrom, "defaults-from", "",
		"settings file used to read defaults from before generating new setting with full options set",
	)

	cmd.Flags().BoolVar(
		&app.JobRuntimeOptions.Print, "print", false,
		"print default settings instead writing to file",
	)

	rootCmd.AddCommand(cmd)
}

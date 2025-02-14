package cmd

import (
	"mtsaver/app"

	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "dump [/path/to/directory]",
		Short: "Prints all archives created for directory",
		Long:  "Prints all archives created for directory with information about archives: size, age, diffs list and so on.",

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := CallParentPreRun(cmd, args); err != nil {
				return err
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			job, err := app.NewJobFromArgs(args)
			if err != nil {
				return err
			}

			job.Dump()

			return nil
		},
	}

	rootCmd.AddCommand(cmd)
}

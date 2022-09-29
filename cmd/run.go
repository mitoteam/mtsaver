package cmd

import (
	"errors"
	"mtsaver/app"

	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "run [/path/to/directory]",
		Short: "Runs backup procedure for path. If no path is given current directory is used.",

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := CallParentPreRun(cmd, args); err != nil {
				return err
			}

			//Options checks
			if app.JobRuntimeOptions.ForceFull && app.JobRuntimeOptions.ForceDiff {
				return errors.New("can not force both full and differential backups simultaneously")
			}

			return nil
		},

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

			if err = job.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(
		&app.JobRuntimeOptions.ForceFull, "force-full", false,
		"Create full archive even if conditions in settings require differential one.",
	)

	cmd.Flags().BoolVar(
		&app.JobRuntimeOptions.ForceDiff, "force-diff", false,
		"Create differential archive even if conditions in settings require full one. This option can not be used if there are no full archives created yet.",
	)

	rootCmd.AddCommand(cmd)
}

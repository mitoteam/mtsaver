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
		Use:   "run [/path/to/directory]",
		Short: "Runs backup procedure for directory",
		Long:  "Runs backup procedure for directory. If no path is given current directory is used.",

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := CallParentPreRun(cmd, args); err != nil {
				return err
			}

			//Options checks
			if mttools.CountValues(true, app.JobRuntimeOptions.ForceFull, app.JobRuntimeOptions.ForceDiff) > 1 {
				return errors.New("can not force both full and differential backups simultaneously")
			}

			//Options messages
			if app.JobRuntimeOptions.ForceFull {
				fmt.Println("Full backup forced.")
			}

			if app.JobRuntimeOptions.ForceDiff {
				fmt.Println("Differential backup forced.")
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			job, err := app.NewJobFromArgs(args)
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

	cmd.Flags().BoolVar(
		&app.JobRuntimeOptions.Solid, "solid", false,
		"Create solid archives.",
	)

	cmd.Flags().StringVar(
		&app.JobRuntimeOptions.Password, "password", "",
		"Set .7z archive password (or override 'password' in settings).",
	)

	cmd.Flags().BoolVar(
		&app.JobRuntimeOptions.EncryptFilenames, "encrypt-filenames", false,
		"Encrypt filenames in archive (used only when password is set).",
	)

	rootCmd.AddCommand(cmd)
}

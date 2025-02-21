package cmd

import (
	"fmt"
	"mtsaver/app"

	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "restore [/path/to/directory]",
		Short: "Unpacks FULL+DIFF archives to specified directory.",
		Long:  "Restores FULL+DIFF archives to specified directory. --to option is required. Directory should not exist or should be empty. If no --latest option provided programs asks interactively for the which archive to restore.",

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

			//do not run if directory has no .mtsaver.yaml and no --settings option specified
			if !job.Settings.LoadedFromFile {
				return fmt.Errorf("Directory %s does not contain %s file", job.Path, app.DefaultSettingsFilename)
			}

			//get all available archives
			job.ScanArchive()

			var to string
			var ja *app.JobArchiveFile

			if app.JobRuntimeOptions.Latest {
				ja = job.Archive.LastFile()
			}

			if err = job.Restore(to, ja); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(
		&app.JobRuntimeOptions.ForceFull, "latest", false,
		"Restore latest available FULL+DIFF pair without asking.",
	)

	cmd.Flags().StringVar(
		&app.JobRuntimeOptions.Password, "to", "",
		"[REQUIRED] Path to directory to unpack archives. Directory should not exist or should be empty.",
	)

	cmd.MarkFlagRequired("to")

	rootCmd.AddCommand(cmd)
}

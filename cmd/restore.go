package cmd

import (
	"fmt"
	"mtsaver/app"
	"os"

	"github.com/mitoteam/mttools"
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

			// check path provided
			if app.JobRuntimeOptions.To == "" {
				return fmt.Errorf("--to option is required")
			}

			if app.JobRuntimeOptions.To, err = mttools.GetDirAbsolutePath(app.JobRuntimeOptions.To); err != nil {
				//ignore "directory does not exists" error
				if err.Error() != fmt.Sprintf("\"%s\" directory does not exists", app.JobRuntimeOptions.To) {
					return err
				}
			}

			job.Log("[%s v%s] Starting directory restore: %s", app.Global.AppName, app.Global.Version, job.Path)

			if mttools.IsDirExists(app.JobRuntimeOptions.To) {
				empty, err := mttools.IsDirEmpty(app.JobRuntimeOptions.To)

				if err != nil {
					return err
				}

				if !empty {
					return fmt.Errorf("Directory %s is not empty", app.JobRuntimeOptions.To)
				}
			} else {
				if err := os.MkdirAll(app.JobRuntimeOptions.To, 0777); err != nil {
					return err
				}

				job.Log("Destination directory created: %s", app.JobRuntimeOptions.To)
			}

			//get all available archives
			job.ScanArchive(app.JobRuntimeOptions.Latest)

			var ja *app.JobArchiveFile

			if app.JobRuntimeOptions.Latest {
				ja = job.Archive.LastFile()
			} else {
				ja = nil
			}

			if err = job.Restore(app.JobRuntimeOptions.To, ja); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(
		&app.JobRuntimeOptions.Latest, "latest", false,
		"Restore latest available FULL+DIFF pair without asking.",
	)

	cmd.Flags().StringVar(
		&app.JobRuntimeOptions.To, "to", "",
		"[REQUIRED] Path to directory to unpack archives. Directory should not exist or should be empty.",
	)

	cmd.MarkFlagRequired("to")

	rootCmd.AddCommand(cmd)
}

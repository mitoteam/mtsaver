package cmd

import (
	"fmt"
	"mtsaver/app"

	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "cleanup [/path/to/directory]",
		Short: "Runs cleanup procedure for directory without creating new archive",
		Long:  "Runs cleanup procedure for directory without creating new archive. If no path is given current directory is used.",

		RunE: func(cmd *cobra.Command, args []string) error {
			job, err := app.NewJobFromArgs(args)
			if err != nil {
				return err
			}

			fmt.Println("Starting cleanup...")

			if err = job.Cleanup(); err != nil {
				return err
			}

			fmt.Println("Done.")

			return nil
		},
	}

	rootCmd.AddCommand(cmd)
}

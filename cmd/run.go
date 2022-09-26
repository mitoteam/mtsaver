package cmd

import (
	mtsaver "mtsaver/main"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run [/path/to/directory]",
	Short: "Runs backup procedure for path. If no path is given current directory is used.",

	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		if len(args) > 0 {
			path = args[0]
		} else {
			path = "." //current directory
		}

		job, err := mtsaver.NewJob(path)
		if err != nil {
			return err
		}

		if err = job.Run(); err != nil {
			return err
		}

		return nil
	},
}

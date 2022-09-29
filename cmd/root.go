// Package cmd provides CLI commands, flags and arguments handling.
// spf13/cobra based.
package cmd

import (
	"log"
	"mtsaver/app"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     app.Global.AppName,
	Version: app.Global.Version,
	Short:   "7-Zip based backup arhives retention",
	Long: `7-Zip based backup arhives retention.
Copyright: MiTo Team, https://mito-team.com`,

	//disable 'completition' subcommand
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},

	Run: func(cmd *cobra.Command, args []string) {
		//show help if no subcommand given
		cmd.Help()
	},

	PersistentPreRunE: app.SetupBeforeCommand,
}

func ExecuteCliApp() {
	rootCmd.PersistentFlags().StringVar(
		&app.Global.SevenZipCmd,
		"7zip",
		"auto",
		"Command to run 7-Zip executable. \"auto\" = try to auto-detect",
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

// CallParentPreRun calls parent command's PersistentPreRun or PersistentPreRunE hooks if they are defined.
func CallParentPreRun(cmd *cobra.Command, args []string) error {
	parent := cmd.Parent()

	if parent == nil {
		return nil
	}

	if handler := parent.PersistentPreRun; handler != nil {
		handler(cmd, args)
	}

	if handler := parent.PersistentPreRunE; handler != nil {
		if err := handler(cmd, args); err != nil {
			return err
		}
	}

	return nil
}

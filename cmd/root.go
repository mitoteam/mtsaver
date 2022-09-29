// Package cmd provides CLI commands, flags and arguments handling.
// spf13/cobra based.
package cmd

import (
	"mtsaver/app"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     app.Global.AppName,
	Version: app.Global.Version,
	Long: app.Global.AppName + ` - differential backup arhives retention tool.

Based on using 7-Zip archiver https://www.7-zip.org

Copyright: MiTo Team, https://mito-team.com`,

	//disable 'completition' subcommand
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},

	Run: func(cmd *cobra.Command, args []string) {
		//show help if no subcommand given
		cmd.Help()
	},

	PersistentPreRunE: app.SetupBeforeCommand,
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&app.Global.SevenZipCmd,
		"7zip",
		"auto",
		"Command to run 7-Zip executable. \"auto\" = try to auto-detect",
	)

	rootCmd.PersistentFlags().StringVar(
		&app.JobRuntimeOptions.SettingsFilename,
		"settings",
		app.DefaultSettingsFilename,
		"Filename or path to directory settings file. Used by 'run', 'info', 'init' commands. If filename only given it is looked for in directory itself.",
	)
}

func Root() *cobra.Command {
	return rootCmd
}

// CallParentPreRun helper function calls parent command's PersistentPreRun
// or PersistentPreRunE hooks if they are defined.
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

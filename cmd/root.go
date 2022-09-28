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

	PersistentPreRun: app.SetupBeforeCommand,
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

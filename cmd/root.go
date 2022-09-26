package cmd

import (
	"log"
	mtsaver "mtsaver/main"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     mtsaver.Global.AppName,
	Version: mtsaver.Global.Version,
	Short:   "7-Zip based backup arhives retention",
	Long: `7-Zip based backup arhives retention.
Copyright: MiTo Team, https://mito-team.com`,

	//disable 'completition' subcommand
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},

	Run: func(cmd *cobra.Command, args []string) {
		//show help if no subcommand given
		cmd.Help()
	},

	PersistentPreRun: mtsaver.SetupBeforeCommand,
}

func ExecuteCliApp() {
	rootCmd.PersistentFlags().StringVar(
		&mtsaver.Global.SevenZipCmd,
		"7zip",
		"auto",
		"Command to run 7-Zip executable. \"auto\" = try to auto-detect",
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

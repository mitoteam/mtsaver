package cmd

import (
	mtsaver "mtsaver/main"

	"github.com/urfave/cli/v2"
)

var CmdRun = cli.Command{
	Name:      "run",
	Usage:     "Runs backup procedure for path",
	ArgsUsage: "/path/to/directory",
	Action: func(ctx *cli.Context) error {
		job, err := mtsaver.NewJob(ctx.Args().Get(0))
		if err != nil {
			return err
		}

		if err = job.Run(); err != nil {
			return err
		}

		return nil
	},
}

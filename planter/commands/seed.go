package commands

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(seedCommand)
}

var seedCommand = &cobra.Command{
	Use:   "seed",
	Short: "Work with seeds",
}

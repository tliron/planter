package commands

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(pluginCommand)
}

var pluginCommand = &cobra.Command{
	Use:   "plugin",
	Short: "Work with plugins",
}

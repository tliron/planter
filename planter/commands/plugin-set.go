package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
)

func init() {
	pluginCommand.AddCommand(pluginCreateCommand)
}

var pluginCreateCommand = &cobra.Command{
	Use:   "set [PLUGIN NAME] [LOCAL FILE PATH]",
	Short: "Set a plugin",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		pluginName := args[0]
		filePath := args[1]
		SetPlugin(pluginName, filePath)
	},
}

func SetPlugin(pluginName string, filePath string) {
	err := NewClient().Planter().SetPlugin(pluginName, filePath)
	util.FailOnError(err)
}

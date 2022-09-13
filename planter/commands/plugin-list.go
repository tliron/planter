package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
)

func init() {
	pluginCommand.AddCommand(pluginListCommand)
}

var pluginListCommand = &cobra.Command{
	Use:   "list",
	Short: "List plugins",
	Run: func(cmd *cobra.Command, args []string) {
		ListPlugins()
	},
}

func ListPlugins() {
	plugins, err := NewClient().Planter().ListPlugins()
	util.FailOnError(err)
	if len(plugins) == 0 {
		return
	}
	// TODO: sort plugins by name?

	for _, plugin := range plugins {
		terminal.Println(plugin)
	}
}

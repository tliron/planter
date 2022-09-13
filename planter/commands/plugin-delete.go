package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
)

func init() {
	pluginCommand.AddCommand(pluginDeleteCommand)
	pluginDeleteCommand.Flags().BoolVarP(&all, "all", "a", false, "delete all plugins")
}

var pluginDeleteCommand = &cobra.Command{
	Use:   "delete [[PLUGIN NAME]]",
	Short: "Delete a plugin (or all plugins)",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			pluginName := args[0]
			DeletePlugin(pluginName)
		} else if all {
			DeleteAllPlugins()
		} else {
			util.Fail("must provide plugin name or specify \"--all\"")
		}
	},
}

func DeletePlugin(pluginName string) {
	err := NewClient().Planter().DeletePlugin(pluginName)
	util.FailOnError(err)
}

func DeleteAllPlugins() {
	planter := NewClient().Planter()
	plugins, err := planter.ListPlugins()
	util.FailOnError(err)
	for _, plugin := range plugins {
		log.Infof("deleting plugin: %s", plugin)
		err := planter.DeletePlugin(plugin)
		util.FailOnError(err)
	}
}

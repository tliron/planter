package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
)

func init() {
	clusterCommand.AddCommand(clusterGetCommand)
}

var clusterGetCommand = &cobra.Command{
	Use:   "get [CLUSTER NAME]",
	Short: "Gets a cluster's kubeconfig",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		clusterName := args[0]
		ClusterKubeConfig(clusterName)
	},
}

func ClusterKubeConfig(clusterName string) {
	content, err := NewClient().Planter().GetClusterKubeConfig(namespace, clusterName)
	util.FailOnError(err)
	terminal.Println(content)
}

package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
)

func init() {
	clusterCommand.AddCommand(clusterConfigCommand)
	clusterConfigCommand.Flags().StringVar(&clusterKubeconfigPath, "cluster-kubeconfig", "", "local path to cluster Kubernetes configuration (file will be copied to remote location)")
	clusterConfigCommand.Flags().StringVar(&clusterContext, "cluster-context", "", "override current context in cluster Kubernetes configuration")
}

var clusterConfigCommand = &cobra.Command{
	Use:   "config [CLUSTER NAME]",
	Short: "Configure a cluster",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		clusterName := args[0]
		ConfigCluster(clusterName)
	},
}

func ConfigCluster(clusterName string) {
	_, err := NewClient().Planter().ConfigClusterFromFile(namespace, clusterName, clusterKubeconfigPath, clusterContext)
	util.FailOnError(err)
}

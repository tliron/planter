package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
)

var clusterKubeconfigPath string
var clusterContext string

func init() {
	clusterCommand.AddCommand(clusterCreateCommand)
	clusterCreateCommand.Flags().StringVar(&clusterKubeconfigPath, "cluster-kubeconfig", "", "local path to cluster Kubernetes configuration (file will be copied to remote location)")
	clusterCreateCommand.Flags().StringVar(&clusterContext, "cluster-context", "", "override current context in cluster Kubernetes configuration")
}

var clusterCreateCommand = &cobra.Command{
	Use:   "create [CLUSTER NAME]",
	Short: "Create a cluster",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		clusterName := args[0]
		CreateCluster(clusterName)
	},
}

func CreateCluster(clusterName string) {
	var err error
	if clusterKubeconfigPath != "" {
		_, err = NewClient().Planter().CreateClusterFromFile(namespace, clusterName, clusterKubeconfigPath, clusterContext)
	} else {
		_, err = NewClient().Planter().CreateClusterWithURL(namespace, clusterName, "", "")
	}
	util.FailOnError(err)
}

package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
)

func init() {
	clusterCommand.AddCommand(clusterDeleteCommand)
	clusterDeleteCommand.Flags().BoolVarP(&all, "all", "a", false, "delete all clusters")
}

var clusterDeleteCommand = &cobra.Command{
	Use:   "delete [[CLUSTER NAME]]",
	Short: "Delete a cluster (or all clusters)",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			clusterName := args[0]
			DeleteCluster(clusterName)
		} else if all {
			DeleteAllClusters()
		} else {
			util.Fail("must provide cluster name or specify \"--all\"")
		}
	},
}

func DeleteCluster(clusterName string) {
	err := NewClient().Planter().DeleteCluster(namespace, clusterName)
	util.FailOnError(err)
}

func DeleteAllClusters() {
	planter := NewClient().Planter()
	clusters, err := planter.ListClusters()
	util.FailOnError(err)
	for _, cluster := range clusters.Items {
		log.Infof("deleting cluster: %s/%s", cluster.Namespace, cluster.Name)
		err := planter.DeleteCluster(cluster.Namespace, cluster.Name)
		util.FailOnError(err)
	}
}

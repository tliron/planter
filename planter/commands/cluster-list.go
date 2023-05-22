package commands

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tliron/go-ard"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/transcribe"
	"github.com/tliron/kutil/util"
	resources "github.com/tliron/planter/resources/planter.nephio.org/v1alpha1"
)

func init() {
	clusterCommand.AddCommand(clusterListCommand)
}

var clusterListCommand = &cobra.Command{
	Use:   "list",
	Short: "List clusters",
	Run: func(cmd *cobra.Command, args []string) {
		ListClusters()
	},
}

func ListClusters() {
	clusters, err := NewClient().Planter().ListClusters()
	util.FailOnError(err)
	if len(clusters.Items) == 0 {
		return
	}
	// TODO: sort clusters by name? they seem already sorted!

	switch format {
	case "":
		table := terminal.NewTable(maxWidth, "Name", "KubeConfigURL", "Context")
		for _, cluster := range clusters.Items {
			table.Add(cluster.Name, cluster.Status.KubeConfigURL, cluster.Status.Context)
		}
		table.Print()

	case "bare":
		for _, cluster := range clusters.Items {
			terminal.Println(cluster.Name)
		}

	default:
		list := make(ard.List, len(clusters.Items))
		for index, cluster := range clusters.Items {
			list[index] = resources.ClusterToARD(&cluster)
		}
		transcribe.Print(list, format, os.Stdout, strict, pretty)
	}
}

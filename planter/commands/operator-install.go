package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
)

var clusterRole string
var sourceRegistry string

func init() {
	operatorCommand.AddCommand(operatorInstallCommand)
	operatorInstallCommand.Flags().StringVarP(&clusterRole, "role", "e", "cluster-admin", "cluster role")
	operatorInstallCommand.Flags().StringVarP(&sourceRegistry, "registry", "g", "docker.io", "source registry host")
	operatorInstallCommand.Flags().BoolVarP(&wait, "wait", "w", false, "wait for installation to succeed")
}

var operatorInstallCommand = &cobra.Command{
	Use:   "install",
	Short: "Install the Planter operator",
	Run: func(cmd *cobra.Command, args []string) {
		err := NewClient().Planter().InstallOperator(sourceRegistry, wait)
		util.FailOnError(err)
	},
}

package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
)

func init() {
	seedCommand.AddCommand(seedDeleteCommand)
	seedDeleteCommand.Flags().BoolVarP(&all, "all", "a", false, "delete all seeds")
}

var seedDeleteCommand = &cobra.Command{
	Use:   "delete [[SEED NAME]]",
	Short: "Delete a seed (or all seeds)",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			seedName := args[0]
			DeleteSeed(seedName)
		} else if all {
			DeleteAllSeeds()
		} else {
			util.Fail("must provide seed name or specify \"--all\"")
		}
	},
}

func DeleteSeed(seedName string) {
	err := NewClient().Planter().DeleteSeed(namespace, seedName)
	util.FailOnError(err)
}

func DeleteAllSeeds() {
	planter := NewClient().Planter()
	seeds, err := planter.ListSeeds()
	util.FailOnError(err)
	for _, seed := range seeds.Items {
		log.Infof("deleting seed: %s/%s", seed.Namespace, seed.Name)
		err := planter.DeleteSeed(seed.Namespace, seed.Name)
		util.FailOnError(err)
	}
}

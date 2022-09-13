package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
)

func init() {
	seedCommand.AddCommand(seedPlantCommand)
}

var seedPlantCommand = &cobra.Command{
	Use:   "plant [SEED NAME]",
	Short: "Plant a seed",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		seedName := args[0]
		PlantSeed(seedName)
	},
}

func PlantSeed(seedName string) {
	_, err := NewClient().Planter().SetSeedPlanted(namespace, seedName, true)
	util.FailOnError(err)
}

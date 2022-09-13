package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
)

func init() {
	seedCommand.AddCommand(seedGetCommand)
	seedGetCommand.Flags().BoolVarP(&planted, "planted", "a", false, "whether to get the planted content of the seed")
}

var seedGetCommand = &cobra.Command{
	Use:   "get [SEED NAME]",
	Short: "Gets a seed's content",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		seedName := args[0]
		SeedContent(seedName)
	},
}

func SeedContent(seedName string) {
	content, err := NewClient().Planter().GetSeedContent(namespace, seedName, planted)
	util.FailOnError(err)
	terminal.Println(content)
}

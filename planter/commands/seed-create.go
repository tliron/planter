package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
)

var url string
var filePath string

func init() {
	seedCommand.AddCommand(seedCreateCommand)
	seedCreateCommand.Flags().StringVarP(&url, "url", "u", "", "remote URL to seed (to be accessed by the operator)")
	seedCreateCommand.Flags().StringVarP(&filePath, "file", "f", "", "local path to seed (file will be copied to remote location)")
	seedCreateCommand.Flags().BoolVarP(&planted, "plant", "a", false, "whether to plant the seed after creating it")
}

var seedCreateCommand = &cobra.Command{
	Use:   "create [SEED NAME]",
	Short: "Create a seed",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		seedName := args[0]
		if (url == "") && (filePath == "") {
			util.Fail("must set either \"--url\" or \"--file\"")
		} else if (url != "") && (filePath != "") {
			util.Fail("cannot set both \"--url\" and \"--file\"")
		}
		CreateSeed(seedName, url, filePath)
	},
}

func CreateSeed(seedName string, url string, localPath string) {
	planter := NewClient().Planter()
	var err error
	if url != "" {
		_, err = planter.CreateSeedWithURL(namespace, seedName, url, planted)
	} else if localPath != "" {
		_, err = planter.CreateSeedFromFile(namespace, seedName, localPath, planted)
	}
	util.FailOnError(err)
}

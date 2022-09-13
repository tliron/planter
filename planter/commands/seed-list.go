package commands

import (
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/transcribe"
	"github.com/tliron/kutil/util"
	resources "github.com/tliron/planter/resources/planter.nephio.org/v1alpha1"
)

func init() {
	seedCommand.AddCommand(seedListCommand)
}

var seedListCommand = &cobra.Command{
	Use:   "list",
	Short: "List seeds",
	Run: func(cmd *cobra.Command, args []string) {
		ListSeeds()
	},
}

func ListSeeds() {
	seeds, err := NewClient().Planter().ListSeeds()
	util.FailOnError(err)
	if len(seeds.Items) == 0 {
		return
	}
	// TODO: sort seeds by name? they seem already sorted!

	switch format {
	case "":
		table := terminal.NewTable(maxWidth, "Name", "URL", "Planted", "PlantedPath")
		for _, seed := range seeds.Items {
			table.Add(seed.Name, seed.Spec.SeedURL, strconv.FormatBool(seed.Spec.Planted), seed.Status.PlantedPath)
		}
		table.Print()

	case "bare":
		for _, service := range seeds.Items {
			terminal.Println(service.Name)
		}

	default:
		list := make(ard.List, len(seeds.Items))
		for index, seed := range seeds.Items {
			list[index] = resources.SeedToARD(&seed)
		}
		transcribe.Print(list, format, os.Stdout, strict, pretty)
	}
}

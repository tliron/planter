package main

import (
	"github.com/tliron/kutil/util"
	"github.com/tliron/planter/planter/commands"

	_ "github.com/tliron/kutil/logging/simple"
)

func main() {
	util.ExitOnSIGTERM()
	commands.Execute()
	util.Exit(0)
}

package commands

import (
	contextpkg "context"

	"github.com/tliron/commonlog"
)

const toolName = "planter"

var context = contextpkg.TODO()

var log = commonlog.GetLogger(toolName)

var tail int
var follow bool
var all bool
var wait bool
var planted bool

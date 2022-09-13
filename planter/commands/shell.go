package commands

import (
	"fmt"
	"io"
	"os"

	"github.com/tliron/kutil/kubernetes"
	"github.com/tliron/kutil/util"
	"github.com/tliron/planter/controller"
	"golang.org/x/term"
)

func Shell(appNameSuffix string, containerName string) {
	// We need stdout to be in "raw" mode
	fd := int(os.Stdout.Fd())
	state, err := term.MakeRaw(fd)
	util.FailOnError(err)
	handle := util.OnExitError(func() error {
		return term.Restore(fd, state)
	})
	err = NewClient().Shell(appNameSuffix, containerName, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		handle.Cancel()
		term.Restore(fd, state)
		util.FailOnError(err)
	}
}

func (self *Client) Shell(appNameSuffix string, containerName string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	appName := fmt.Sprintf("%s-%s", controller.NamePrefix, appNameSuffix)

	if podName, err := kubernetes.GetFirstPodName(self.Context, self.Kubernetes, self.Namespace, appName); err == nil {
		return kubernetes.Exec(self.REST, self.Config, self.Namespace, podName, containerName, stdin, stdout, stderr, true, "sh")
	} else {
		return err
	}
}

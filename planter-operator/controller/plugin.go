package controller

import (
	"os/exec"
	"sort"
	"strings"

	"github.com/tliron/kutil/util"
)

func (self *Controller) execPlugins(content string) (string, error) {
	if plugins, err := self.Client.ListPluginPaths(); err == nil {
		sort.Strings(plugins)
		for _, plugin := range plugins {
			self.Log.Infof("running plugin: %s", plugin)
			command := exec.Command(plugin)
			command.Stdin = strings.NewReader(content)
			if content_, err := command.Output(); err == nil {
				content = util.BytesToString(content_)
			} else {
				return "", err
			}
		}
	} else {
		return "", err
	}

	return content, nil
}

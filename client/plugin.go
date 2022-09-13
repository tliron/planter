package client

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tliron/kutil/kubernetes"
)

var executable int64 = 0700

func (self *Client) SetPlugin(name string, path string) error {
	appName := fmt.Sprintf("%s-%s", self.NamePrefix, "operator")
	cachePath := filepath.Join(self.getPluginCachePath(name), filepath.Base(path))

	if file, err := os.Open(path); err == nil {
		defer file.Close()

		if podNames, err := kubernetes.GetPodNames(self.Context, self.Kubernetes, self.Namespace, appName); err == nil {
			for _, podName := range podNames {
				self.Log.Infof("setting plugin %q in operator pod: %s/%s", name, self.Namespace, podName)
				if err := self.WriteToContainer(self.Namespace, podName, "operator", file, cachePath, &executable); err != nil {
					return err
				}
			}
		} else {
			return err
		}
	} else {
		return err
	}

	return nil
}

func (self *Client) DeletePlugin(name string) error {
	appName := fmt.Sprintf("%s-%s", self.NamePrefix, "operator")
	cachePath := self.getPluginCachePath(name)
	if podNames, err := kubernetes.GetPodNames(self.Context, self.Kubernetes, self.Namespace, appName); err == nil {
		for _, podName := range podNames {
			if err := self.Exec(self.Namespace, podName, "operator", nil, nil, "rm", "--force", "--recursive", cachePath); err != nil {
				return err
			}
		}
	} else {
		return err
	}

	return nil
}

func (self *Client) ListPlugins() ([]string, error) {
	appName := fmt.Sprintf("%s-%s", self.NamePrefix, "operator")
	if podName, err := kubernetes.GetFirstPodName(self.Context, self.Kubernetes, self.Namespace, appName); err == nil {
		var buffer bytes.Buffer
		if err := self.Exec(self.Namespace, podName, "operator", nil, &buffer, "ls", "-1", filepath.Join(self.CachePath, "plugins")); err == nil {
			var names []string
			for _, filename := range strings.Split(strings.TrimRight(buffer.String(), "\n"), "\n") {
				names = append(names, filename)
			}
			return names, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (self *Client) ListPluginPaths() ([]string, error) {
	appName := fmt.Sprintf("%s-%s", self.NamePrefix, "operator")
	if podName, err := kubernetes.GetFirstPodName(self.Context, self.Kubernetes, self.Namespace, appName); err == nil {
		var buffer bytes.Buffer
		if err := self.Exec(self.Namespace, podName, "operator", nil, &buffer, "find", filepath.Join(self.CachePath, "plugins"), "-type", "f"); err == nil {
			var names []string
			for _, filename := range strings.Split(strings.TrimRight(buffer.String(), "\n"), "\n") {
				names = append(names, filename)
			}
			return names, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (self *Client) getPluginCachePath(pluginName string) string {
	return filepath.Join(self.CachePath, "plugins", pluginName)
}

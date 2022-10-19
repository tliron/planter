package client

import (
	"fmt"
	"strings"

	"github.com/tliron/kutil/kubernetes"
)

func (self *Client) GetContent(url string) (string, error) {
	appName := fmt.Sprintf("%s-operator", self.NamePrefix)

	// TODO: we need a custom executable to fetch the URL content
	if podName, err := kubernetes.GetFirstPodName(self.Context, self.Kubernetes, self.Namespace, appName); err == nil {
		var builder strings.Builder
		if err := self.Exec(self.Namespace, podName, "operator", nil, &builder, "cat", url); err == nil {
			return strings.TrimRight(builder.String(), "\n"), nil
		} else {
			return "", fmt.Errorf("no content at: %s", url)
		}
	} else {
		return "", err
	}
}

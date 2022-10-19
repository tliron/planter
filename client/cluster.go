package client

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tliron/kutil/kubernetes"
	resources "github.com/tliron/planter/resources/planter.nephio.org/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (self *Client) GetCluster(namespace string, clusterName string) (*resources.Cluster, error) {
	if namespace == "" {
		namespace = self.Namespace
	}

	if cluster, err := self.Planter.PlanterV1alpha1().Clusters(namespace).Get(self.Context, clusterName, meta.GetOptions{}); err == nil {
		// When retrieved from cache the GVK may be empty
		if cluster.Kind == "" {
			cluster = cluster.DeepCopy()
			cluster.APIVersion, cluster.Kind = resources.ClusterGVK.ToAPIVersionAndKind()
		}
		return cluster, nil
	} else {
		return nil, err
	}
}

func (self *Client) GetClusterKubeConfig(namespace string, clusterName string) (string, error) {
	if namespace == "" {
		namespace = self.Namespace
	}

	if cluster, err := self.GetCluster(namespace, clusterName); err == nil {
		if cluster.Status.KubeConfigURL != "" {
			return self.GetContent(cluster.Status.KubeConfigURL)
		} else {
			return "", fmt.Errorf("cluster not configured: %s/%s", namespace, clusterName)
		}
	} else {
		return "", err
	}
}

func (self *Client) ListClusters() (*resources.ClusterList, error) {
	// TODO: all clusters in cluster mode
	return self.Planter.PlanterV1alpha1().Clusters(self.Namespace).List(self.Context, meta.ListOptions{})
}

func (self *Client) UpdateClusterSpec(cluster *resources.Cluster) (*resources.Cluster, error) {
	if cluster_, err := self.Planter.PlanterV1alpha1().Clusters(cluster.Namespace).Update(self.Context, cluster, meta.UpdateOptions{}); err == nil {
		// When retrieved from cache the GVK may be empty
		if cluster_.Kind == "" {
			cluster_ = cluster_.DeepCopy()
			cluster_.APIVersion, cluster_.Kind = resources.ClusterGVK.ToAPIVersionAndKind()
		}
		return cluster_, nil
	} else {
		return cluster, err
	}
}

func (self *Client) UpdateClusterStatus(cluster *resources.Cluster) (*resources.Cluster, error) {
	if cluster_, err := self.Planter.PlanterV1alpha1().Clusters(cluster.Namespace).UpdateStatus(self.Context, cluster, meta.UpdateOptions{}); err == nil {
		// When retrieved from cache the GVK may be empty
		if cluster_.Kind == "" {
			cluster_ = cluster_.DeepCopy()
			cluster_.APIVersion, cluster_.Kind = resources.ClusterGVK.ToAPIVersionAndKind()
		}
		return cluster_, nil
	} else {
		return cluster, err
	}
}

func (self *Client) DeleteCluster(namespace string, clusterName string) error {
	if namespace == "" {
		namespace = self.Namespace
	}

	return self.Planter.PlanterV1alpha1().Clusters(namespace).Delete(self.Context, clusterName, meta.DeleteOptions{})
}

func (self *Client) CreateClusterFromFile(namespace string, clusterName string, path string, context string) (*resources.Cluster, error) {
	if namespace == "" {
		namespace = self.Namespace
	}

	if cachePath, err := self.writeClusterConfig(namespace, clusterName, path); err == nil {
		return self.CreateClusterWithURL(namespace, clusterName, cachePath, context)
	} else {
		return nil, err
	}
}

func (self *Client) ConfigClusterFromFile(namespace string, clusterName string, path string, context string) (*resources.Cluster, error) {
	if namespace == "" {
		namespace = self.Namespace
	}

	if cachePath, err := self.writeClusterConfig(namespace, clusterName, path); err == nil {
		if cluster, err := self.GetCluster(namespace, clusterName); err == nil {
			cluster.Status.KubeConfigURL = cachePath
			cluster.Status.Context = context
			return self.UpdateClusterStatus(cluster)
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (self *Client) CreateClusterWithURL(namespace string, clusterName string, url string, context string) (*resources.Cluster, error) {
	if namespace == "" {
		namespace = self.Namespace
	}

	cluster := &resources.Cluster{
		ObjectMeta: meta.ObjectMeta{
			Name:      clusterName,
			Namespace: namespace,
		},
	}

	if cluster_, err := self.createCluster(namespace, clusterName, cluster); err == nil {
		cluster_.Status.KubeConfigURL = url
		cluster_.Status.Context = context
		return self.UpdateClusterStatus(cluster_)
	} else {
		return nil, err
	}
}

func (self *Client) createCluster(namespace string, clusterName string, cluster *resources.Cluster) (*resources.Cluster, error) {
	if cluster, err := self.Planter.PlanterV1alpha1().Clusters(namespace).Create(self.Context, cluster, meta.CreateOptions{}); err == nil {
		return cluster, nil
	} else if errors.IsAlreadyExists(err) {
		return self.Planter.PlanterV1alpha1().Clusters(namespace).Get(self.Context, clusterName, meta.GetOptions{})
	} else {
		return nil, err
	}
}

func (self *Client) writeClusterConfig(namespace string, clusterName string, path string) (string, error) {
	cachePath := self.getClusterCachePath(namespace, clusterName)

	if file, err := os.Open(path); err == nil {
		defer file.Close()

		appName := fmt.Sprintf("%s-%s", self.NamePrefix, "operator")
		if podNames, err := kubernetes.GetPodNames(self.Context, self.Kubernetes, self.Namespace, appName); err == nil {
			for _, podName := range podNames {
				if err := self.WriteToContainer(self.Namespace, podName, "operator", file, cachePath, nil); err != nil {
					return "", err
				}
			}
		} else {
			return "", err
		}
	} else {
		return "", err
	}

	return cachePath, nil
}

func (self *Client) getClusterCachePath(namespace string, clusterName string) string {
	return filepath.Join(self.CachePath, "clusters", namespace, clusterName+".yaml")
}

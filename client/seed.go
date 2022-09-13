package client

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tliron/kutil/kubernetes"
	resources "github.com/tliron/planter/resources/planter.nephio.org/v1alpha1"
	kuberneteserrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (self *Client) GetSeed(namespace string, seedName string) (*resources.Seed, error) {
	if namespace == "" {
		namespace = self.Namespace
	}

	if seed, err := self.Planter.PlanterV1alpha1().Seeds(namespace).Get(self.Context, seedName, meta.GetOptions{}); err == nil {
		// When retrieved from cache the GVK may be empty
		if seed.Kind == "" {
			seed = seed.DeepCopy()
			seed.APIVersion, seed.Kind = resources.SeedGVK.ToAPIVersionAndKind()
		}
		return seed, nil
	} else {
		return nil, err
	}
}

func (self *Client) GetSeedContent(namespace string, seedName string, planted bool) (string, error) {
	if namespace == "" {
		namespace = self.Namespace
	}

	if seed, err := self.GetSeed(namespace, seedName); err == nil {
		url := seed.Spec.SeedURL
		if planted {
			url = seed.Status.PlantedPath
			if url == "" {
				return "", errors.New("seed not planted")
			}
		}

		return self.GetContent(url)
	} else {
		return "", err
	}
}

func (self *Client) ListSeeds() (*resources.SeedList, error) {
	// TODO: all seeds in cluster mode
	return self.Planter.PlanterV1alpha1().Seeds(self.Namespace).List(self.Context, meta.ListOptions{})
}

func (self *Client) UpdateSeedSpec(seed *resources.Seed) (*resources.Seed, error) {
	if seed_, err := self.Planter.PlanterV1alpha1().Seeds(seed.Namespace).Update(self.Context, seed, meta.UpdateOptions{}); err == nil {
		// When retrieved from cache the GVK may be empty
		if seed_.Kind == "" {
			seed_ = seed_.DeepCopy()
			seed_.APIVersion, seed_.Kind = resources.SeedGVK.ToAPIVersionAndKind()
		}
		return seed_, nil
	} else {
		return seed, err
	}
}

func (self *Client) UpdateSeedStatus(seed *resources.Seed) (*resources.Seed, error) {
	if seed_, err := self.Planter.PlanterV1alpha1().Seeds(seed.Namespace).UpdateStatus(self.Context, seed, meta.UpdateOptions{}); err == nil {
		// When retrieved from cache the GVK may be empty
		if seed_.Kind == "" {
			seed_ = seed_.DeepCopy()
			seed_.APIVersion, seed_.Kind = resources.SeedGVK.ToAPIVersionAndKind()
		}
		return seed_, nil
	} else {
		return seed, err
	}
}

func (self *Client) DeleteSeed(namespace string, seedName string) error {
	if namespace == "" {
		namespace = self.Namespace
	}

	return self.Planter.PlanterV1alpha1().Seeds(namespace).Delete(self.Context, seedName, meta.DeleteOptions{})
}

func (self *Client) CreateSeedFromFile(namespace string, seedName string, path string, planted bool) (*resources.Seed, error) {
	if namespace == "" {
		namespace = self.Namespace
	}

	cachePath := self.getSeedCachePath(seedName)

	if file, err := os.Open(path); err == nil {
		defer file.Close()

		appName := fmt.Sprintf("%s-%s", self.NamePrefix, "operator")
		if podNames, err := kubernetes.GetPodNames(self.Context, self.Kubernetes, self.Namespace, appName); err == nil {
			for _, podName := range podNames {
				if err := self.WriteToContainer(self.Namespace, podName, "operator", file, cachePath, nil); err != nil {
					return nil, err
				}
			}
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}

	return self.CreateSeedWithURL(namespace, seedName, cachePath, planted)
}

func (self *Client) CreateSeedWithURL(namespace string, seedName string, url string, planted bool) (*resources.Seed, error) {
	if namespace == "" {
		namespace = self.Namespace
	}

	seed := &resources.Seed{
		ObjectMeta: meta.ObjectMeta{
			Name:      seedName,
			Namespace: namespace,
		},
		Spec: resources.SeedSpec{
			SeedURL: url,
			Planted: planted,
		},
	}

	return self.createSeed(namespace, seedName, seed)
}

func (self *Client) SetSeedPlanted(namespace string, seedName string, planted bool) (*resources.Seed, error) {
	if namespace == "" {
		namespace = self.Namespace
	}

	if seed, err := self.GetSeed(namespace, seedName); err == nil {
		if seed.Spec.Planted != planted {
			seed.Spec.Planted = planted
			return self.UpdateSeedSpec(seed)
		} else {
			return seed, nil
		}
	} else {
		return nil, err
	}
}

func (self *Client) PlantSeed(seed *resources.Seed, content string) (*resources.Seed, error) {
	cachePath := self.getSeedPlantedCachePath(seed.Namespace, seed.Name)

	appName := fmt.Sprintf("%s-%s", self.NamePrefix, "operator")
	if podNames, err := kubernetes.GetPodNames(self.Context, self.Kubernetes, self.Namespace, appName); err == nil {
		for _, podName := range podNames {
			if err := self.WriteToContainer(self.Namespace, podName, "operator", strings.NewReader(content), cachePath, nil); err != nil {
				return nil, err
			}
		}
	} else {
		return nil, err
	}

	seed.Status.PlantedPath = cachePath
	return self.UpdateSeedStatus(seed)
}

func (self *Client) createSeed(namespace string, seedName string, seed *resources.Seed) (*resources.Seed, error) {
	if cluster, err := self.Planter.PlanterV1alpha1().Seeds(namespace).Create(self.Context, seed, meta.CreateOptions{}); err == nil {
		return cluster, nil
	} else if kuberneteserrors.IsAlreadyExists(err) {
		return self.Planter.PlanterV1alpha1().Seeds(namespace).Get(self.Context, seedName, meta.GetOptions{})
	} else {
		return nil, err
	}
}

func (self *Client) getSeedCachePath(seedName string) string {
	return filepath.Join(self.CachePath, "seeds", self.Namespace, seedName+".yaml")
}

func (self *Client) getSeedPlantedCachePath(namespace string, seedName string) string {
	return filepath.Join(self.CachePath, "planted", namespace, seedName+".yaml")
}

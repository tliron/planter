package controller

import (
	"strings"

	"github.com/tliron/kutil/ard"
	kuberneteserrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

//
// Resources
//

type Resources struct {
	resources map[string]ard.StringMap
	queue     []ard.StringMap
	clients   *ClusterClients
}

func NewResources(content string, controller *Controller) (*Resources, error) {
	//transcribe.PrintYAML(resource, os.Stdout, false, false)
	if content_, err := ard.ReadAllYAML(strings.NewReader(content)); err == nil {
		content__, _ := ard.NormalizeStringMaps(content_)
		content_ = content__.(ard.List)

		self := Resources{
			resources: make(map[string]ard.StringMap),
			clients:   controller.NewClusterClients(),
		}

		for _, resource := range content_ {
			if resource_, ok := resource.(ard.StringMap); ok {
				self.queue = append(self.queue, resource_)
				self.resources[getID(resource_)] = resource_
			}
		}

		return &self, nil
	} else {
		return nil, err
	}
}

func (self *Resources) Reconcile(resource ard.StringMap) error {
	id := getID(resource)
	self.clients.controller.Log.Infof("reconciling: %s", id)

	for _, dependency := range getDependencies(resource) {
		if self.queued(dependency) {
			self.clients.controller.Log.Infof("reconciling dependency: %s", dependency)
			if err := self.Reconcile(self.resources[dependency]); err != nil {
				return err
			}
			self.unqueue(dependency)
		}
	}

	object := &unstructured.Unstructured{Object: resource}

	clusterName := getClusterName(resource)
	if clusterName == "SELF" {
		self.clients.controller.Log.Infof("creating resource: %s/%s %s/%s", object.GetAPIVersion(), object.GetKind(), object.GetNamespace(), object.GetName())
		if _, err := self.clients.controller.Dynamic.CreateResource(object); err != nil {
			if kuberneteserrors.IsAlreadyExists(err) {
				self.clients.controller.Log.Infof("resource already exists: %s/%s %s/%s", object.GetAPIVersion(), object.GetKind(), object.GetNamespace(), object.GetName())
			} else {
				return err
			}
		}
	} else if clusterName != "" {
		cluster := self.getCluster(clusterName)
		if cluster != nil {
			if client, err := self.clients.Get(clusterName); err == nil {
				self.clients.controller.Log.Infof("creating resource at %s: %s/%s %s/%s", clusterName, object.GetAPIVersion(), object.GetKind(), object.GetNamespace(), object.GetName())
				if _, err := client.Dynamic(object.GetNamespace()).CreateResource(object); err != nil {
					if kuberneteserrors.IsAlreadyExists(err) {
						self.clients.controller.Log.Infof("resource already exists: %s/%s %s/%s", object.GetAPIVersion(), object.GetKind(), object.GetNamespace(), object.GetName())
					} else {
						return err
					}
				}
			} else {
				if !kuberneteserrors.IsNotFound(err) {
					self.clients.controller.Log.Warningf("cluster not found: %s", clusterName)
					return err
				}
			}
		} else {
			self.clients.controller.Log.Warningf("unknown cluster: %s", clusterName)
		}
	}

	return nil
}

func (self *Resources) next() ard.StringMap {
	if len(self.queue) > 0 {
		resource := self.queue[0]
		self.queue = self.queue[1:]
		return resource
	} else {
		return nil
	}
}

func (self *Resources) queued(id string) bool {
	for _, resource := range self.queue {
		if getID(resource) == id {
			return true
		}
	}
	return false
}

func (self *Resources) unqueue(id string) bool {
	for index, resource := range self.queue {
		if getID(resource) == id {
			self.queue = append(self.queue[0:index], self.queue[index+1:]...)
			return true
		}
	}
	return false
}

func (self *Resources) getCluster(clusterName string) ard.StringMap {
	if cluster, ok := self.resources["planter.nephio.org/v1alpha1|Cluster|"+clusterName]; ok {
		return cluster
	} else {
		return nil
	}
}

// Utils

func getID(resource ard.StringMap) string {
	resource_ := ard.NewNode(resource)
	apiVersion, _ := resource_.Get("apiVersion").String()
	kind, _ := resource_.Get("kind").String()
	namespace, _ := resource_.Get("metadata").Get("namespace").String()
	name, _ := resource_.Get("metadata").Get("name").String()
	return apiVersion + "|" + kind + "|" + namespace + "|" + name
}

func getClusterName(resource ard.StringMap) string {
	clusterName, _ := ard.NewNode(resource).Get("metadata").Get("annotations").Get("planter.nephio.org/cluster").String()
	return clusterName
}

func getDependencies(resource ard.StringMap) []string {
	dependencies, _ := ard.NewNode(resource).Get("metadata").Get("annotations").Get("planter.nephio.org/dependencies").String()
	r := strings.Split(dependencies, ",")
	clusterName := getClusterName(resource)
	if clusterName != "" {
		dependency := "planter.nephio.org/v1alpha1|Cluster|" + clusterName
		r = append(r, dependency)
	}
	return r
}

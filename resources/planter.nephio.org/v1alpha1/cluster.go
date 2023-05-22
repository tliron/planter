package v1alpha1

import (
	"fmt"

	"github.com/tliron/go-ard"
	group "github.com/tliron/planter/resources/planter.nephio.org"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ClusterGVK = SchemeGroupVersion.WithKind(ClusterKind)

const (
	ClusterKind     = "Cluster"
	ClusterListKind = "ClusterList"

	ClusterSingular  = "cluster"
	ClusterPlural    = "clusters"
	ClusterShortName = "pc" // = Planter Cluster
)

//
// Cluster
//

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Cluster struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterSpec   `json:"spec"`
	Status ClusterStatus `json:"status"`
}

type ClusterSpec struct {
	Nodes []ClusterNode `json:"nodes,omitempty"` // Grouped node configurations
	WANs  []string      `json:"wans"`            // WANs for reaching this cluster
}

type ClusterNode struct {
	Count        int               `json:"count"`                  // Amount of nodes in this group
	Labels       map[string]string `json:"labels"`                 // Labels for the nodes in this group
	KernelConfig string            `json:"kernelConfig,omitempty"` // Optional name of KernelConfig resource for the nodes in this group
	MemoryConfig string            `json:"memoryConfig,omitempty"` // Optional name of MemoryConfig resource for the nodes in this group
}

type ClusterStatus struct {
	KubeConfigURL string `json:"kubeconfigUrl,omitempty"` // Full URL of kubeconfig (YAML file)
	Context       string `json:"context,omitempty"`       // Name of context to use in kubeconfig
}

//
// ClusterList
//

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ClusterList struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata"`

	Items []Cluster `json:"items"`
}

//
// ClusterCustomResourceDefinition
//

// See: assets/kubernetes/custom-resource-definitions.yaml

var ClusterResourcesName = fmt.Sprintf("%s.%s", ClusterPlural, group.GroupName)

var ClusterCustomResourceDefinition = apiextensions.CustomResourceDefinition{
	ObjectMeta: meta.ObjectMeta{
		Name: ClusterResourcesName,
	},
	Spec: apiextensions.CustomResourceDefinitionSpec{
		Group: group.GroupName,
		Names: apiextensions.CustomResourceDefinitionNames{
			Singular: ClusterSingular,
			Plural:   ClusterPlural,
			Kind:     ClusterKind,
			ListKind: ClusterListKind,
			ShortNames: []string{
				ClusterShortName,
			},
			Categories: []string{
				"all", // will appear in "kubectl get all"
			},
		},
		Scope: apiextensions.NamespaceScoped,
		Versions: []apiextensions.CustomResourceDefinitionVersion{
			{
				Name:    Version,
				Served:  true,
				Storage: true, // one and only one version must be marked with storage=true
				Subresources: &apiextensions.CustomResourceSubresources{ // requires CustomResourceSubresources feature gate enabled
					Status: &apiextensions.CustomResourceSubresourceStatus{},
				},
				Schema: &apiextensions.CustomResourceValidation{
					OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
						Description: "Planter cluster",
						Type:        "object",
						Required:    []string{"spec"},
						Properties: map[string]apiextensions.JSONSchemaProps{
							"spec": {
								Type: "object",
								Properties: map[string]apiextensions.JSONSchemaProps{
									"nodes": {
										Description: "Grouped node configurations",
										Type:        "array",
										Items: &apiextensions.JSONSchemaPropsOrArray{
											Schema: &apiextensions.JSONSchemaProps{
												Type: "object",
												Properties: map[string]apiextensions.JSONSchemaProps{
													"count": {
														Description: "Amount of nodes in this group",
														Type:        "integer",
													},
													"labels": {
														Description: "Labels for the nodes in this group",
														Type:        "object",
														Items: &apiextensions.JSONSchemaPropsOrArray{
															Schema: &apiextensions.JSONSchemaProps{
																Type: "string",
															},
														},
													},
													"kernelConfig": {
														Description: "Optional name of KernelConfig resource for the nodes in this group",
														Type:        "string",
													},
													"memoryConfig": {
														Description: "Optional name of MemoryConfig resource for the nodes in this group",
														Type:        "string",
													},
												},
											},
										},
									},
									"wans": {
										Description: "WANs for reaching this cluster",
										Type:        "array",
										Items: &apiextensions.JSONSchemaPropsOrArray{
											Schema: &apiextensions.JSONSchemaProps{
												Type: "string",
											},
										},
									},
								},
							},
							"status": {
								Type: "object",
								Properties: map[string]apiextensions.JSONSchemaProps{
									"kubeconfigUrl": {
										Description: "Full URL of kubeconfig (YAML file)",
										Type:        "string",
									},
									"context": {
										Description: "Name of context to use in kubeconfig",
										Type:        "string",
									},
								},
							},
						},
					},
				},
			},
		},
	},
}

func ClusterToARD(cluster *Cluster) ard.StringMap {
	map_ := make(ard.StringMap)
	map_["name"] = cluster.Name
	var nodes []ard.StringMap
	if cluster.Spec.Nodes != nil {
		for _, node := range cluster.Spec.Nodes {
			nodes = append(nodes, ard.StringMap{
				"count":        node.Count,
				"labels":       node.Labels,
				"kernelConfig": node.KernelConfig,
				"memoryConfig": node.MemoryConfig,
			})
		}
	}
	map_["nodes"] = nodes
	map_["wans"] = cluster.Spec.WANs
	map_["kubeconfigUrl"] = cluster.Status.KubeConfigURL
	map_["context"] = cluster.Status.Context
	return map_
}

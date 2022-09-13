package v1alpha1

import (
	"fmt"

	group "github.com/tliron/planter/resources/planter.nephio.org"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	NetworkKind     = "Network"
	NetworkListKind = "NetworkList"

	NetworkSingular = "network"
	NetworkPlural   = "networks"
)

//
// NetworkCustomResourceDefinition
//

// See: assets/kubernetes/custom-resource-definitions.yaml

var NetworkResourcesName = fmt.Sprintf("%s.%s", NetworkPlural, group.GroupName)

var NetworkCustomResourceDefinition = apiextensions.CustomResourceDefinition{
	ObjectMeta: meta.ObjectMeta{
		Name: NetworkResourcesName,
	},
	Spec: apiextensions.CustomResourceDefinitionSpec{
		Group: group.GroupName,
		Names: apiextensions.CustomResourceDefinitionNames{
			Singular: NetworkSingular,
			Plural:   NetworkPlural,
			Kind:     NetworkKind,
			ListKind: NetworkListKind,
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
						Description: "Planter network",
						Type:        "object",
						Required:    []string{"spec"},
						Properties: map[string]apiextensions.JSONSchemaProps{
							"spec": {
								Type: "object",
								Properties: map[string]apiextensions.JSONSchemaProps{
									"provider": {
										Description: "Network provider",
										Type:        "string",
									},
								},
							},
							"status": {
								Type:       "object",
								Properties: map[string]apiextensions.JSONSchemaProps{},
							},
						},
					},
				},
			},
		},
	},
}

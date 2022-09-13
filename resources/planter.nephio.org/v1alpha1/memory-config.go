package v1alpha1

import (
	"fmt"

	group "github.com/tliron/planter/resources/planter.nephio.org"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	MemoryConfigKind     = "MemoryConfig"
	MemoryConfigListKind = "MemoryConfigList"

	MemoryConfigSingular = "memoryconfig"
	MemoryConfigPlural   = "memoryconfigs"
)

//
// MemoryConfigCustomResourceDefinition
//

// See: assets/kubernetes/custom-resource-definitions.yaml

var MemoryConfigResourcesName = fmt.Sprintf("%s.%s", MemoryConfigPlural, group.GroupName)

var MemoryConfigCustomResourceDefinition = apiextensions.CustomResourceDefinition{
	ObjectMeta: meta.ObjectMeta{
		Name: MemoryConfigResourcesName,
	},
	Spec: apiextensions.CustomResourceDefinitionSpec{
		Group: group.GroupName,
		Names: apiextensions.CustomResourceDefinitionNames{
			Singular: MemoryConfigSingular,
			Plural:   MemoryConfigPlural,
			Kind:     MemoryConfigKind,
			ListKind: MemoryConfigListKind,
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
						Description: "Planter memory config",
						Type:        "object",
						Required:    []string{"spec"},
						Properties: map[string]apiextensions.JSONSchemaProps{
							"spec": {
								Type: "object",
								Properties: map[string]apiextensions.JSONSchemaProps{
									"numa": {
										Description: "Whether NUMA should be enabled",
										Type:        "boolean",
									},
									"hugePages": {
										Description: "Size of huge pages to be used",
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

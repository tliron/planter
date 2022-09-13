package v1alpha1

import (
	"fmt"

	group "github.com/tliron/planter/resources/planter.nephio.org"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	KernelConfigKind     = "KernelConfig"
	KernelConfigListKind = "KernelConfigList"

	KernelConfigSingular = "kernelconfig"
	KernelConfigPlural   = "kernelconfigs"
)

//
// KernelConfigCustomResourceDefinition
//

// See: assets/kubernetes/custom-resource-definitions.yaml

var KernelConfigResourcesName = fmt.Sprintf("%s.%s", KernelConfigPlural, group.GroupName)

var KernelConfigCustomResourceDefinition = apiextensions.CustomResourceDefinition{
	ObjectMeta: meta.ObjectMeta{
		Name: KernelConfigResourcesName,
	},
	Spec: apiextensions.CustomResourceDefinitionSpec{
		Group: group.GroupName,
		Names: apiextensions.CustomResourceDefinitionNames{
			Singular: KernelConfigSingular,
			Plural:   KernelConfigPlural,
			Kind:     KernelConfigKind,
			ListKind: KernelConfigListKind,
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
						Description: "Planter kernel config",
						Type:        "object",
						Required:    []string{"spec"},
						Properties: map[string]apiextensions.JSONSchemaProps{
							"spec": {
								Type: "object",
								Properties: map[string]apiextensions.JSONSchemaProps{
									"realtime": {
										Description: "Whether a realtime kernel should be used",
										Type:        "boolean",
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

package v1alpha1

import (
	"fmt"

	group "github.com/tliron/planter/resources/planter.nephio.org"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	WANKind     = "WAN"
	WANListKind = "WANList"

	WANSingular = "wan"
	WANPlural   = "wans"
)

//
// WANCustomResourceDefinition
//

// See: assets/kubernetes/custom-resource-definitions.yaml

var WANResourcesName = fmt.Sprintf("%s.%s", WANPlural, group.GroupName)

var WANCustomResourceDefinition = apiextensions.CustomResourceDefinition{
	ObjectMeta: meta.ObjectMeta{
		Name: WANResourcesName,
	},
	Spec: apiextensions.CustomResourceDefinitionSpec{
		Group: group.GroupName,
		Names: apiextensions.CustomResourceDefinitionNames{
			Singular: WANSingular,
			Plural:   WANPlural,
			Kind:     WANKind,
			ListKind: WANListKind,
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
						Description: "Planter WAN",
						Type:        "object",
						Required:    []string{"spec"},
						Properties: map[string]apiextensions.JSONSchemaProps{
							"spec": {
								Type: "object",
								Properties: map[string]apiextensions.JSONSchemaProps{
									"provider": {
										Description: "WAN provider",
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

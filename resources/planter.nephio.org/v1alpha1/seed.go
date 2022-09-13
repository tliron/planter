package v1alpha1

import (
	"fmt"

	"github.com/tliron/kutil/ard"
	group "github.com/tliron/planter/resources/planter.nephio.org"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var SeedGVK = SchemeGroupVersion.WithKind(SeedKind)

const (
	SeedKind     = "Seed"
	SeedListKind = "SeedList"

	SeedSingular  = "seed"
	SeedPlural    = "seeds"
	SeedShortName = "ps" // = Planter Seed
)

//
// Seed
//

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Seed struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`

	Spec   SeedSpec   `json:"spec"`
	Status SeedStatus `json:"status"`
}

type SeedSpec struct {
	SeedURL string `json:"seedUrl"` // Full URL of seed (YAML file)
	Planted bool   `json:"planted"` // Whether the seed should be planted
}

type SeedStatus struct {
	PlantedPath string `json:"plantedPath"` // Planted seed path or empty if not planted
}

//
// SeedList
//

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SeedList struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata"`

	Items []Seed `json:"items"`
}

//
// SeedCustomResourceDefinition
//

// See: assets/kubernetes/custom-resource-definitions.yaml

var SeedResourcesName = fmt.Sprintf("%s.%s", SeedPlural, group.GroupName)

var SeedCustomResourceDefinition = apiextensions.CustomResourceDefinition{
	ObjectMeta: meta.ObjectMeta{
		Name: SeedResourcesName,
	},
	Spec: apiextensions.CustomResourceDefinitionSpec{
		Group: group.GroupName,
		Names: apiextensions.CustomResourceDefinitionNames{
			Singular: SeedSingular,
			Plural:   SeedPlural,
			Kind:     SeedKind,
			ListKind: SeedListKind,
			ShortNames: []string{
				SeedShortName,
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
						Description: "Planter seed",
						Type:        "object",
						Required:    []string{"spec"},
						Properties: map[string]apiextensions.JSONSchemaProps{
							"spec": {
								Type:     "object",
								Required: []string{"seedUrl"},
								Properties: map[string]apiextensions.JSONSchemaProps{
									"seedUrl": {
										Description: "Full URL of seed (YAML file)",
										Type:        "string",
									},
									"planted": {
										Description: "Whether the seed should be planted",
										Type:        "boolean",
									},
								},
							},
							"status": {
								Type: "object",
								Properties: map[string]apiextensions.JSONSchemaProps{
									"plantedPath": {
										Description: "Planted seed path or empty if not planted",
										Type:        "string",
									},
								},
							},
						},
					},
				},
				AdditionalPrinterColumns: []apiextensions.CustomResourceColumnDefinition{
					{
						Name:        "PlantedPath",
						Description: "Planted seed path or empty if not planted",
						Type:        "string",
						JSONPath:    ".status.plantedPath",
					},
				},
			},
		},
	},
}

func SeedToARD(service *Seed) ard.StringMap {
	map_ := make(ard.StringMap)
	map_["name"] = service.Name
	map_["url"] = service.Spec.SeedURL
	map_["planted"] = service.Spec.Planted
	map_["plantedPath"] = service.Status.PlantedPath
	return map_
}

package client

import (
	contextpkg "context"

	"github.com/tliron/commonlog"
	planterpkg "github.com/tliron/planter/apis/clientset/versioned"
	apiextensionspkg "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	kubernetespkg "k8s.io/client-go/kubernetes"
	restpkg "k8s.io/client-go/rest"
)

//
// Client
//

type Client struct {
	Kubernetes    kubernetespkg.Interface
	APIExtensions apiextensionspkg.Interface
	Planter       planterpkg.Interface
	REST          restpkg.Interface
	Config        *restpkg.Config

	ClusterRole       string
	Namespace         string
	NamePrefix        string
	PartOf            string
	ManagedBy         string
	OperatorImageName string
	CachePath         string

	Context contextpkg.Context
	Log     commonlog.Logger
}

func NewClient(kubernetes kubernetespkg.Interface, apiExtensions apiextensionspkg.Interface, planter planterpkg.Interface, rest restpkg.Interface, config *restpkg.Config, context contextpkg.Context, clusterRole string, namespace string, namePrefix string, partOf string, managedBy string, operatorImageName string, cachePath string, logName string) *Client {
	return &Client{
		Kubernetes:        kubernetes,
		APIExtensions:     apiExtensions,
		Planter:           planter,
		REST:              rest,
		Config:            config,
		ClusterRole:       clusterRole,
		Namespace:         namespace,
		NamePrefix:        namePrefix,
		PartOf:            partOf,
		ManagedBy:         managedBy,
		OperatorImageName: operatorImageName,
		CachePath:         cachePath,
		Context:           context,
		Log:               commonlog.GetLogger(logName),
	}
}

package commands

import (
	contextpkg "context"
	"fmt"

	kubernetesutil "github.com/tliron/kutil/kubernetes"
	"github.com/tliron/kutil/util"
	planterpkg "github.com/tliron/planter/apis/clientset/versioned"
	clientpkg "github.com/tliron/planter/client"
	"github.com/tliron/planter/controller"
	apiextensionspkg "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	kubernetespkg "k8s.io/client-go/kubernetes"
	restpkg "k8s.io/client-go/rest"
)

//
// Client
//

type Client struct {
	Config     *restpkg.Config
	Kubernetes kubernetespkg.Interface
	REST       restpkg.Interface
	Context    contextpkg.Context
	Namespace  string
}

func NewClient() *Client {
	config, err := kubernetesutil.NewConfigFromFlags(masterUrl, kubeconfigPath, kubeconfigContext, log)
	util.FailOnError(err)

	kubernetes, err := kubernetespkg.NewForConfig(config)
	util.FailOnError(err)

	namespace_ := planterNamespace
	if namespace_ == "" {
		namespace_ = namespace
	}
	if namespace_ == "" {
		if namespace__, ok := kubernetesutil.GetConfiguredNamespace(kubeconfigPath, kubeconfigContext); ok {
			namespace_ = namespace__
		}
		if namespace_ == "" {
			util.Fail("could not discover namespace and \"--namespace\" not provided")
		}
	}

	return &Client{
		Config:     config,
		Kubernetes: kubernetes,
		REST:       kubernetes.CoreV1().RESTClient(),
		Context:    context,
		Namespace:  namespace_,
	}
}

func (self *Client) Planter() *clientpkg.Client {
	apiExtensions, err := apiextensionspkg.NewForConfig(self.Config)
	util.FailOnError(err)

	planter, err := planterpkg.NewForConfig(self.Config)
	util.FailOnError(err)

	return clientpkg.NewClient(
		self.Kubernetes,
		apiExtensions,
		planter,
		self.REST,
		self.Config,
		self.Context,
		clusterRole,
		self.Namespace,
		controller.NamePrefix,
		controller.PartOf,
		controller.ManagedBy,
		controller.OperatorImageName,
		controller.CacheDirectory,
		fmt.Sprintf("%s.client", toolName),
	)
}

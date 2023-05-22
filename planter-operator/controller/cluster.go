package controller

import (
	contextpkg "context"
	"strings"

	kubernetesutil "github.com/tliron/kutil/kubernetes"
	"github.com/tliron/kutil/util"
	dynamicpkg "k8s.io/client-go/dynamic"
	kubernetespkg "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

//
// ClusterClient
//

type ClusterClient struct {
	kubernetes kubernetespkg.Interface
	dynamic    dynamicpkg.Interface
	context    contextpkg.Context
}

func (self *Controller) NewClusterClient(namespace string, clusterName string) (*ClusterClient, error) {
	configContent, err := self.Client.GetClusterKubeConfig(namespace, clusterName)
	if err != nil {
		return nil, err
	}

	clientConfig, err := clientcmd.NewClientConfigFromBytes(util.StringToBytes(configContent))
	if err != nil {
		return nil, err
	}

	config, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	kubernetes, err := kubernetespkg.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	dynamic, err := dynamicpkg.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &ClusterClient{
		kubernetes: kubernetes,
		dynamic:    dynamic,
		context:    self.Context,
	}, nil
}

func (self *ClusterClient) Dynamic(namespace string) *kubernetesutil.Dynamic {
	return kubernetesutil.NewDynamic("planter", self.dynamic, self.kubernetes.Discovery(), namespace, self.context)
}

//
// ClusterClients
//

type ClusterClients struct {
	controller *Controller
	clients    map[string]*ClusterClient
}

func (self *Controller) NewClusterClients() *ClusterClients {
	return &ClusterClients{
		controller: self,
		clients:    make(map[string]*ClusterClient),
	}
}

func (self *ClusterClients) Get(clusterName string) (*ClusterClient, error) {
	if client, ok := self.clients[clusterName]; ok {
		return client, nil
	} else {
		split := strings.Split(clusterName, "|")
		namespace, name := split[0], split[1]
		if client, err := self.controller.NewClusterClient(namespace, name); err == nil {
			self.clients[clusterName] = client
			return client, nil
		} else {
			return nil, err
		}
	}
}

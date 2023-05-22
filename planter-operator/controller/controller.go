package controller

import (
	contextpkg "context"
	"fmt"
	"time"

	"github.com/tliron/commonlog"
	kubernetesutil "github.com/tliron/kutil/kubernetes"
	planterpkg "github.com/tliron/planter/apis/clientset/versioned"
	planterinformers "github.com/tliron/planter/apis/informers/externalversions"
	planterlisters "github.com/tliron/planter/apis/listers/planter.nephio.org/v1alpha1"
	clientpkg "github.com/tliron/planter/client"
	planterresources "github.com/tliron/planter/resources/planter.nephio.org/v1alpha1"
	apiextensionspkg "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	dynamicpkg "k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	restpkg "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
)

//
// Controller
//

type Controller struct {
	Config      *restpkg.Config
	Dynamic     *kubernetesutil.Dynamic
	Kubernetes  kubernetes.Interface
	Planter     planterpkg.Interface
	Client      *clientpkg.Client
	CachePath   string
	StopChannel <-chan struct{}

	Processors *kubernetesutil.Processors
	Events     record.EventRecorder

	KubernetesInformerFactory informers.SharedInformerFactory
	PlanterInformerFactory    planterinformers.SharedInformerFactory

	Seeds planterlisters.SeedLister

	Context contextpkg.Context
	Log     commonlog.Logger
}

func NewController(context contextpkg.Context, toolName string, clusterRole string, namespace string, dynamic dynamicpkg.Interface, kubernetes kubernetes.Interface, apiExtensions apiextensionspkg.Interface, planter planterpkg.Interface, config *restpkg.Config, cachePath string, informerResyncPeriod time.Duration, stopChannel <-chan struct{}) *Controller {
	log := commonlog.GetLoggerf("%s.controller", toolName)

	self := Controller{
		Config:      config,
		Dynamic:     kubernetesutil.NewDynamic(toolName, dynamic, kubernetes.Discovery(), namespace, context),
		Kubernetes:  kubernetes,
		Planter:     planter,
		CachePath:   cachePath,
		StopChannel: stopChannel,
		Processors:  kubernetesutil.NewProcessors(toolName),
		Events:      kubernetesutil.CreateEventRecorder(kubernetes, "Planter", log),
		Context:     context,
		Log:         log,
	}

	self.Client = clientpkg.NewClient(
		kubernetes,
		apiExtensions,
		planter,
		kubernetes.CoreV1().RESTClient(),
		config,
		context,
		clusterRole,
		namespace,
		NamePrefix,
		PartOf,
		ManagedBy,
		OperatorImageName,
		CacheDirectory,
		fmt.Sprintf("%s.client", toolName),
	)

	self.KubernetesInformerFactory = informers.NewSharedInformerFactory(kubernetes, informerResyncPeriod)
	self.PlanterInformerFactory = planterinformers.NewSharedInformerFactory(planter, informerResyncPeriod)

	// Informers
	seedInformer := self.PlanterInformerFactory.Planter().V1alpha1().Seeds()

	// Listers
	self.Seeds = seedInformer.Lister()

	// Processors

	processorPeriod := 5 * time.Second

	self.Processors.Add(planterresources.SeedGVK, kubernetesutil.NewProcessor(
		toolName,
		"seeds",
		seedInformer.Informer(),
		processorPeriod,
		func(name string, namespace string) (any, error) {
			return self.Client.GetSeed(namespace, name)
		},
		func(object any) (bool, error) {
			return self.processSeed(object.(*planterresources.Seed))
		},
	))

	return &self
}

func (self *Controller) Run(concurrency uint, startup func()) error {
	defer utilruntime.HandleCrash()

	self.Log.Info("starting informer factories")
	self.KubernetesInformerFactory.Start(self.StopChannel)
	self.PlanterInformerFactory.Start(self.StopChannel)

	self.Log.Info("waiting for processor informer caches to sync")
	utilruntime.HandleError(self.Processors.WaitForCacheSync(self.StopChannel))

	self.Log.Infof("starting processors (concurrency=%d)", concurrency)
	self.Processors.Start(concurrency, self.StopChannel)
	defer self.Processors.ShutDown()

	if startup != nil {
		go startup()
	}

	<-self.StopChannel

	self.Log.Info("shutting down")

	return nil
}

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	versioned "github.com/tliron/planter/apis/clientset/versioned"
	internalinterfaces "github.com/tliron/planter/apis/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/tliron/planter/apis/listers/planter.nephio.org/v1alpha1"
	planternephioorgv1alpha1 "github.com/tliron/planter/resources/planter.nephio.org/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// SeedInformer provides access to a shared informer and lister for
// Seeds.
type SeedInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.SeedLister
}

type seedInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewSeedInformer constructs a new informer for Seed type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewSeedInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredSeedInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredSeedInformer constructs a new informer for Seed type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredSeedInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.PlanterV1alpha1().Seeds(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.PlanterV1alpha1().Seeds(namespace).Watch(context.TODO(), options)
			},
		},
		&planternephioorgv1alpha1.Seed{},
		resyncPeriod,
		indexers,
	)
}

func (f *seedInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredSeedInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *seedInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&planternephioorgv1alpha1.Seed{}, f.defaultInformer)
}

func (f *seedInformer) Lister() v1alpha1.SeedLister {
	return v1alpha1.NewSeedLister(f.Informer().GetIndexer())
}

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"
	json "encoding/json"
	"fmt"

	planternephioorgv1alpha1 "github.com/tliron/planter/apis/applyconfiguration/planter.nephio.org/v1alpha1"
	v1alpha1 "github.com/tliron/planter/resources/planter.nephio.org/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeSeeds implements SeedInterface
type FakeSeeds struct {
	Fake *FakePlanterV1alpha1
	ns   string
}

var seedsResource = v1alpha1.SchemeGroupVersion.WithResource("seeds")

var seedsKind = v1alpha1.SchemeGroupVersion.WithKind("Seed")

// Get takes name of the seed, and returns the corresponding seed object, and an error if there is any.
func (c *FakeSeeds) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Seed, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(seedsResource, c.ns, name), &v1alpha1.Seed{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Seed), err
}

// List takes label and field selectors, and returns the list of Seeds that match those selectors.
func (c *FakeSeeds) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.SeedList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(seedsResource, seedsKind, c.ns, opts), &v1alpha1.SeedList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.SeedList{ListMeta: obj.(*v1alpha1.SeedList).ListMeta}
	for _, item := range obj.(*v1alpha1.SeedList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested seeds.
func (c *FakeSeeds) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(seedsResource, c.ns, opts))

}

// Create takes the representation of a seed and creates it.  Returns the server's representation of the seed, and an error, if there is any.
func (c *FakeSeeds) Create(ctx context.Context, seed *v1alpha1.Seed, opts v1.CreateOptions) (result *v1alpha1.Seed, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(seedsResource, c.ns, seed), &v1alpha1.Seed{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Seed), err
}

// Update takes the representation of a seed and updates it. Returns the server's representation of the seed, and an error, if there is any.
func (c *FakeSeeds) Update(ctx context.Context, seed *v1alpha1.Seed, opts v1.UpdateOptions) (result *v1alpha1.Seed, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(seedsResource, c.ns, seed), &v1alpha1.Seed{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Seed), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeSeeds) UpdateStatus(ctx context.Context, seed *v1alpha1.Seed, opts v1.UpdateOptions) (*v1alpha1.Seed, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(seedsResource, "status", c.ns, seed), &v1alpha1.Seed{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Seed), err
}

// Delete takes name of the seed and deletes it. Returns an error if one occurs.
func (c *FakeSeeds) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(seedsResource, c.ns, name, opts), &v1alpha1.Seed{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeSeeds) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(seedsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.SeedList{})
	return err
}

// Patch applies the patch and returns the patched seed.
func (c *FakeSeeds) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Seed, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(seedsResource, c.ns, name, pt, data, subresources...), &v1alpha1.Seed{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Seed), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied seed.
func (c *FakeSeeds) Apply(ctx context.Context, seed *planternephioorgv1alpha1.SeedApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Seed, err error) {
	if seed == nil {
		return nil, fmt.Errorf("seed provided to Apply must not be nil")
	}
	data, err := json.Marshal(seed)
	if err != nil {
		return nil, err
	}
	name := seed.Name
	if name == nil {
		return nil, fmt.Errorf("seed.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(seedsResource, c.ns, *name, types.ApplyPatchType, data), &v1alpha1.Seed{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Seed), err
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *FakeSeeds) ApplyStatus(ctx context.Context, seed *planternephioorgv1alpha1.SeedApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Seed, err error) {
	if seed == nil {
		return nil, fmt.Errorf("seed provided to Apply must not be nil")
	}
	data, err := json.Marshal(seed)
	if err != nil {
		return nil, err
	}
	name := seed.Name
	if name == nil {
		return nil, fmt.Errorf("seed.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(seedsResource, c.ns, *name, types.ApplyPatchType, data, "status"), &v1alpha1.Seed{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Seed), err
}

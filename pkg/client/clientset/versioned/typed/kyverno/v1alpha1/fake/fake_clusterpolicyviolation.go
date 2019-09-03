/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/nirmata/kyverno/pkg/api/kyverno/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeClusterPolicyViolations implements ClusterPolicyViolationInterface
type FakeClusterPolicyViolations struct {
	Fake *FakeKyvernoV1alpha1
}

var clusterpolicyviolationsResource = schema.GroupVersionResource{Group: "kyverno.io", Version: "v1alpha1", Resource: "clusterpolicyviolations"}

var clusterpolicyviolationsKind = schema.GroupVersionKind{Group: "kyverno.io", Version: "v1alpha1", Kind: "ClusterPolicyViolation"}

// Get takes name of the clusterPolicyViolation, and returns the corresponding clusterPolicyViolation object, and an error if there is any.
func (c *FakeClusterPolicyViolations) Get(name string, options v1.GetOptions) (result *v1alpha1.ClusterPolicyViolation, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(clusterpolicyviolationsResource, name), &v1alpha1.ClusterPolicyViolation{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterPolicyViolation), err
}

// List takes label and field selectors, and returns the list of ClusterPolicyViolations that match those selectors.
func (c *FakeClusterPolicyViolations) List(opts v1.ListOptions) (result *v1alpha1.ClusterPolicyViolationList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(clusterpolicyviolationsResource, clusterpolicyviolationsKind, opts), &v1alpha1.ClusterPolicyViolationList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ClusterPolicyViolationList{ListMeta: obj.(*v1alpha1.ClusterPolicyViolationList).ListMeta}
	for _, item := range obj.(*v1alpha1.ClusterPolicyViolationList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested clusterPolicyViolations.
func (c *FakeClusterPolicyViolations) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(clusterpolicyviolationsResource, opts))
}

// Create takes the representation of a clusterPolicyViolation and creates it.  Returns the server's representation of the clusterPolicyViolation, and an error, if there is any.
func (c *FakeClusterPolicyViolations) Create(clusterPolicyViolation *v1alpha1.ClusterPolicyViolation) (result *v1alpha1.ClusterPolicyViolation, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(clusterpolicyviolationsResource, clusterPolicyViolation), &v1alpha1.ClusterPolicyViolation{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterPolicyViolation), err
}

// Update takes the representation of a clusterPolicyViolation and updates it. Returns the server's representation of the clusterPolicyViolation, and an error, if there is any.
func (c *FakeClusterPolicyViolations) Update(clusterPolicyViolation *v1alpha1.ClusterPolicyViolation) (result *v1alpha1.ClusterPolicyViolation, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(clusterpolicyviolationsResource, clusterPolicyViolation), &v1alpha1.ClusterPolicyViolation{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterPolicyViolation), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeClusterPolicyViolations) UpdateStatus(clusterPolicyViolation *v1alpha1.ClusterPolicyViolation) (*v1alpha1.ClusterPolicyViolation, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(clusterpolicyviolationsResource, "status", clusterPolicyViolation), &v1alpha1.ClusterPolicyViolation{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterPolicyViolation), err
}

// Delete takes name of the clusterPolicyViolation and deletes it. Returns an error if one occurs.
func (c *FakeClusterPolicyViolations) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(clusterpolicyviolationsResource, name), &v1alpha1.ClusterPolicyViolation{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeClusterPolicyViolations) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(clusterpolicyviolationsResource, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.ClusterPolicyViolationList{})
	return err
}

// Patch applies the patch and returns the patched clusterPolicyViolation.
func (c *FakeClusterPolicyViolations) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.ClusterPolicyViolation, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(clusterpolicyviolationsResource, name, pt, data, subresources...), &v1alpha1.ClusterPolicyViolation{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterPolicyViolation), err
}

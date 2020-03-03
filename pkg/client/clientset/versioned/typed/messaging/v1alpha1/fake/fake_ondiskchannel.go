/*
Copyright 2020 The Knative Authors

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
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
	v1alpha1 "knative.dev/eventing/pkg/apis/messaging/v1alpha1"
)

// FakeOnDiskChannels implements OnDiskChannelInterface
type FakeOnDiskChannels struct {
	Fake *FakeMessagingV1alpha1
	ns   string
}

var ondiskchannelsResource = schema.GroupVersionResource{Group: "messaging.knative.dev", Version: "v1alpha1", Resource: "ondiskchannels"}

var ondiskchannelsKind = schema.GroupVersionKind{Group: "messaging.knative.dev", Version: "v1alpha1", Kind: "OnDiskChannel"}

// Get takes name of the onDiskChannel, and returns the corresponding onDiskChannel object, and an error if there is any.
func (c *FakeOnDiskChannels) Get(name string, options v1.GetOptions) (result *v1alpha1.OnDiskChannel, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(ondiskchannelsResource, c.ns, name), &v1alpha1.OnDiskChannel{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.OnDiskChannel), err
}

// List takes label and field selectors, and returns the list of OnDiskChannels that match those selectors.
func (c *FakeOnDiskChannels) List(opts v1.ListOptions) (result *v1alpha1.OnDiskChannelList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(ondiskchannelsResource, ondiskchannelsKind, c.ns, opts), &v1alpha1.OnDiskChannelList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.OnDiskChannelList{ListMeta: obj.(*v1alpha1.OnDiskChannelList).ListMeta}
	for _, item := range obj.(*v1alpha1.OnDiskChannelList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested onDiskChannels.
func (c *FakeOnDiskChannels) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(ondiskchannelsResource, c.ns, opts))

}

// Create takes the representation of a onDiskChannel and creates it.  Returns the server's representation of the onDiskChannel, and an error, if there is any.
func (c *FakeOnDiskChannels) Create(onDiskChannel *v1alpha1.OnDiskChannel) (result *v1alpha1.OnDiskChannel, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(ondiskchannelsResource, c.ns, onDiskChannel), &v1alpha1.OnDiskChannel{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.OnDiskChannel), err
}

// Update takes the representation of a onDiskChannel and updates it. Returns the server's representation of the onDiskChannel, and an error, if there is any.
func (c *FakeOnDiskChannels) Update(onDiskChannel *v1alpha1.OnDiskChannel) (result *v1alpha1.OnDiskChannel, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(ondiskchannelsResource, c.ns, onDiskChannel), &v1alpha1.OnDiskChannel{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.OnDiskChannel), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeOnDiskChannels) UpdateStatus(onDiskChannel *v1alpha1.OnDiskChannel) (*v1alpha1.OnDiskChannel, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(ondiskchannelsResource, "status", c.ns, onDiskChannel), &v1alpha1.OnDiskChannel{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.OnDiskChannel), err
}

// Delete takes name of the onDiskChannel and deletes it. Returns an error if one occurs.
func (c *FakeOnDiskChannels) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(ondiskchannelsResource, c.ns, name), &v1alpha1.OnDiskChannel{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeOnDiskChannels) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(ondiskchannelsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.OnDiskChannelList{})
	return err
}

// Patch applies the patch and returns the patched onDiskChannel.
func (c *FakeOnDiskChannels) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.OnDiskChannel, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(ondiskchannelsResource, c.ns, name, pt, data, subresources...), &v1alpha1.OnDiskChannel{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.OnDiskChannel), err
}

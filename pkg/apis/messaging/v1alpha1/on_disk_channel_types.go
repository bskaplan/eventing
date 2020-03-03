/*
 * Copyright 2019 The Knative Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	eventingduck "knative.dev/eventing/pkg/apis/duck/v1alpha1"
	eventingduckv1beta1 "knative.dev/eventing/pkg/apis/duck/v1beta1"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	duckv1alpha1 "knative.dev/pkg/apis/duck/v1alpha1"
	"knative.dev/pkg/kmeta"
)

// +genclient
// +genreconciler
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OnDiskChannel is a resource representing an in memory channel
type OnDiskChannel struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the desired state of the Channel.
	Spec OnDiskChannelSpec `json:"spec,omitempty"`

	// Status represents the current state of the Channel. This data may be out of
	// date.
	// +optional
	Status OnDiskChannelStatus `json:"status,omitempty"`
}

var (
	// Check that OnDiskChannel can be validated and defaulted.
	_ apis.Validatable = (*OnDiskChannel)(nil)
	_ apis.Defaultable = (*OnDiskChannel)(nil)

	// Check that OnDiskChannel can return its spec untyped.
	_ apis.HasSpec = (*OnDiskChannel)(nil)

	_ runtime.Object = (*OnDiskChannel)(nil)

	// Check that we can create OwnerReferences to an OnDiskChannel.
	_ kmeta.OwnerRefable = (*OnDiskChannel)(nil)
)

// OnDiskChannelSpec defines which subscribers have expressed interest in
// receiving events from this OnDiskChannel.
// arguments for a Channel.
type OnDiskChannelSpec struct {
	// Channel conforms to Duck type Subscribable.
	Subscribable *eventingduck.Subscribable `json:"subscribable,omitempty"`

	// For round tripping (v1beta1 <-> v1alpha1>
	Delivery *eventingduckv1beta1.DeliverySpec `json:"delivery,omitempty"`
}

// ChannelStatus represents the current state of a Channel.
type OnDiskChannelStatus struct {
	// inherits duck/v1 Status, which currently provides:
	// * ObservedGeneration - the 'Generation' of the Service that was last processed by the controller.
	// * Conditions - the latest available observations of a resource's current state.
	duckv1.Status `json:",inline"`

	// OnDiskChannel is Addressable. It currently exposes the endpoint as a
	// fully-qualified DNS name which will distribute traffic over the
	// provided targets from inside the cluster.
	//
	// It generally has the form {channel}.{namespace}.svc.{cluster domain name}
	duckv1alpha1.AddressStatus `json:",inline"`

	// Subscribers is populated with the statuses of each of the Channelable's subscribers.
	eventingduck.SubscribableTypeStatus `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OnDiskChannelList is a collection of in-memory channels.
type OnDiskChannelList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OnDiskChannel `json:"items"`
}

// GetGroupVersionKind returns GroupVersionKind for OnDiskChannels
func (imc *OnDiskChannel) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("OnDiskChannel")
}

// GetUntypedSpec returns the spec of the OnDiskChannel.
func (i *OnDiskChannel) GetUntypedSpec() interface{} {
	return i.Spec
}

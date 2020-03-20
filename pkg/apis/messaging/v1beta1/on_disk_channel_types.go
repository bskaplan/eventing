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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	eventingduckv1beta1 "knative.dev/eventing/pkg/apis/duck/v1beta1"
	"knative.dev/pkg/apis"
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
	// Channel conforms to Duck type Channelable.
	eventingduckv1beta1.ChannelableSpec `json:",inline"`
}

// ChannelStatus represents the current state of a Channel.
type OnDiskChannelStatus struct {
	// Channel conforms to Duck type Channelable.
	eventingduckv1beta1.ChannelableStatus `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OnDiskChannelList is a collection of in-memory channels.
type OnDiskChannelList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OnDiskChannel `json:"items"`
}


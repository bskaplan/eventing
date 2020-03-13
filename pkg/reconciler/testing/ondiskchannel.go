/*
Copyright 2019 The Knative Authors

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

package testing

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/types"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	duckv1alpha1 "knative.dev/eventing/pkg/apis/duck/v1alpha1"
	"knative.dev/eventing/pkg/apis/eventing"
	"knative.dev/eventing/pkg/apis/messaging/v1alpha1"
	"knative.dev/pkg/apis"
)

// OnDiskChannelOption enables further configuration of a OnDiskChannel.
type OnDiskChannelOption func(*v1alpha1.OnDiskChannel)

// NewOnDiskChannel creates an OnDiskChannel with OnDiskChannelOptions.
func NewOnDiskChannel(name, namespace string, odcopt ...OnDiskChannelOption) *v1alpha1.OnDiskChannel {
	odc := &v1alpha1.OnDiskChannel{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.OnDiskChannelSpec{},
	}
	for _, opt := range odcopt {
		opt(odc)
	}
	odc.SetDefaults(context.Background())
	return odc
}

func WithInitOnDiskChannelConditions(odc *v1alpha1.OnDiskChannel) {
	odc.Status.InitializeConditions()
}

func WithOnDiskChannelGeneration(gen int64) OnDiskChannelOption {
	return func(s *v1alpha1.OnDiskChannel) {
		s.Generation = gen
	}
}

func WithOnDiskChannelStatusObservedGeneration(gen int64) OnDiskChannelOption {
	return func(s *v1alpha1.OnDiskChannel) {
		s.Status.ObservedGeneration = gen
	}
}

func WithOnDiskChannelDeleted(odc *v1alpha1.OnDiskChannel) {
	deleteTime := metav1.NewTime(time.Unix(1e9, 0))
	odc.ObjectMeta.SetDeletionTimestamp(&deleteTime)
}

func WithOnDiskChannelSubscribers(subscribers []duckv1alpha1.SubscriberSpec) OnDiskChannelOption {
	return func(odc *v1alpha1.OnDiskChannel) {
		odc.Spec.Subscribable = &duckv1alpha1.Subscribable{Subscribers: subscribers}
	}
}

func WithOnDiskChannelDeploymentFailed(reason, message string) OnDiskChannelOption {
	return func(odc *v1alpha1.OnDiskChannel) {
		odc.Status.MarkDispatcherFailed(reason, message)
	}
}

func WithOnDiskChannelDeploymentUnknown(reason, message string) OnDiskChannelOption {
	return func(odc *v1alpha1.OnDiskChannel) {
		odc.Status.MarkDispatcherUnknown(reason, message)
	}
}

func WithOnDiskChannelDeploymentReady() OnDiskChannelOption {
	return func(odc *v1alpha1.OnDiskChannel) {
		odc.Status.PropagateDispatcherStatus(&appsv1.StatefulSetStatus{Conditions: []appsv1.StatefulSetCondition{{Type: appsv1.DeploymentAvailable, Status: corev1.ConditionTrue}}})
	}
}

func WithOnDiskChannelServicetNotReady(reason, message string) OnDiskChannelOption {
	return func(odc *v1alpha1.OnDiskChannel) {
		odc.Status.MarkServiceFailed(reason, message)
	}
}

func WithOnDiskChannelServiceReady() OnDiskChannelOption {
	return func(odc *v1alpha1.OnDiskChannel) {
		odc.Status.MarkServiceTrue()
	}
}

func WithOnDiskChannelChannelServiceNotReady(reason, message string) OnDiskChannelOption {
	return func(odc *v1alpha1.OnDiskChannel) {
		odc.Status.MarkChannelServiceFailed(reason, message)
	}
}

func WithOnDiskChannelChannelServiceReady() OnDiskChannelOption {
	return func(odc *v1alpha1.OnDiskChannel) {
		odc.Status.MarkChannelServiceTrue()
	}
}

func WithOnDiskChannelEndpointsNotReady(reason, message string) OnDiskChannelOption {
	return func(odc *v1alpha1.OnDiskChannel) {
		odc.Status.MarkEndpointsFailed(reason, message)
	}
}

func WithOnDiskChannelEndpointsReady() OnDiskChannelOption {
	return func(odc *v1alpha1.OnDiskChannel) {
		odc.Status.MarkEndpointsTrue()
	}
}

func WithOnDiskChannelAddress(a string) OnDiskChannelOption {
	return func(odc *v1alpha1.OnDiskChannel) {
		odc.Status.SetAddress(&apis.URL{
			Scheme: "http",
			Host:   a,
		})
	}
}

func WithOnDiskChannelReady(host string) OnDiskChannelOption {
	return func(odc *v1alpha1.OnDiskChannel) {
		odc.Status.SetAddress(&apis.URL{
			Scheme: "http",
			Host:   host,
		})
		odc.Status.MarkChannelServiceTrue()
		odc.Status.MarkEndpointsTrue()
		odc.Status.MarkServiceTrue()
	}
}

func WithOnDiskChannelReadySubscriber(uid string) OnDiskChannelOption {
	return WithOnDiskChannelReadySubscriberAndGeneration(uid, 0)
}

func WithOnDiskChannelReadySubscriberAndGeneration(uid string, observedGeneration int64) OnDiskChannelOption {
	return func(c *v1alpha1.OnDiskChannel) {
		if c.Status.GetSubscribableTypeStatus() == nil { // Both the SubscribableStatus fields are nil
			c.Status.SetSubscribableTypeStatus(duckv1alpha1.SubscribableStatus{})
		}
		c.Status.SubscribableTypeStatus.AddSubscriberToSubscribableStatus(duckv1alpha1.SubscriberStatus{
			UID:                types.UID(uid),
			ObservedGeneration: observedGeneration,
			Ready:              corev1.ConditionTrue,
		})
	}
}

func WithOnDiskChannelStatusSubscribers(subscriberStatuses []duckv1alpha1.SubscriberStatus) OnDiskChannelOption {
	return func(odc *v1alpha1.OnDiskChannel) {
		odc.Status.SetSubscribableTypeStatus(duckv1alpha1.SubscribableStatus{
			Subscribers: subscriberStatuses})
	}
}

func WithOnDiskScopeAnnotation(value string) OnDiskChannelOption {
	return func(odc *v1alpha1.OnDiskChannel) {
		if odc.Annotations == nil {
			odc.Annotations = make(map[string]string)
		}
		odc.Annotations[eventing.ScopeAnnotationKey] = value
	}
}

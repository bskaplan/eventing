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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	time "time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	messagingv1alpha1 "knative.dev/eventing/pkg/apis/messaging/v1alpha1"
	versioned "knative.dev/eventing/pkg/client/clientset/versioned"
	internalinterfaces "knative.dev/eventing/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "knative.dev/eventing/pkg/client/listers/messaging/v1alpha1"
)

// OnDiskChannelInformer provides access to a shared informer and lister for
// OnDiskChannels.
type OnDiskChannelInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.OnDiskChannelLister
}

type onDiskChannelInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewOnDiskChannelInformer constructs a new informer for OnDiskChannel type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewOnDiskChannelInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredOnDiskChannelInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredOnDiskChannelInformer constructs a new informer for OnDiskChannel type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredOnDiskChannelInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.MessagingV1alpha1().OnDiskChannels(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.MessagingV1alpha1().OnDiskChannels(namespace).Watch(options)
			},
		},
		&messagingv1alpha1.OnDiskChannel{},
		resyncPeriod,
		indexers,
	)
}

func (f *onDiskChannelInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredOnDiskChannelInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *onDiskChannelInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&messagingv1alpha1.OnDiskChannel{}, f.defaultInformer)
}

func (f *onDiskChannelInformer) Lister() v1alpha1.OnDiskChannelLister {
	return v1alpha1.NewOnDiskChannelLister(f.Informer().GetIndexer())
}

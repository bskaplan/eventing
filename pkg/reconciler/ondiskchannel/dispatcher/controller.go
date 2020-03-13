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

package dispatcher

import (
	"context"
	"time"

	"go.uber.org/zap"
	"k8s.io/client-go/tools/cache"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/injection"
	pkgreconciler "knative.dev/pkg/reconciler"
	tracingconfig "knative.dev/pkg/tracing/config"

	"knative.dev/eventing/pkg/apis/eventing"
	"knative.dev/eventing/pkg/channel"
	"knative.dev/eventing/pkg/channel/swappable"
	ondiskchannelinformer "knative.dev/eventing/pkg/client/injection/informers/messaging/v1alpha1/ondiskchannel"
	"knative.dev/eventing/pkg/ondiskchannel"
	"knative.dev/eventing/pkg/reconciler"
	"knative.dev/eventing/pkg/tracing"
)

const (
	// ReconcilerName is the name of the reconciler.
	ReconcilerName = "OnDiskChannels"

	// controllerAgentName is the string used by this controller to identify
	// itself when creating events.
	controllerAgentName = "on-disk-channel-dispatcher"

	readTimeout  = 15 * time.Minute
	writeTimeout = 15 * time.Minute
	port         = 8080
)

// NewController initializes the controller and is called by the generated code.
// Registers event handlers to enqueue events.
func NewController(
	ctx context.Context,
	cmw configmap.Watcher,
) *controller.Impl {
	base := reconciler.NewBase(ctx, controllerAgentName, cmw)

	// Setup trace publishing.
	iw := cmw.(*configmap.InformedWatcher)
	if err := tracing.SetupDynamicPublishing(base.Logger, iw, "odc-dispatcher", tracingconfig.ConfigName); err != nil {
		base.Logger.Fatalw("Error setting up trace publishing", zap.Error(err))
	}

	sh, err := swappable.NewEmptyHandler(base.Logger.Desugar())
	if err != nil {
		base.Logger.Fatal("Error creating swappable.Handler", zap.Error(err))
	}

	args := &ondiskchannel.OnDiskDispatcherArgs{
		Port:         port,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		Handler:      sh,
		Logger:       base.Logger.Desugar(),
	}
	onDiskDispatcher := ondiskchannel.NewDispatcher(args)

	ondiskchannelInformer := ondiskchannelinformer.Get(ctx)
	informer := ondiskchannelInformer.Informer()

	r := &Reconciler{
		Base:                    base,
		dispatcher:              onDiskDispatcher,
		ondiskchannelLister:   ondiskchannelInformer.Lister(),
		ondiskchannelInformer: informer,
	}
	r.impl = controller.NewImpl(r, r.Logger, ReconcilerName)

	// Nothing to filer, enqueue all odcs if configmap updates.
	noopFilter := func(interface{}) bool { return true }
	resyncODCs := configmap.TypeFilter(channel.EventDispatcherConfig{})(func(string, interface{}) {
		r.impl.FilteredGlobalResync(noopFilter, informer)
	})
	// Watch for configmap changes and trigger odc reconciliation by enqueuing imcs.
	configStore := channel.NewEventDispatcherConfigStore(base.Logger, resyncODCs)
	configStore.WatchConfigs(cmw)
	r.configStore = configStore

	r.Logger.Info("Setting up event handlers")

	// Watch for inmemory channels.
	r.ondiskchannelInformer.AddEventHandler(
		cache.FilteringResourceEventHandler{
			FilterFunc: filterWithAnnotation(injection.HasNamespaceScope(ctx)),
			Handler:    controller.HandleAll(r.impl.Enqueue),
		})

	// Start the dispatcher.
	go func() {
		err := onDiskDispatcher.Start(ctx)
		if err != nil {
			r.Logger.Error("Failed stopping onDiskDispatcher.", zap.Error(err))
		}
	}()

	return r.impl
}

func filterWithAnnotation(namespaced bool) func(obj interface{}) bool {
	if namespaced {
		return pkgreconciler.AnnotationFilterFunc(eventing.ScopeAnnotationKey, "namespace", false)
	}
	return pkgreconciler.AnnotationFilterFunc(eventing.ScopeAnnotationKey, "cluster", true)
}

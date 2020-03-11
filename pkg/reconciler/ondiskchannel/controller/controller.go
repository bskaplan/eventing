package controller

import (
	"context"

	"github.com/kelseyhightower/envconfig"
	"k8s.io/client-go/tools/cache"
	"knative.dev/eventing/pkg/reconciler"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/system"

	"knative.dev/eventing/pkg/client/injection/informers/messaging/v1alpha1/ondiskchannel"
	ondiskchannelreconciler "knative.dev/eventing/pkg/client/injection/reconciler/messaging/v1alpha1/ondiskchannel"
	"knative.dev/pkg/client/injection/kube/informers/apps/v1/statefulset"
	"knative.dev/pkg/client/injection/kube/informers/core/v1/endpoints"
	"knative.dev/pkg/client/injection/kube/informers/core/v1/persistentvolumeclaim"
	"knative.dev/pkg/client/injection/kube/informers/core/v1/service"
	"knative.dev/pkg/client/injection/kube/informers/core/v1/serviceaccount"
	"knative.dev/pkg/client/injection/kube/informers/rbac/v1/rolebinding"
)

const (
	// ReconcilerName is the name of the reconciler
	ReconcilerName = "OnDiskChannels"

	// controllerAgentName is the string used by this controller to identify
	// itself when creating events.
	controllerAgentName = "on-disk-channel-controller"

	// TODO: this should be passed in on the env.
	dispatcherName = "odc-dispatcher"
)

type envConfig struct {
	Image string `envconfig:"DISPATCHER_IMAGE" required:"true"`
}

// NewController initializes the controller and is called by the generated code.
// Registers event handlers to enqueue events.
func NewController(
	ctx context.Context,
	cmw configmap.Watcher,
) *controller.Impl {

	ondiskchannelInformer := ondiskchannel.Get(ctx)
	statefulsetInformer := statefulset.Get(ctx)
	serviceInformer := service.Get(ctx)
	endpointsInformer := endpoints.Get(ctx)
	serviceAccountInformer := serviceaccount.Get(ctx)
	roleBindingInformer := rolebinding.Get(ctx)
	persistentVolumeClaimInformer := persistentvolumeclaim.Get(ctx)

	r := &Reconciler{
		Base: reconciler.NewBase(ctx, controllerAgentName, cmw),

		systemNamespace:         system.Namespace(),
		ondiskchannelLister:   	 ondiskchannelInformer.Lister(),
		ondiskchannelInformer:   ondiskchannelInformer.Informer(),
		statefulsetLister:       statefulsetInformer.Lister(),
		serviceLister:           serviceInformer.Lister(),
		endpointsLister:         endpointsInformer.Lister(),
		serviceAccountLister:    serviceAccountInformer.Lister(),
		roleBindingLister:       roleBindingInformer.Lister(),
		persistentVolumeClaimLister: persistentVolumeClaimInformer.Lister(),
	}

	env := &envConfig{}
	if err := envconfig.Process("", env); err != nil {
		r.Logger.Panicf("unable to process in-memory channel's required environment variables: %v", err)
	}

	if env.Image == "" {
		r.Logger.Panic("unable to process in-memory channel's required environment variables (missing DISPATCHER_IMAGE)")
	}

	r.dispatcherImage = env.Image

	impl := ondiskchannelreconciler.NewImpl(ctx, r)

	r.Logger.Info("Setting up event handlers")
	ondiskchannelInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	// Set up watches for dispatcher resources we care about, since any changes to these
	// resources will affect our Channels. So, set up a watch here, that will cause
	// a global Resync for all the channels to take stock of their health when these change.

	// Call GlobalResync on ondiskchannels.
	grCh := func(obj interface{}) {
		impl.GlobalResync(ondiskchannelInformer.Informer())
	}

	statefulsetInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.FilterWithName(dispatcherName),
		Handler:    controller.HandleAll(grCh),
	})
	serviceInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.FilterWithName(dispatcherName),
		Handler:    controller.HandleAll(grCh),
	})
	endpointsInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.FilterWithName(dispatcherName),
		Handler:    controller.HandleAll(grCh),
	})
	serviceAccountInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.FilterWithName(dispatcherName),
		Handler:    controller.HandleAll(grCh),
	})
	roleBindingInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.FilterWithName(dispatcherName),
		Handler:    controller.HandleAll(grCh),
	})
	persistentVolumeClaimInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler(
		FilterFunc: conroller.FilterWithName(dispatcherName),
		Handler: controller.HandleAll(grCh))
	))

	return impl
}
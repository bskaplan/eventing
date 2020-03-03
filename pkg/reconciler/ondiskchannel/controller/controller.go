import (
	"context"

	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	kindreconciler "knative.dev/eventing/pkg/client/injection/reconciler/messaging/v1alpha1/ondiskchannel"
)

func NewController(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
	logger := logging.FromContext(ctx)

	// TODO(you): Access informers

	r := &Reconciler{
		// TODO(you): Pass listers, clients, and other stuff.
	}
	impl := kindreconciler.NewImpl(ctx, r)

	// TODO(you): Set up event handlers.

	return impl
}
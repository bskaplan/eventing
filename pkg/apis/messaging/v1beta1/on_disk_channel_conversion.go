package v1beta1

import (
	"context"
	"fmt"

	"knative.dev/pkg/apis"
)

// ConvertTo implements apis.Convertible
func (source *OnDiskChannel) ConvertTo(ctx context.Context, sink apis.Convertible) error {
	return fmt.Errorf("v1beta1 is the highest known version, got: %T", sink)
}

// ConvertFrom implements apis.Convertible
func (sink *OnDiskChannel) ConvertFrom(ctx context.Context, source apis.Convertible) error {
	return fmt.Errorf("v1beta1 is the highest known version, got: %T", source)
}
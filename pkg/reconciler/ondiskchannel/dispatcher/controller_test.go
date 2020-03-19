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
	"os"
	"testing"

	"knative.dev/pkg/configmap"
	. "knative.dev/pkg/reconciler/testing"

	// Fake injection informers
	_ "knative.dev/eventing/pkg/client/injection/informers/messaging/v1alpha1/ondiskchannel/fake"
)

func TestNew(t *testing.T) {
	ctx, _ := SetupFakeContext(t)

	os.Setenv("SCOPE", "cluster")
	c := NewController(ctx, &configmap.InformedWatcher{})

	if c == nil {
		t.Fatal("Expected NewController to return a non-nil value")
	}
}

func TestNewInNamespace(t *testing.T) {
	ctx, _ := SetupFakeContext(t)

	os.Setenv("SCOPE", "namespace")
	c := NewController(ctx, &configmap.InformedWatcher{})

	if c == nil {
		t.Fatal("Expected NewController to return a non-nil value")
	}
}
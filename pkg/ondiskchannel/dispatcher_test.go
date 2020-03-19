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

package ondiskchannel

import (
	"testing"
	"time"

	"knative.dev/eventing/pkg/channel/swappable"

	logtesting "knative.dev/pkg/logging/testing"
)

func TestNewDispatcher(t *testing.T) {
	logger := logtesting.TestLogger(t).Desugar()
	sh, err := swappable.NewEmptyHandler(logger)

	if err != nil {
		t.Fatalf("Failed to create handler")
	}

	args := &OnDiskDispatcherArgs{
		Port:         8080,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
		Handler:      sh,
		Logger:       logger,
		DbLoc:		  "/tmp/odc-test",
	}

	d := NewDispatcher(args)

	if d == nil {
		t.Fatalf("Failed to create with NewDispatcher")
	}
}
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
	"context"
	"encoding/json"
	"errors"
	"time"

	cloudevents "github.com/cloudevents/sdk-go"
	badger "github.com/dgraph-io/badger"
	"go.uber.org/zap"
	"knative.dev/eventing/pkg/channel/multichannelfanout"
	"knative.dev/eventing/pkg/channel/swappable"
	"knative.dev/eventing/pkg/kncloudevents"
)

const DefaultDbLoc = "/data/events-cache"

type Dispatcher interface {
	UpdateConfig(config *multichannelfanout.Config) error
}

type OnDiskDispatcher struct {
	handler      *swappable.Handler
	ceClient     cloudevents.Client
	writeTimeout time.Duration
	logger       *zap.Logger
	db           *badger.DB
}

type OnDiskDispatcherArgs struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Handler      *swappable.Handler
	Logger       *zap.Logger
	DbLoc        string
}

func (d *OnDiskDispatcher) UpdateConfig(config *multichannelfanout.Config) error {
	return d.handler.UpdateConfig(config)
}
func (d *OnDiskDispatcher) processAndRetry(ctx context.Context, event cloudevents.Event, resp *cloudevents.EventResponse) error {

	// Write to Badger event:
	err := d.db.Update(func(txn *badger.Txn) error {
		data, err := json.Marshal(event)
		if err != nil {
			return err
		}
		return txn.Set([]byte(event.ID()), data)
	})
	err = d.handler.ServeHTTP(ctx, event, resp)

	// attempt one retry
	if err != nil {
		err = d.handler.ServeHTTP(ctx, event, resp)
	}

	if err != nil {
		// Delete from Badger
		err = d.db.Update(func(txn *badger.Txn) error {
			return txn.Delete([]byte(event.ID()))
		})
	}
	return err
}

func (d *OnDiskDispatcher) retryWaiting(ctx context.Context) error {
	// TODO: make this asynchronous with a reasonable limit on concurrent ops
	err := d.db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		var globalErr error
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := item.Key()
			err := item.Value(func(data []byte) error {
				var event cloudevents.Event
				err := json.Unmarshal(data, &event)
				if err != nil {
					txn.Delete(key)
					// TODO: log failed key
					return err
				}
				// we don't do anything with the response
				response := cloudevents.EventResponse{}
				err = d.handler.ServeHTTP(ctx, event, &response)
				// TODO: continue retrying in event of a failure
				txn.Delete(key)
				return err
			})
			if err != nil && globalErr == nil {
				globalErr = err
			}
		}
		return globalErr
	})
	return err
}

// Start starts the inmemory dispatcher's message processing.
// This is a blocking call.
func (d *OnDiskDispatcher) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// start off retrying stuff that's already waiting
	_ = d.retryWaiting(ctx)

	errCh := make(chan error, 1)
	go func() {
		errCh <- d.ceClient.StartReceiver(ctx, d.processAndRetry)
	}()

	// Stop either if the receiver stops (sending to errCh) or if the context Done channel is closed.
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		break
	}

	// Done channel has been closed, we need to gracefully shutdown d.ceClient. The cancel() method will start its
	// shutdown, if it hasn't finished in a reasonable amount of time, just return an error.
	cancel()
	select {
	case err := <-errCh:
		return err
	case <-time.After(d.writeTimeout):
		return errors.New("timeout shutting down ceClient")
	}
}

func NewDispatcher(args *OnDiskDispatcherArgs) *OnDiskDispatcher {
	// TODO set read and write timeouts and port?
	ceClient, err := kncloudevents.NewDefaultClient()
	if err != nil {
		args.Logger.Fatal("failed to create cloudevents client", zap.Error(err))
	}

	dbLoc := args.DbLoc
	if dbLoc == "" {
		dbLoc = DefaultDbLoc
	}
	// TODO: Once knative switches to modules, add support for InMemory mode for testing
	db, err := badger.Open(badger.DefaultOptions(dbLoc))
	if err != nil {
		args.Logger.Fatal("failed to initialize badger db", zap.Error(err))
	}
	dispatcher := &OnDiskDispatcher{
		handler:  args.Handler,
		ceClient: ceClient,
		logger:   args.Logger,
		db:       db,
	}

	return dispatcher
}

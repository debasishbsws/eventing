/*
Copyright 2022 The Knative Authors

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

package subscriptionsapi

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"go.uber.org/zap"
	"knative.dev/pkg/logging"

	"knative.dev/eventing/pkg/eventfilter"
)

type allFilter []eventfilter.Filter

// NewAllFilter returns an event filter which passes if all the contained filters pass
func NewAllFilter(filters ...eventfilter.Filter) eventfilter.Filter {
	return append(allFilter{}, filters...)
}

func (filter allFilter) Filter(ctx context.Context, event cloudevents.Event) eventfilter.FilterResult {
	res := eventfilter.NoFilter
	logging.FromContext(ctx).Debugw("Performing an ALL match ", zap.Any("filters", filter), zap.Any("event", event))
	for _, f := range filter {
		res = res.And(f.Filter(ctx, event))
		// Short circuit to optimize it
		if res == eventfilter.FailFilter {
			return eventfilter.FailFilter
		}
	}
	return res
}

func (filter allFilter) Cleanup() {}

var _ eventfilter.Filter = allFilter{}

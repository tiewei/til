/*
Copyright 2021 TriggerMesh Inc.

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

package channels

import (
	"strconv"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/internal/sdk/k8s"
	"bridgedl/translation"
)

type PubSub struct{}

var (
	_ translation.Decodable    = (*PubSub)(nil)
	_ translation.Translatable = (*PubSub)(nil)
	_ translation.Addressable  = (*PubSub)(nil)
)

// Spec implements translation.Decodable.
func (*PubSub) Spec() hcldec.Spec {
	return &hcldec.AttrSpec{
		Name:     "subscribers",
		Type:     cty.Set(k8s.DestinationCty),
		Required: true,
	}
}

// Manifests implements translation.Translatable.
func (*PubSub) Manifests(id string, config, _ cty.Value) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	ch := k8s.NewChannel(name)
	manifests = append(manifests, ch)

	for i, subscribersIter := 0, config.ElementIterator(); subscribersIter.Next(); i++ {
		subscriberName := name + "-s" + strconv.Itoa(i)

		_, subscriber := subscribersIter.Element()
		subs := k8s.NewSubscription(subscriberName, name, subscriber, cty.NullVal(k8s.DestinationCty))

		manifests = append(manifests, subs)
	}

	return manifests
}

// Address implements translation.Addressable.
func (*PubSub) Address(id string, _, _ cty.Value) cty.Value {
	return k8s.NewDestination(k8s.APIMessaging, "Channel", k8s.RFC1123Name(id))
}

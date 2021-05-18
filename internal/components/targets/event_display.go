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

package targets

import (
	"github.com/zclconf/go-cty/cty"

	"bridgedl/internal/sdk/k8s"
	"bridgedl/translation"
)

type EventDisplay struct{}

var (
	_ translation.Translatable = (*EventDisplay)(nil)
	_ translation.Addressable  = (*EventDisplay)(nil)
)

// Manifests implements translation.Translatable.
func (*EventDisplay) Manifests(id string, _, _ cty.Value) []interface{} {
	var manifests []interface{}

	// Knative v0.23.0
	// https://prow.knative.dev/?job=ci-knative-eventing-auto-release
	const img = "gcr.io/knative-releases/knative.dev/eventing/cmd/event_display" +
		"@sha256:d53673872272dbe8cb044a567e43c4b5e9dfb2bb2bc5b61cf72ff4911ed93979"

	const public = false

	ksvc := k8s.NewKnService(k8s.RFC1123Name(id), img, public)

	return append(manifests, ksvc)
}

// Address implements translation.Addressable.
func (*EventDisplay) Address(id string, _, _ cty.Value) cty.Value {
	return k8s.NewDestination(k8s.APIServing, "Service", k8s.RFC1123Name(id))
}

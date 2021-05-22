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

	"bridgedl/config/globals"
	"bridgedl/internal/sdk/k8s"
	"bridgedl/translation"
)

type Sockeye struct{}

var (
	_ translation.Translatable = (*Sockeye)(nil)
	_ translation.Addressable  = (*Sockeye)(nil)
)

// Manifests implements translation.Translatable.
func (*Sockeye) Manifests(id string, _, _ cty.Value, _ globals.Accessor) []interface{} {
	var manifests []interface{}

	// https://github.com/n3wscott/sockeye/releases
	const img = "docker.io/n3wscott/sockeye:v0.7.0" +
		"@sha256:e603d8494eeacce966e57f8f508e4c4f6bebc71d095e3f5a0a1abaf42c5f0e48"

	const public = true

	ksvc := k8s.NewKnService(k8s.RFC1123Name(id), img, public)

	return append(manifests, ksvc)
}

// Address implements translation.Addressable.
func (*Sockeye) Address(id string, _, _ cty.Value) cty.Value {
	return k8s.NewDestination(k8s.APIServing, "Service", k8s.RFC1123Name(id))
}

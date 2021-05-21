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
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/config/globals"
	"bridgedl/internal/sdk/k8s"
	"bridgedl/translation"
)

type Container struct{}

var (
	_ translation.Decodable    = (*Container)(nil)
	_ translation.Translatable = (*Container)(nil)
	_ translation.Addressable  = (*Container)(nil)
)

// Spec implements translation.Decodable.
func (*Container) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"image": &hcldec.AttrSpec{
			Name:     "image",
			Type:     cty.String,
			Required: true,
		},
		"public": &hcldec.AttrSpec{
			Name:     "public",
			Type:     cty.Bool,
			Required: false,
		},
		"env_vars": &hcldec.ValidateSpec{
			Wrapped: &hcldec.AttrSpec{
				Name:     "env_vars",
				Type:     cty.Map(cty.DynamicPseudoType),
				Required: false,
			},
			Func: validateContainerAttrEnvVars,
		},
		"env_var": &hcldec.BlockObjectSpec{
			TypeName:   "env_var",
			LabelNames: []string{"name"},
			Nested: &hcldec.ValidateSpec{
				Wrapped: &hcldec.AttrSpec{
					Name:     "value",
					Type:     cty.DynamicPseudoType,
					Required: true,
				},
				Func: validateContainerAttrEnvVars,
			},
		},
	}
}

// Manifests implements translation.Translatable.
func (*Container) Manifests(id string, config, eventDst cty.Value, _ globals.Accessor) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	img := config.GetAttr("image").AsString()
	public := config.GetAttr("public").True()

	envVars := make(map[string]cty.Value)
	// Compact syntax (map)
	if v := config.GetAttr("env_vars"); !v.IsNull() {
		for name, value := range v.AsValueMap() {
			envVars[name] = value
		}
	}
	// Verbose syntax (block)
	for name, value := range config.GetAttr("env_var").AsValueMap() {
		envVars[name] = value
	}

	var ksvcOpts []k8s.KnServiceOption
	ksvcOpts = append(ksvcOpts, k8s.EnvVars(envVars))

	ksvc := k8s.NewKnService(name, img, public, ksvcOpts...)
	manifests = append(manifests, ksvc)

	if !eventDst.IsNull() {
		ch := k8s.NewChannel(name)
		subscriber := k8s.NewDestination(k8s.APIServing, "Service", name)
		subs := k8s.NewSubscription(name, name, subscriber, k8s.ReplyDest(eventDst))
		manifests = append(manifests, ch, subs)
	}

	return manifests
}

// Address implements translation.Addressable.
func (*Container) Address(id string, _, eventDst cty.Value, _ globals.Accessor) cty.Value {
	name := k8s.RFC1123Name(id)

	if eventDst.IsNull() {
		return k8s.NewDestination(k8s.APIServing, "Service", name)
	}
	return k8s.NewDestination(k8s.APIMessaging, "Channel", name)
}

func validateContainerAttrEnvVars(val cty.Value) hcl.Diagnostics {
	var diags hcl.Diagnostics

	switch {
	case val.IsNull():
		return diags

	case val.CanIterateElements():
		envVarsIter := val.ElementIterator()
		for envVarsIter.Next() {
			_, v := envVarsIter.Element()
			diags = diags.Extend(validateContainerAttrEnvVar(v))
		}

	default:
		diags = diags.Extend(validateContainerAttrEnvVar(val))
	}

	return diags
}

func validateContainerAttrEnvVar(val cty.Value) hcl.Diagnostics {
	var diags hcl.Diagnostics

	if !(k8s.IsSecretKeySelector(val) || val.Type() == cty.String) {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid value type",
			Detail:   "The value of an environment variable must be either a secret reference or a string.",
		})
	}

	return diags
}

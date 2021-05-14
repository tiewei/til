package routers

import (
	"strconv"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/internal/sdk/k8s"
	"bridgedl/internal/sdk/validation"
	"bridgedl/translation"
)

type ContentBased struct{}

var (
	_ translation.Decodable    = (*ContentBased)(nil)
	_ translation.Translatable = (*ContentBased)(nil)
	_ translation.Addressable  = (*ContentBased)(nil)
)

// Spec implements translation.Decodable.
func (*ContentBased) Spec() hcldec.Spec {
	// NOTE(antoineco): see the following implementation to get a sense of
	// how HCL blocks map to hcldec.Specs and cty.Types:
	// https://pkg.go.dev/github.com/hashicorp/terraform@v0.14.7/configs/configschema#Block.DecoderSpec
	return &hcldec.BlockSetSpec{
		TypeName: "route",
		Nested: &hcldec.ObjectSpec{
			"attributes": &hcldec.ValidateSpec{
				Wrapped: &hcldec.AttrSpec{
					Name:     "attributes",
					Type:     cty.Map(cty.String),
					Required: false,
				},
				Func: validation.ContainsCEContextAttributes,
			},
			"condition": &hcldec.AttrSpec{
				Name:     "condition",
				Type:     cty.String,
				Required: false,
			},
			"to": &hcldec.AttrSpec{
				Name:     "to",
				Type:     k8s.DestinationCty,
				Required: true,
			},
		},
		MinItems: 1,
	}

	/*
		Example of value decoded from the spec above, for a config with
		two "route" blocks:

		v: (set.Set) {
		 vals: (map[int][]interface {}) (len=2) {
		  (int) 1734954449: ([]interface {}) (len=1) {
		   (map[string]interface {}) (len=2) {
		    "attributes": (map[string]interface {}) (len=1) {
		     "type": (string) "corp.acme.my.processing"
		    },
		    "to": (map[string]interface {}) (len=2) {
		     "ref": (map[string]interface {}) (len=3) {
		      "apiVersion": (string) "eventing.knative.dev",
		      "kind": (string) "KafkaSink",
		      "name": (string) "my-kafka-topic"
		     }
		    }
		   }
		  },
		  (int) 3503683627: ([]interface {}) (len=1 cap=1) {
		   (map[string]interface {}) (len=2) {
		    "attributes": (map[string]interface {}) (len=1) {
		     "type": (string) "com.amazon.sqs.message"
		    },
		    "to": (map[string]interface {}) (len=2) {
		     "ref": (map[string]interface {}) (len=3) {
		      "apiVersion": (string) "flow.triggermesh.io/v1alpha1",
		      "kind": (string) "Transformation",
		      "name": (string) "my-transformation"
		     }
		    }
		   }
		  }
		 }
		}
	*/
}

// Manifests implements translation.Translatable.
func (*ContentBased) Manifests(id string, config, _ cty.Value) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	broker := k8s.NewBroker(name)
	manifests = append(manifests, broker)

	i := 0
	for routeIter := config.ElementIterator(); routeIter.Next(); {
		_, route := routeIter.Element()

		routeName := name + "-r" + strconv.Itoa(i)

		routeDst := route.GetAttr("to")
		filterAttr := attributesFromRoute(route)

		// By default, the Trigger's subscriber is the "to" destination.
		// If a "condition" is set, a Filter object is interpolated
		// between the Trigger and the "to" destination.
		triggerSubsDst := routeDst

		if v := route.GetAttr("condition"); !v.IsNull() {
			const filterAPIGroup = "routing.triggermesh.io/v1alpha1"
			const filterKind = "Filter"

			triggerSubsDst = k8s.NewDestination(filterAPIGroup, filterKind, routeName)

			filter := k8s.NewObject(filterAPIGroup, filterKind, routeName)

			expr := v.AsString()
			filter.SetNestedField(expr, "spec", "expression")

			sink := k8s.DecodeDestination(routeDst)
			filter.SetNestedMap(sink, "spec", "sink", "ref")

			manifests = append(manifests, filter.Unstructured())
		}

		trigger := k8s.NewTrigger(routeName, name, triggerSubsDst, filterAttr)

		manifests = append(manifests, trigger)

		i++
	}

	return manifests
}

// Address implements translation.Addressable.
func (*ContentBased) Address(id string, _, _ cty.Value) cty.Value {
	return k8s.NewDestination(k8s.APIEventing, "Broker", k8s.RFC1123Name(id))
}

func attributesFromRoute(route cty.Value) map[string]interface{} {
	routeAttr := route.GetAttr("attributes")
	if routeAttr.IsNull() {
		return nil
	}

	filterAttr := make(map[string]interface{}, routeAttr.LengthInt())

	routeAttrIter := routeAttr.ElementIterator()
	for routeAttrIter.Next() {
		attr, val := routeAttrIter.Element()
		filterAttr[attr.AsString()] = val.AsString()
	}

	return filterAttr
}

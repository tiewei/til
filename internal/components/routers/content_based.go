package routers

import (
	"strconv"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/k8s"
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
			"attributes": &hcldec.AttrSpec{
				Name:     "attributes",
				Type:     cty.Map(cty.String),
				Required: true,
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

	namePrefix := k8s.RFC1123Name(id)

	broker := k8s.NewBroker(namePrefix)
	manifests = append(manifests, broker)

	i := 0
	routeIter := config.ElementIterator()
	for routeIter.Next() {
		_, route := routeIter.Element()

		routeDst := route.GetAttr("to")
		filterAttr := attributesFromRoute(route)

		name := namePrefix + "-r" + strconv.Itoa(i)
		trigger := k8s.NewTrigger(name, namePrefix, routeDst, filterAttr)

		manifests = append(manifests, trigger)

		i++
	}

	return manifests
}

// Address implements translation.Addressable.
func (*ContentBased) Address(id string) cty.Value {
	return k8s.NewDestination("eventing.knative.dev/v1", "Broker", k8s.RFC1123Name(id))
}

func attributesFromRoute(route cty.Value) map[string]interface{} {
	routeAttr := route.GetAttr("attributes")

	filterAttr := make(map[string]interface{}, routeAttr.LengthInt())

	routeAttrIter := routeAttr.ElementIterator()
	for routeAttrIter.Next() {
		attr, val := routeAttrIter.Element()
		filterAttr[attr.AsString()] = val.AsString()
	}

	return filterAttr
}

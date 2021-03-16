package core

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/graph"
	"bridgedl/translation"
)

// AttachableSpecVertex is implemented by all types used as graph.Vertex that
// can have a hcldec.Spec attached.
type AttachableSpecVertex interface {
	MessagingComponentVertex

	AttachSpec(hcldec.Spec)
}

// DecodableConfigVertex is implemented by all types used as graph.Vertex that
// may contain a HCL configuration body that can be decoded.
type DecodableConfigVertex interface {
	DecodedConfig(*hcl.EvalContext) (cty.Value, hcl.Diagnostics)

	// If a type can decode a configuration, it must also be able to attach
	// a decode spec.
	AttachableSpecVertex
}

// AttachSpecsTransformer is a GraphTransformer that attaches a decode spec to
// all graph vertices that support it.
type AttachSpecsTransformer struct{}

var _ GraphTransformer = (*AttachSpecsTransformer)(nil)

// Transform implements GraphTransformer.
func (t *AttachSpecsTransformer) Transform(g *graph.DirectedGraph) hcl.Diagnostics {
	var diags hcl.Diagnostics

	for _, v := range g.Vertices() {
		dec, ok := v.(DecodableConfigVertex)
		if !ok {
			continue
		}

		// some component types may not have any configuration body to
		// decode at all, in which case there is simply no spec to attach
		decType, ok := dec.Implementation().(translation.Decodable)
		if !ok {
			continue
		}

		dec.AttachSpec(decType.Spec())
	}

	return diags
}

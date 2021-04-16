package core

import (
	"github.com/hashicorp/hcl/v2"

	"bridgedl/config"
	"bridgedl/config/addr"
	"bridgedl/graph"
)

// MessagingComponentVertex is implemented by all messaging components of a
// Bridge which are represented by a graph.Vertex.
type MessagingComponentVertex interface {
	ComponentAddr() addr.MessagingComponent
	Implementation() interface{}
}

// AddComponentsTransformer is a GraphTransformer that adds all messaging
// components described in a Bridge as vertices of a graph, without connecting
// them.
type AddComponentsTransformer struct {
	Bridge *config.Bridge
}

var _ GraphTransformer = (*AddComponentsTransformer)(nil)

// Transform implements GraphTransformer.
func (t *AddComponentsTransformer) Transform(g *graph.DirectedGraph) hcl.Diagnostics {
	var diags hcl.Diagnostics

	for _, ch := range t.Bridge.Channels {
		v := &ChannelVertex{
			Addr: addr.Channel{
				Identifier: ch.Identifier,
			},
			Channel: ch,
		}
		g.Add(v)
	}

	for _, rtr := range t.Bridge.Routers {
		v := &RouterVertex{
			Addr: addr.Router{
				Identifier: rtr.Identifier,
			},
			Router: rtr,
		}
		g.Add(v)
	}

	for _, trsf := range t.Bridge.Transformers {
		v := &TransformerVertex{
			Addr: addr.Transformer{
				Identifier: trsf.Identifier,
			},
			Transformer: trsf,
		}
		g.Add(v)
	}

	for _, src := range t.Bridge.Sources {
		v := &SourceVertex{
			Source: src,
		}
		g.Add(v)
	}

	for _, trg := range t.Bridge.Targets {
		v := &TargetVertex{
			Addr: addr.Target{
				Identifier: trg.Identifier,
			},
			Target: trg,
		}
		g.Add(v)
	}

	return diags
}

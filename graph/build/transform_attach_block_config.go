package build

import (
	"github.com/hashicorp/hcl/v2"

	"bridgedl/graph"
	"bridgedl/translate"
)

// AttachTranslatorsTransformer is a GraphTransformer that attaches a block
// translator to all graph vertices that support it.
type AttachTranslatorsTransformer struct {
	Translators *translate.TranslatorProviders
}

var _ GraphTransformer = (*AttachTranslatorsTransformer)(nil)

// Transform implements GraphTransformer.
func (t *AttachTranslatorsTransformer) Transform(g *graph.DirectedGraph) hcl.Diagnostics {
	var diags hcl.Diagnostics

	for _, v := range g.Vertices() {
		attch, ok := v.(AttachableTranslatorVertex)
		if !ok {
			continue
		}

		transl := attch.FindTranslator(t.Translators)
		if transl == nil {
			continue
		}

		attch.AttachTranslator(transl)
	}

	return diags
}

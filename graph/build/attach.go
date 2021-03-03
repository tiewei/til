package build

import (
	"github.com/hashicorp/hcl/v2"

	"bridgedl/translate"
)

// AttachableTranslatorVertex is implemented by all types used as graph.Vertex
// that can have a translate.BlockTranslator attached.
type AttachableTranslatorVertex interface {
	TranslatorFinder
	AttachTranslator(translate.BlockTranslator)
}

// TranslatorFinder can look up a suitable translate.BlockTranslator in a
// collection of translate.TranslatorProviders.
type TranslatorFinder interface {
	FindTranslator(*translate.TranslatorProviders) (translate.BlockTranslator, hcl.Diagnostics)
}

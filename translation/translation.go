package translation

import "github.com/hashicorp/hcl/v2/hcldec"

// TranslatorProvider provides access to BlockTranslators for a set of
// component types.
//
// The scope of a TranslatorProvider applies to a single block type. That is,
// an instance of such provider supports *either* channels, or routers, etc.
// but not several of those types.
type TranslatorProvider interface {
	// Returns a translator for the given component type, if the provider
	// supports such type.
	Translator(componentType string) BlockTranslator
}

// BlockTranslator exposes the methods necessary for translating configuration
// blocks from a Bridge Description File into concrete resource definitions,
// such as Kubernetes objects.
type BlockTranslator interface {
	SpecProvider
	KubernetesTranslator
}

// SpecProvider provides access to specs that allow component-specific block
// configurations to be decoded into concrete values.
type SpecProvider interface {
	Spec() hcldec.Spec
}

// KubernetesTranslator translates component blocks into Kubernetes manifests.
type KubernetesTranslator interface {
	Manifests(interface{}) []interface{}
}

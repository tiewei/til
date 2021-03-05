package translation

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
	SchemaProvider
	KubernetesTranslator
}

// SchemaProvider provides access to Go types that represent the specific
// configurations of blocks in Bridge Description File.
//
// It is assumed that returned instances can be decoded from a hcl.Body using
// gohcl.DecodeBody. For instance, structs fields are expected to have `hcl`
// tags.
type SchemaProvider interface {
	ConcreteConfig() interface{}
}

// KubernetesTranslator translates instances of Go types into Kubernetes
// manifests.
type KubernetesTranslator interface {
	Manifests(interface{}) []interface{}
}

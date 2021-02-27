package translate

// TranslatorProvider exposes the types of components supported for a specific
// block type. Component type strings define the "type" labels that can be used
// in HCL blocks.
type TranslatorProvider interface {
	// Returns a translator for the given component type, if available.
	Translator(componentType string) BlockTranslator
}

// BlockTranslator translates HCL blocks into Kubernetes manifests.
type BlockTranslator interface {
	ConcreteConfig() interface{}
	K8SManifests() []interface{}
}

// TranslatorProviders encapsulates TranslatorProviders for all block types
// that support translation.
type TranslatorProviders struct {
	Channels     TranslatorProvider
	Routers      TranslatorProvider
	Transformers TranslatorProvider
	Sources      TranslatorProvider
	Targets      TranslatorProvider
	Functions    TranslatorProvider
}

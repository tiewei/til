package file

import (
	"fmt"
	"io"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"

	"bridgedl/config"
	"bridgedl/fs"
)

// Parser can parse Bridge Description Files and decode them into actual Bridge structs.
// The embedded FS interface is mainly useful for tests.
type Parser struct {
	*hclparse.Parser
	FS fs.FS
}

// NewParser returns an new Parser initialized with a fs.FS backed by the OS.
func NewParser() *Parser {
	return &Parser{
		Parser: hclparse.NewParser(),
		FS:     &fs.OSFS{},
	}
}

// LoadBridge parses the Bridge Description File at the given path and decodes
// it into a Bridge struct.
func (p *Parser) LoadBridge(filePath string) (*config.Bridge, hcl.Diagnostics) {
	hclFile, diags := p.ParseHCLFile(filePath)
	if diags.HasErrors() {
		return nil, diags
	}

	brg := &config.Bridge{}
	diags = decodeBridge(hclFile.Body, brg)

	return brg, diags
}

// ParseHCLFile reads and parses the contents of a HCL file.
//
// This method overrides (*hclparse.Parser).ParseHCLFile in order to use the
// embedded fs.FS interface instead of calling OS functions directly.
func (p *Parser) ParseHCLFile(filePath string) (*hcl.File, hcl.Diagnostics) {
	f, err := p.FS.Open(filePath)
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to open file",
			Detail: fmt.Sprintf("The configuration file %q could not be opened. "+
				"The error was: %s", filePath, err),
		}}
	}
	defer f.Close()

	src, err := io.ReadAll(f)
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to read file",
			Detail: fmt.Sprintf("The configuration file %q could not be read. "+
				"The error was: %s", filePath, err),
		}}
	}

	return p.ParseHCL(src, filePath)
}

/*
Copyright 2021 TriggerMesh Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/hashicorp/hcl/v2"

	"til/cli"
	"til/config/file"
	"til/core"
	"til/encoding"
	"til/graph/dot"
)

// CLI subcommands
const (
	cmdGenerate = "generate"
	cmdValidate = "validate"
	cmdGraph    = "graph"
)

// usage is a usageFn for the top level command.
func usage(cliName string) string {
	return "Interpreter for TriggerMesh's Integration Language.\n" +
		"\n" +
		"USAGE:\n" +
		"    " + cliName + " <command>\n" +
		"\n" +
		"COMMANDS:\n" +
		"    " + cmdGenerate + "     Generate Kubernetes manifests for deploying a Bridge.\n" +
		"    " + cmdValidate + "     Validate a Bridge description.\n" +
		"    " + cmdGraph + "        Represent a Bridge as a directed graph in DOT format.\n"
}

// usageGenerate is a usageFn for the "generate" subcommand.
func usageGenerate(cmd string) string {
	return "Generates the Kubernetes manifests which allow a Bridge to be deployed " +
		"to TriggerMesh, and writes them to standard output.\n" +
		"\n" +
		"USAGE:\n" +
		"    " + cmd + " FILE [OPTION]...\n" +
		"\n" +
		"OPTIONS:\n" +
		"    --bridge     Output a Bridge object instead of a List-manifest.\n" +
		"    --yaml       Output generated manifests in YAML format.\n"
}

// usageValidate is a usageFn for the "validate" subcommand.
func usageValidate(cmd string) string {
	return "Verifies that a Bridge is syntactically valid and can be generated. " +
		"Returns with an exit code of 0 in case of success, with an exit code of 1 " +
		"otherwise.\n" +
		"\n" +
		"USAGE:\n" +
		"    " + cmd + " FILE\n"
}

// usageGraph is a usageFn for the "usage" subcommand.
func usageGraph(cmd string) string {
	return "Generates a DOT representation of a Bridge and writes it to standard " +
		"output.\n" +
		"\n" +
		"USAGE:\n" +
		"    " + cmd + " FILE\n"
}

// usageFn returns the usage text for a program or subcommand.
type usageFn func(cmd string) string

// setUsageFn uses the given usageFn to set the Usage function of the provided
// flag.FlagSet.
func setUsageFn(f *flag.FlagSet, u usageFn) {
	f.Usage = func() {
		fmt.Fprint(f.Output(), u(f.Name()))
	}
}

var (
	_ cli.Command = (*GenerateCommand)(nil)
	_ cli.Command = (*ValidateCommand)(nil)
	_ cli.Command = (*GraphCommand)(nil)
)

type GenerateCommand struct {
	// flags
	bridge bool
	yaml   bool
}

// Run implements cli.Command.
func (c *GenerateCommand) Run(ctx context.Context, args []string) error {
	flagSet := cli.FlagSetFromContext(ctx)
	setUsageFn(flagSet, usageGenerate)

	flagSet.BoolVar(&c.bridge, "bridge", false, "")
	flagSet.BoolVar(&c.yaml, "yaml", false, "")

	pos, flags := splitArgs(1, args)
	_ = flagSet.Parse(flags) // ignore err; the FlagSet uses ExitOnError

	if len(pos) != 1 {
		return fmt.Errorf("unexpected number of positional arguments.\n\n%s", usageGenerate(flagSet.Name()))
	}
	filePath := pos[0]

	// value to use as the Bridge identifier in case none is defined in the
	// parsed Bridge description
	const defaultBridgeIdentifier = "til_generated"

	ui := cli.UIFromContext(ctx)

	p := file.NewParser()
	brg, diags := p.LoadBridge(filePath)
	dw := newDiagnosticTextWriter(ui.ErrWriter, p.Files())
	if diags.HasErrors() {
		_ = dw.WriteDiagnostics(diags)
		return errLoadBridge
	}

	cctx, diags := core.NewContext(brg)
	if diags.HasErrors() {
		_ = dw.WriteDiagnostics(diags)
		return errInitContext
	}

	manifests, diags := cctx.Generate()
	if diags.HasErrors() {
		_ = dw.WriteDiagnostics(diags)
		return errGenerate
	}

	brgID := brg.Identifier
	if brgID == "" {
		brgID = defaultBridgeIdentifier
	}
	s := encoding.NewSerializer(brgID)

	var w encoding.ManifestsWriterFunc

	switch {
	case c.bridge && c.yaml:
		w = s.WriteBridgeYAML
	case c.bridge:
		w = s.WriteBridgeJSON
	case c.yaml:
		w = s.WriteManifestsYAML
	default:
		w = s.WriteManifestsJSON
	}

	return w(ui.StdWriter, manifests)
}

type ValidateCommand struct{}

// Run implements Command.
func (c *ValidateCommand) Run(ctx context.Context, args []string) error {
	flagSet := cli.FlagSetFromContext(ctx)
	setUsageFn(flagSet, usageValidate)

	pos, flags := splitArgs(1, args)
	_ = flagSet.Parse(flags) // ignore err; the FlagSet uses ExitOnError

	if len(pos) != 1 {
		return fmt.Errorf("unexpected number of positional arguments.\n\n%s", usageValidate(flagSet.Name()))
	}
	filePath := pos[0]

	ui := cli.UIFromContext(ctx)

	p := file.NewParser()
	brg, diags := p.LoadBridge(filePath)
	dw := newDiagnosticTextWriter(ui.ErrWriter, p.Files())
	if diags.HasErrors() {
		_ = dw.WriteDiagnostics(diags)
		return errLoadBridge
	}

	cctx, diags := core.NewContext(brg)
	if diags.HasErrors() {
		_ = dw.WriteDiagnostics(diags)
		return errInitContext
	}

	if _, diags := cctx.Generate(); diags.HasErrors() {
		_ = dw.WriteDiagnostics(diags)
		return errGenerate
	}

	return nil
}

type GraphCommand struct{}

// Run implements Command.
func (c *GraphCommand) Run(ctx context.Context, args []string) error {
	flagSet := cli.FlagSetFromContext(ctx)
	setUsageFn(flagSet, usageGraph)

	pos, flags := splitArgs(1, args)
	_ = flagSet.Parse(flags) // ignore err; the FlagSet uses ExitOnError

	if len(pos) != 1 {
		return fmt.Errorf("unexpected number of positional arguments.\n\n%s", usageGraph(flagSet.Name()))
	}
	filePath := pos[0]

	ui := cli.UIFromContext(ctx)

	p := file.NewParser()
	brg, diags := p.LoadBridge(filePath)
	dw := newDiagnosticTextWriter(ui.ErrWriter, p.Files())
	if diags.HasErrors() {
		_ = dw.WriteDiagnostics(diags)
		return errLoadBridge
	}

	cctx, diags := core.NewContext(brg)
	if diags.HasErrors() {
		_ = dw.WriteDiagnostics(diags)
		return errInitContext
	}

	g, diags := cctx.Graph()
	if diags.HasErrors() {
		_ = dw.WriteDiagnostics(diags)
		return errors.New("failed to build bridge graph. See error diagnostics")
	}

	dg, err := dot.Marshal(g)
	if err != nil {
		return fmt.Errorf("marshaling graph to DOT: %w", err)
	}

	stdout := cli.UIFromContext(ctx).StdWriter
	if _, err := stdout.Write(dg); err != nil {
		return fmt.Errorf("writing generated DOT graph: %w", err)
	}

	return nil
}

// splitArgs attempts to separate n positional arguments from the rest of the
// given arguments list. The caller is responsible for ensuring that the
// correct number of positional arguments could be extracted.
//
// It is meant as a helper to implement CLI commands of the shape:
//   cmd ARG1 ARG2 [flags]
func splitArgs(n int, args []string) ( /*positional*/ []string /*flags*/, []string) {
	if len(args) == 0 {
		return nil, nil
	}

	// no positional, or user passed only flags (e.g. "cmd -h")
	if n == 0 || (args[0] != "" && args[0][0] == '-') {
		return nil, args
	}

	if len(args) <= n {
		return args, nil
	}

	return args[:n], args[n:]
}

// newDiagnosticTextWriter returns a hcl.DiagnosticWriter that writes
// diagnostics to the given writer as formatted text.
func newDiagnosticTextWriter(out io.Writer, files map[string]*hcl.File) hcl.DiagnosticWriter {
	const outputWidth = 0
	const enableColor = true
	return hcl.NewDiagnosticTextWriter(out, files, outputWidth, enableColor)
}

// Errors for common operations performed by commands.
var (
	errLoadBridge  = errors.New("failed to load bridge. See error diagnostics")
	errInitContext = errors.New("failed to initialize command context. See error diagnostics")
	errGenerate    = errors.New("failed to generate bridge manifests. See error diagnostics")
)
